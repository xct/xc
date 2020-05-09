package shell

import (
	"net"
	"os/exec"
)

/*
 It would be awesome if someone implemented this as well ;-)
*/

// Shell ...
func Shell() *exec.Cmd {
	panic("Not implemented")
}

// Exec ...
func Exec(command string, conn net.Conn) {
	panic("Not implemented")
}

// RunAsPs ...
func RunAs(user string, pass string, domain string, c net.Conn) {
	panic("Not implemented")
}

// RunAsPs ...
func RunAsPS(user string, pass string, domain string, c net.Conn) {
	panic("Not implemented")
}

// ExecSC ....
func ExecSC(sc []byte) {
	panic("Not implemented")
}

// ExecOut ...
func ExecOut(command string) (string, error) {
	panic("Not implemented")
}

// ExecPSOut ...
func ExecPSOut(command string) (string, error) {
	panic("Not implemented")
}

// ExecDebug ...
func ExecDebug(cmd string) (string, error) {
	panic("Not implemented")
}

// ExecPSDebug ...
func ExecPSDebug(cmd string) (string, error) {
	panic("Not implemented")
}

// Powershell ...
func Powershell() *exec.Cmd {
	panic("Not implemented")
}

// CopySelf ...
func CopySelf() string {
	panic("Not implemented")
}

// ExecSilent ...
func ExecSilent(cmd string, c net.Conn) {
	panic("Not implemented")
}

// Seppuku ...
func Seppuku(c net.Conn) {
	panic("Not implemented")
}
