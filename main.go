package main

import (
	// "kara/menus"
	// "fmt"
	"kara/menus"
	"kara/repository"
	"os"
	// "os"
)


func main() {

	repo := &repository.Repository{}
	repo.InitRepository()
	err := repo.InitUser()

	if err != nil {
		if err.Error() == "no token" {
			repo.AuthenticateUser()
		}else {
			os.Exit(1)
		}
	}
	menu := &menus.MenuInterface{}
	menu.InitMenu()
	repo.GetArgs(menu)
	
	// pages := menu.Pages
	// err = menu.App.SetRoot(pages, true).SetFocus(pages).EnableMouse(true).Run()

	// if err != nil {
	// 	fmt.Printf("\n\nAn error occured:: %v", err)
	// }

}