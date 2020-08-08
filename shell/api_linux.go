package shell

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

func CreateProcessAsUser(name string, path string, cmdline string) error {
	args := strings.Split(cmdline, " ")

	user, err := user.Lookup(name)
	if err != nil {
		return err
	}
	uid, err := strconv.ParseUint(user.Uid, 10, 32)

	var cred = &syscall.Credential{uint32(uid), uint32(uid), []uint32{}, false}
	var sysproc = &syscall.SysProcAttr{Credential: cred, Noctty: true}
	var attr = os.ProcAttr{
		Dir: ".",
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			nil,
			nil,
		},
		Sys: sysproc,
	}
	proc, err := os.StartProcess(path, args, &attr)
	if err == nil {
		err = proc.Release()
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
	return nil
}
