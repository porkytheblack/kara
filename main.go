package main

import (
	"fmt"
	"kara/menus"
)


func main() {

	menus := &menus.MenuInterface{}

	menus.InitMenu()

	menus.ConstructComponentsList()
	menus.ConstructInitialScreen()
	

	pages := menus.Pages
	
	err := menus.App.SetRoot(pages, true).SetFocus(pages).EnableMouse(true).Run()

	if err != nil {
		fmt.Printf("An error occured %v", err)
	}


	
}