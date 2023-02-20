package utils

import (
	"fmt"
	"os/exec"
	"strings"
)


func HasUpstream (branch string) bool {

	upstream, err := exec.Command("git", "config", fmt.Sprintf("branch.%s.remote",strings.TrimSpace(branch))).CombinedOutput()

	if err != nil {
		fmt.Printf("An error occured")
		return false
	}

	if string(upstream) == "" {
		return false
	}

	return true

}

func CreateCommitMessage (Type string, Name string, Description string) string {
	str := fmt.Sprintf(`%s:%s-%s`, Type, Name, Description)
	return str
}

func BranchExists (branch_name string) bool {
	_, err := exec.Command("git", "show-ref", "--verify", fmt.Sprintf("refs/heads/%s", branch_name)).CombinedOutput()
	return err == nil
}