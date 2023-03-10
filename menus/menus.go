package menus

import (
	// "fmt"

	"fmt"
	"kara/commands"
	"kara/utils"

	"github.com/rivo/tview"
)

type MenuInterface struct  {
	Pages *tview.Pages
	App *tview.Application
	CurrentPage string
}

func (menu *MenuInterface) InitMenu () {
	menu.CurrentPage = "Instruction Page"
	menu.Pages = tview.NewPages()
	menu.App = tview.NewApplication()

	
}

func (menu *MenuInterface) ConstructInitialScreen () {

	options := tview.NewList().
	// AddItem("Updated Commit(Component)", "",'1', func() {
	// 	// fmt.Printf("Selected Update")
	// 	menu.Pages.SwitchToPage("ComponentList")
	// 	menu.CurrentPage = "ComponentList"
	// }).
	AddItem("Create Commit(Component)", "", '1', func() {
		// fmt.Printf("Selected create commit")
		menu.ConstructUpdateComponent("Create")
	}).
	AddItem("Switch to Feature(Branch)", "", '2', func() {
		// fmt.Printf("Selected create branch")
		menu.ConstructChooseBranch()
	}).
	AddItem("Create Feature(Branch)", "", '3', func() {
		// fmt.Printf("Selected create branch")
		menu.ConstructCreateBranch()
	}).
	AddItem("Exit", "", '4', func() {
		menu.App.Stop()
	})


	menu.Pages.AddPage("Instruction Page", options, true, menu.CurrentPage == "Instruction Page")

}

func (menu *MenuInterface) ConstructComponentsList() {
	options := tview.NewList().
	AddItem("Component A", "Component Description", '1', func() {
		commands.CreateAndPush("component a", "feature", "component", "some stuff", true)
	}).
	AddItem("Component B", "Component Description", '2', nil).
	AddItem("Component C", "Component Description", '3', nil).
	AddItem("Component D", "Component Description", '4', nil)

	err := menu.App.SetRoot(options, true).SetFocus(options).Run()

	if err != nil {
		panic(err)
	}
}


func (menu *MenuInterface) ConstructRemoteBranchList() {
	options := tview.NewList(). 
	AddItem("Branch A", "Branch A", '1', func () {
		commands.CreateAndChangeBranch("branch_name", "Initial Commit")
	}).
	AddItem("Branch B", "Branch B", '2', func () {

	})

	err := menu.App.SetRoot(options, true).SetFocus(options).Run()

	if err != nil {
		panic(err)
	}
}

func (menu *MenuInterface) ConstructUpdateComponent(formType string) {
	var description string;
	var name string;
	var commit_type string;
	form := tview.NewForm().
	AddInputField("Component Name", "", 40, nil, func(text string){
		name = text
	}).
	AddTextArea("Description", "", 40,0, 0, func (text string){
		description = text
	}). 
	AddDropDown("Type", []string{"Feature", "Fix", "Chore", "Refactor"}, 0, func(option string, optionIndex int){
		commit_type = option
	}).
	AddButton("Done", func() {
		err := commands.CreateAndPush(description, commit_type, name, "", true)

		if err != nil {
			fmt.Printf("An error occured %v", err)
			menu.App.Stop()
		}

		menu.App.Stop()
	}). 
	AddButton("Cancel", func () {
		menu.App.Stop()
	})
	
	menu.Pages.AddAndSwitchToPage("UpdateComponent", form, true)

}

func (menu *MenuInterface) ConstructCreateBranch() {
	var branch_name string;
	var description string;
	form := tview.NewForm().
	AddInputField("Feature Name", "", 40, nil, func(text string){
		branch_name = text
	}).
	AddTextArea("Feature Description", "", 40, 0, 0,func (text string) {
		description = text
	}).
	AddButton("Done", func() {
		err := commands.CreateAndChangeBranch(branch_name, utils.CreateCommitMessage("new_branch", "init", description))

		if err != nil {
			fmt.Printf("An error occured : %s", err)
			menu.App.Stop()
		}

		menu.App.Stop()
	}).
	AddButton("Cancel", func () {
		menu.App.Stop()
	})

	menu.Pages.AddAndSwitchToPage("CreateFeatureBranch", form, true)

}

func (menu *MenuInterface) ConstructChooseBranch() {
	list := tview.NewList()
	
	for i, value := range (commands.GetLocalBranches()) {
		branch := value
		list.AddItem(branch, "", rune(fmt.Sprintf(`%c`, i)[0]), func() {
			err := commands.SwitchBranch(branch)
			if err != nil {
				menu.App.Stop()
				fmt.Printf("An error occured switching to %s ::%s",branch, err)
			}
			menu.App.Stop()
		})
	}

	err := menu.App.SetRoot(list, true).SetFocus(list).Run()

	if err != nil {
		panic(err)
	}
}

