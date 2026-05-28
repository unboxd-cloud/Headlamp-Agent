package setup

import (
	"os/exec"
)

type CheckResult struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func CheckCommand(id string, name string, binary string) CheckResult {
	path, err := exec.LookPath(binary)
	if err != nil {
		return CheckResult{ID: id, Name: name, OK: false, Message: binary + " not found"}
	}

	return CheckResult{ID: id, Name: name, OK: true, Message: path}
}

func LocalOperatorChecks() []CheckResult {
	return []CheckResult{
		CheckCommand("go", "Go", "go"),
		CheckCommand("podman", "Podman", "podman"),
	}
}
