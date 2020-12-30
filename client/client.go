package client

import (
	"fmt"
	"log"
	"net"
	"os"

	"../utils"
	"github.com/hashicorp/yamux"
)

var activeForwards []utils.Forward

// splitArgs
func splitArgs(cmd string) []string {
	args := []string{}
	current := ""
	squote := 0
	dquote := 0
	last := ""
	for _, rune := range cmd {
		char := string(rune)
		if char == "'" && last != "\\" {
			squote++
			if squote == 2 {
				squote = 0
			}
		} else if char == "\"" && last != "\\" {
			dquote++
			if dquote == 2 {
				dquote = 0
			}
		} else if char == " " && squote == 0 && dquote == 0 && last != "\\" {
			args = append(args, current)
			current = ""
		} else {
			if char != "\\" || (last == "\\" && char == "\\") {
				current += char
			}
			last = char
		}
	}
	if len(args) == 0 {
		args = append(args, cmd)
	} else {
		args = append(args, current)
	}
	return args
}

func prompt(c net.Conn) {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "?"
	}
	fmt.Fprintf(c, fmt.Sprintf("[xc: %s]: ", cwd))
}

func lfwd(fwd utils.Forward, s *yamux.Session, c net.Conn) {
	go func() {
		proxy, err := s.Accept()
		if err != nil {
			//log.Println(err)
			return
		}
		fwdCon, err := net.Dial("tcp", fmt.Sprintf("%s:%s", fwd.Addr, fwd.RPort))
		if err != nil {
			//log.Println(err)
			return
		}
		defer fwdCon.Close()
		go utils.CopyIO(fwdCon, proxy)
		go utils.CopyIO(proxy, fwdCon)
	}()
	for {
		select {
		case <-fwd.Quit:
			return
		}
	}
}

// opens the listening socket on the client side
func rfwd(fwd utils.Forward, session *yamux.Session, c net.Conn) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", fwd.LPort))
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		for {
			fwdCon, err := ln.Accept()
			if err != nil && fwdCon != nil {
				defer fwdCon.Close()
				if err != nil {
					//log.Println(err)
				}
				proxy, err := session.Open()
				if err != nil {
					//log.Println(err)
				}
				go utils.CopyIO(fwdCon, proxy)
				go utils.CopyIO(proxy, fwdCon)
			}
		}
	}()
	// Wait for exit signal
	for {
		select {
		case <-fwd.Quit:
			ln.Close()
			return
		}
	}
}
