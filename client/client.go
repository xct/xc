package client

import (
	"fmt"
	"log"
	"net"
	"os"

	"../utils"
	"github.com/hashicorp/yamux"
)

func prompt(c net.Conn) {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "?"
	}
	fmt.Fprintf(c, fmt.Sprintf("[xc: %s]: ", cwd))
}

func lfwd(host string, port string, s *yamux.Session, c net.Conn) {
	for {
		proxy, err := s.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		fwdCon, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		if err != nil {
			log.Println(err)
			return
		}
		defer fwdCon.Close()
		go utils.CopyIO(fwdCon, proxy)
		go utils.CopyIO(proxy, fwdCon)
	}
}

// opens the listening socket on the client side
func rfwd(port string, session *yamux.Session, c net.Conn) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println(err)
	}
	c.Write([]byte(fmt.Sprintf("Client listening on %s\n", port)))
	for {
		fwdCon, err := ln.Accept()
		defer fwdCon.Close()
		if err != nil {
			log.Println(err)
			c.Write([]byte(fmt.Sprintf("Remote forwarding caused an error %s\n", err)))
		}
		proxy, err := session.Open()
		if err != nil {
			log.Println(err)
		}
		go utils.CopyIO(fwdCon, proxy)
		go utils.CopyIO(proxy, fwdCon)
	}
}
