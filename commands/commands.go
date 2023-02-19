package commands

import (
	"fmt"
	"os/exec"
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

func CreateAndPush() error {
	err := exec.Command("git", "add", ".").Run()
	if err != nil {
		fmt.Printf("Error occured adding changes %v", err)
		return err
	}
	err = exec.Command("git", "commit", "-m", fmt.Sprintf("")).Run()

	if err != nil {
		fmt.Printf("Error occured making commit %v", err)
	}

	err = exec.Command("git", "push")
}