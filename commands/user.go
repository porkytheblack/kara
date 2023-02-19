package commands

import "os/exec"

func CurrentUser() string {
	user_identifier, err := exec.Command("git", "config", "user.email").CombinedOutput()

	if err != nil {
		return ""
	}
	str := string(user_identifier)
	return str
}