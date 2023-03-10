package commands

import (
	"errors"
	"fmt"
	"kara/utils"
	"os/exec"
	"strings"
)

func CurrentRepository() (string, error) {

	url, err := exec.Command("git", "remote", "get-url", "origin").CombinedOutput()

	if err != nil {
		// 2  : origin branch does not exist
		fmt.Printf("Current error:: %s", err)
		return "", err
	}

	return string(url), nil
}


func CurrentBranch() (string, error) {
	branch, err := exec.Command("git", "branch", "--show-current").CombinedOutput()

	if err != nil {
		fmt.Printf("Error getting current branch:: %v", err)
		return "", err
	}

	return string(branch), nil;

}

func CurrentCommit() (string, error) {
	branch, err := exec.Command("git", "rev-parse", "HEAD").CombinedOutput()

	if err != nil {
		fmt.Printf("Error getting current commit hash %v", err )
		return "", err
	}

	return string(branch), nil
}

func PushToRemote () error {
	fmt.Print("Pushing to remote")
	branch, _ := CurrentBranch()
	fmt.Printf("Current Branch is %s", branch)
	if !utils.HasUpstream(branch) {
		err := exec.Command("git", "push", "--set-upstream", "origin", branch).Run()
		fmt.Printf("Set upsteam and pushed %v", err)
		return err
	}
	err := exec.Command("git", "push").Run()
	fmt.Printf("Pushed %v", err)
	if err != nil {
		fmt.Printf("An error occured while trying to push %v", err)
		return err
	}

	return nil
}

func CreateAndPush(commit_message string, commit_type string, commit_name string, description string, push bool) error {
	err := exec.Command("git", "add", ".").Run()
	if err != nil {
		fmt.Printf("Error occured adding changes %v", err)
		return err
	}
	err = exec.Command("git", "commit" , "-m", utils.CreateCommitMessage(commit_type, commit_name, commit_message), "-m", fmt.Sprintf(`
		# Description
		%s
	`, description)).Run()

	if err != nil {
		fmt.Printf("Error occured making commit %v", err)
		return err
	}
	fmt.Println("Starting the push to remote")

	if push {
		err = PushToRemote()

		if err != nil {
			fmt.Printf("An error occured %s", err)
			return err
		}
	}
	

	return nil
}

func ChangeBranch (branch_name string) error {

	if !utils.BranchExists(branch_name) {
		return errors.New("branch does not exist")
	}

	err := exec.Command("git", "checkout", branch_name).Run()

	if err != nil {
		fmt.Printf("An Error occured %v", err)
		return fmt.Errorf("an Error occured %v", err)
	}

	return nil
}


func SwitchBranch (branch_name string) error {

	current_branch, _ := CurrentBranch()

	switch branch_name {
		case current_branch:
			return nil
		default:
			return ChangeBranch(branch_name)
	}
}

func CreateAndChangeBranch (branch_name string, message string) error {
	if utils.BranchExists(branch_name) {
		return errors.New("branch exists")
	}
	
	err := exec.Command("git", "checkout", "-b", branch_name).Run()

	if err != nil {
		return err
	}

	err = exec.Command("git", "commit", "--allow-empty", "-m", message).Run()

	if err != nil {
		return err
	}

	err = exec.Command("git", "push", "-u", "origin", branch_name).Run()

	if err != nil {
		return err
	}


	return nil
}


func GetLocalBranches() []string {
	br, err := exec.Command("git", "branch").CombinedOutput()

	if err != nil {
		fmt.Printf("An Error occured fetching branches")
		return []string{}
	}

	lines := strings.Split(string(br), "\n")

	branches := make([]string, 0, len(lines))

	for _, line := range lines {
		branch := strings.TrimPrefix(strings.TrimSpace(line), "* ")
		
		if branch != "" {
			branches = append(branches, branch)
		}
	}

	return branches
}