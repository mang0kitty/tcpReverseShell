package rsh

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type PowerShell struct {
	powerShell string
}
func Powershell() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}
func (p *PowerShell) Execute(args ...string) error {
	cmd := run(exec.Command(p.powerShell, args...))
	err :=cmd.Run()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

    return err

}