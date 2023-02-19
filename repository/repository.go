package repository

import (
	"bufio"
	"errors"
	"fmt"
	"kara/commands"
	"log"
	"os"
	"strings"

	"github.com/boltdb/bolt"
	"golang.org/x/crypto/ssh/terminal"
)

type Repository struct {
	DB	*bolt.DB
	UserBucket	*bolt.Bucket
	CurrentJwt	*string
}

func (repo *Repository) InitRepository () {

	db, err := bolt.Open("karadb", 0600, nil)
	// defer db.Close()
	if err != nil {
		log.Fatalf("Unable to initialize %s", err)
	}

	err = db.Update(func (tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("users"))

		if err != nil {
			log.Printf("An error occured: %s", err)
			return err
		}

		repo.UserBucket = bucket
		return nil

	})

	if err != nil {
		log.Fatalf("Unable to initialize user db")
	}

	repo.DB = db
}



func (repo *Repository) InitUser( ) error {
	user := commands.CurrentUser()

	err := repo.DB.View(func ( tx *bolt.Tx ) error {
		token := string(tx.Bucket([]byte("users")).Get([]byte(user)))
		fmt.Printf("token is %s", token)
		if token == "" {
			// fmt.Println("Error with token")
			return errors.New("no token")
		}

		repo.CurrentJwt = &token
		return nil
	})
	fmt.Println(err)
	if err != nil {
		// fmt.Printf("Error initializing user:: %s", err )
		return err
	}


	return nil
}


func (repo *Repository) AuthenticateUser( ) {
	reader := bufio.NewReader(os.Stdin)
	var mail string;
	var pass string;

	for {

		fmt.Printf("Enter your email:: ")
		email, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input error!")
			continue;
		}

		email = strings.TrimSuffix(email, "\n")

		if email == "" {
			fmt.Println("Email is required")
			continue
		}
		mail = email;

		break
		
	}

	for {
		fmt.Printf("Enter your password:: ")
		bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))

		if err != nil {
			fmt.Println("Input error")
			continue
		}


		if string(bytePassword) == "" {
			fmt.Println("Password is required")
			continue
		}

		pass = string(bytePassword)

		break
	}

	fmt.Printf("\n%s\n%s", mail,pass)

	// make an api request to get token

	user := commands.CurrentUser()
	

	err := repo.DB.Update(func(tx *bolt.Tx) error {

		err := tx.Bucket([]byte("users")).Put([]byte(user), []byte("token got from login request"))

		if err != nil {
			fmt.Printf("An error occured %s", err)
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("\nUnable to log user in %s", err)
		os.Exit(1)
	}

	fmt.Println("\nSuccessfully Authenticated user")

	
}


func (repo *Repository) RemoveUser() error {
	user := commands.CurrentUser()
	err := repo.DB.Update(func (tx *bolt.Tx) error {
		err := tx.Bucket([]byte("users")).Delete([]byte(user))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}