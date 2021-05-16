package bouncer

import "os/exec"

var (
	// execCmd represents the execution environment for commands
	execCmd ExecCmd = execShellCmd
)

// ExecCmd represents the call to execute a command
type ExecCmd func(name string, arg ...string) (out []byte, err error)

// execShellCmd performs the actual command execution via exec.Command(..)
func execShellCmd(name string, arg ...string) (out []byte, err error) {
	cmd := exec.Command(name, arg...)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}
