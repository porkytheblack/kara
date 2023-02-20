package main

import (
	"fmt"
	"kara/menus"
	"kara/repository"
)


func main() {

	repo := &repository.Repository{}

	repo.InitRepository()
	repo.RemoveUser()
	err := repo.InitUser()

	if err != nil {
		if err.Error() == "no token" {
			repo.AuthenticateUser()
			
		} else {
			return
		}
	}
	

	menus := &menus.MenuInterface{}
	menus.InitMenu()
	menus.ConstructComponentsList()
	menus.ConstructInitialScreen()

	pages := menus.Pages
	
	err = menus.App.SetRoot(pages, true).SetFocus(pages).EnableMouse(true).Run()

	if err != nil {
		fmt.Printf("An error occured %v", err)
	}


	
}