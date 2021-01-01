package client

import (
	"bufio"
	"fmt"
	"net"
	"os/user"
	"strconv"

	"../plugins"
	"../shell"
	"github.com/hashicorp/yamux"
)

var signalSession *yamux.Session

// Run runs the mainloop of the shell
func Run(s *yamux.Session, c net.Conn) {
	defer c.Close()

	// open 2nd channel
	signalSession, err := s.Accept()
	if err != nil {
		fmt.Println(err)
	}
	signalScanner := bufio.NewScanner(signalSession)

	scanner := bufio.NewScanner(c)
	usr, _ := user.Current()
	homedir := usr.HomeDir

	// init
	plugins.Init(c)
	prompt(c)

	go func() {
		for signalScanner.Scan() {
			command := signalScanner.Text()
			//fmt.Printf("Command %s\n", command)
			switch command {
			case "!sigint":
				if shellPipeReader != nil && shellPipeWriter != nil && shellCmd != nil {
					//log.Println("Quitting shell (Ctrl+C)")
					shellCmd.Process.Kill()
					shellPipeReader.Close()
					shellPipeWriter.Close()
					shellCmd = nil
				}
			default:
				// pass
			}
		}
	}()

	for scanner.Scan() {
		command := scanner.Text()
		if len(command) > 1 {
			argv := splitArgs(command)
			// we only handle commands here that do something on the client side
			// commands that are shared between different os
			help := usage()
			help += "\n"
			help += " !ssh <port>\n"
			help += "   * starts sshd with the configured keys on the specified port\n"
			handled := handleSharedCommand(s, c, argv, help, homedir)
			// os specific commands
			if !handled {
				switch argv[0] {
				case "!ssh":
					if len(argv) == 2 {
						port, err := strconv.Atoi(argv[1])
						if err == nil {
							shell.StartSSHServer(port, c)
						} else {
							fmt.Println(err)
						}
					} else {
						c.Write([]byte("Usage: !ssh <port>\n"))
					}
					prompt(c)
				default:
					shell.Exec(command, c)
					prompt(c)
				}
			}
		} else {
			prompt(c)
		}
	}
}
