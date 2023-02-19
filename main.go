package main

import (
	"kara/repository"
)

func main(){
	repo := &repository.Repository{}
	repo.InitRepository()
	repo.RemoveUser()
	err := repo.InitUser()
	if err != nil {
		if err.Error() == string("no token") {
			repo.AuthenticateUser()
		}
	}

}	

