package repository

import (
	"bufio"
	"errors"
	"fmt"
	"kara/commands"
	"kara/menus"
	"log"
	"os"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/urfave/cli/v2"
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

	fmt.Printf("\n\n\n%s\n\n\n%s", mail,pass)

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

	fmt.Println("\n\nSuccessfully Authenticated user")	
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

func (repo *Repository) GetArgs (menus *menus.MenuInterface) {
	app := &cli.App{
		Name: "Kara",
		Usage: "Github Conventional Commits",
		Commands: []*cli.Command{
			{
				Name: "authenticate",
				Aliases: []string{"auth", "-a"},
				Usage: "Authenticate User",
				Action: func(c *cli.Context) error {
					repo.AuthenticateUser()
					return nil
				},
			},
			{
				Name: "create",
				Usage: "Create",
				Subcommands: []*cli.Command{
					{
						Name: "component",
						Usage: "Component",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "title",
								Aliases: []string{"t"},
								Usage: " Title ",
							},
							&cli.StringFlag{
								Name: "message",
								Aliases: []string{"m"},
								Usage: "Message",
								Value: "",
							},	
							&cli.StringFlag{
								Name: "category",
								Aliases: []string{"c"},
								Usage: "Category e.g feature | chore | refactor | fix",
								Value: "feature",
							},
							&cli.BoolFlag{
								Name: "remote",
								Aliases: []string{"r"},
								Usage: "Create from a remote component",
							},
							&cli.StringFlag{
								Name: "description",
								Aliases:  []string{"d"},
								Usage: "Add a description to the component creation",
							},
						},
						Action: func(ctx *cli.Context) error {
							use_remote := ctx.Bool("remote")
							t := ctx.Value("title")
							title := strings.TrimSpace(fmt.Sprintf("%s",t))
							m := ctx.Value("message")
							message := strings.TrimSpace(fmt.Sprintf("%s", m))
							c := ctx.Value("category")
							category := strings.TrimSpace(fmt.Sprintf("%s", c))
							d := ctx.Value("description")
							description := strings.TrimSpace(fmt.Sprintf("%s", d))
						
							if use_remote {
								menus.ConstructComponentsList()
							} else {
								err := commands.CreateAndPush(message, category, title, description, false)
								if err != nil {
									fmt.Printf("An error occured, Unable to commit %s", err)
									os.Exit(1)
								}
							}

							return nil
						},
					},
					{
						Name: "feature",
						Usage: "Create a feature(branch)",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "name",
								Aliases: []string{"n"},
								Usage: "Feature Name",
							},
							&cli.StringFlag{
								Name: "description",
								Aliases: []string{"d"},
								Usage: "Feature Description",
							},
							&cli.StringFlag{
								Name: "category",
								Aliases: []string{"c"},
								Usage: "Feature Category e.g feature | chore | refactor | fix",
								Value: "feature",
							},
							&cli.BoolFlag{
								Name: "remote",
								Aliases: []string{"r"},
								Usage: "Create from remote",
							},
						},
						Action: func(ctx *cli.Context) error {
							use_remote := ctx.Bool("remote")
							n := ctx.Value("name")
							d := ctx.Value("description")
							c := ctx.Value("category")
							name := fmt.Sprintf("%s", n)
							description := fmt.Sprintf("%s", d)
							category := fmt.Sprintf("%s", c)
							if use_remote {
								menus.ConstructRemoteBranchList()
							} else {
								commands.CreateAndChangeBranch(fmt.Sprintf("%s/%s", category, name), description)
							}

							// get all feature stuff
							return nil
						},
					},
				},
			},
			{
				Name: "switch",
				Usage: "switch to a feature(branch)",
				Action: func(ctx *cli.Context) error {
					menus.ConstructChooseBranch()
					return nil
				},
			},
		},
	}	

	err := app.Run(os.Args)

	if err != nil {
		fmt.Printf("An error occured %v", err)
		os.Exit(1)
	}

	

}