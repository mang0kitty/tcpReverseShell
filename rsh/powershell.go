package rsh

import (
	"io"
	"os/exec"
)

type AppRunner struct {
	Stdin io.Reader
	Stdout io.Writer
}

func (p *AppRunner) Execute(app string, args ...string) error {
	ps, err := exec.LookPath(app)
	if err != nil {
		return err
	}

	cmd := exec.Command(ps, args...)

	cmd.Stdin = p.Stdin
	cmd.Stdout = p.Stdout
	cmd.Stderr = p.Stdout

	return cmd.Run()
}