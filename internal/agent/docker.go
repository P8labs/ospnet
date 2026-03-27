package agent

import (
	"fmt"
	"os/exec"
)

func RunContainer(req RunRequest) error {
	cmd := exec.Command(
		"docker", "run", "-d",
		"-p", fmt.Sprintf("%d:80", req.Port),
		"--name", req.Name,
		req.Image,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker error: %s", string(output))
	}

	return nil
}
