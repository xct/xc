package shell

import (
	"os"
	"regexp"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	advapi32                    = syscall.NewLazyDLL("advapi32.dll")
	procCreateProcessWithLogonW = advapi32.NewProc("CreateProcessWithLogonW")
)

const (
	logonWithProfile        uint32 = 0x00000001
	logonNetCredentialsOnly uint32 = 0x00000002
	createDefaultErrorMode  uint32 = 0x04000000
	createNewProcessGroup   uint32 = 0x00000200
)

// CreateProcessWithLogonW ...
func CreateProcessWithLogonW(
	username *uint16,
	domain *uint16,
	password *uint16,
	logonFlags uint32,
	applicationName *uint16,
	commandLine *uint16,
	creationFlags uint32,
	environment *uint16,
	currentDirectory *uint16,
	startupInfo *syscall.StartupInfo,
	processInformation *syscall.ProcessInformation) error {
	r1, _, e1 := procCreateProcessWithLogonW.Call(
		uintptr(unsafe.Pointer(username)),
		uintptr(unsafe.Pointer(domain)),
		uintptr(unsafe.Pointer(password)),
		uintptr(logonFlags),
		uintptr(unsafe.Pointer(applicationName)),
		uintptr(unsafe.Pointer(commandLine)),
		uintptr(creationFlags),
		uintptr(unsafe.Pointer(environment)), // env
		uintptr(unsafe.Pointer(currentDirectory)),
		uintptr(unsafe.Pointer(startupInfo)),
		uintptr(unsafe.Pointer(processInformation)))
	runtime.KeepAlive(username)
	runtime.KeepAlive(domain)
	runtime.KeepAlive(password)
	runtime.KeepAlive(applicationName)
	runtime.KeepAlive(commandLine)
	runtime.KeepAlive(environment)
	runtime.KeepAlive(currentDirectory)
	runtime.KeepAlive(startupInfo)
	runtime.KeepAlive(processInformation)
	if int(r1) == 0 {
		return os.NewSyscallError("CreateProcessWithLogonW", e1)
	}
	return nil
}

// ListToEnvironmentBlock ...
func ListToEnvironmentBlock(list *[]string) *uint16 {
	if list == nil {
		return nil
	}
	size := 1
	for _, v := range *list {
		size += len(syscall.StringToUTF16(v))
	}
	result := make([]uint16, size)
	tail := 0
	for _, v := range *list {
		uline := syscall.StringToUTF16(v)
		copy(result[tail:], uline)
		tail += len(uline)
	}
	result[tail] = 0
	return &result[0]
}

// CreateProcessWithLogon creates a process giving user credentials
// Ref: https://github.com/hosom/honeycred/blob/master/honeycred.go
func CreateProcessWithLogon(username string, password string, domain string, path string, cmdLine string) error {
	user := syscall.StringToUTF16Ptr(username)
	dom := syscall.StringToUTF16Ptr(domain)
	pass := syscall.StringToUTF16Ptr(password)
	logonFlags := logonWithProfile // changed
	applicationName := syscall.StringToUTF16Ptr(path)
	commandLine := syscall.StringToUTF16Ptr(cmdLine)
	creationFlags := createDefaultErrorMode
	environment := ListToEnvironmentBlock(nil)
	currentDirectory := syscall.StringToUTF16Ptr(`c:\programdata`)
	startupInfo := &syscall.StartupInfo{}
	processInfo := &syscall.ProcessInformation{}

	err := CreateProcessWithLogonW(
		user,
		dom,
		pass,
		logonFlags,
		applicationName,
		commandLine,
		creationFlags,
		environment,
		currentDirectory,
		startupInfo,
		processInfo)
	return err
}

// GetBuild ...
func GetBuild(raw string) string {
	// Microsoft Windows [Version 10.0.18363.778]
	var re = regexp.MustCompile(`(?P<build>[\d+\.]+)`)
	version := re.FindString(raw)
	return version
}

// GetHotfixes ...
func GetHotfixes(raw string) []string {
	// HOSTNAME Update KB4537572 NT AUTHORITY\SYSTEM 3/31/2020 12:00:00 AM
	kbs := []string{}
	var re = regexp.MustCompile(`(?m)(?P<kb>KB\d+)`)
	for _, match := range re.FindAllString(raw, -1) {
		kbs = append(kbs, match)
	}
	return kbs
}
