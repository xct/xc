package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"time"

	"../utils"
	"github.com/hashicorp/yamux"
)

var gc net.Conn
var gs *yamux.Session

type augReader struct {
	innerReader io.Reader
	augFunc     func([]byte) []byte
}

type augWriter struct {
	innerWriter io.Writer
	augFunc     func([]byte) []byte
}

func (r *augReader) Read(buf []byte) (int, error) {
	tmpBuf := make([]byte, len(buf))
	n, err := r.innerReader.Read(tmpBuf)
	copy(buf[:n], r.augFunc(tmpBuf[:n]))
	return n, err
}

func (w *augWriter) Write(buf []byte) (int, error) {
	return w.innerWriter.Write(w.augFunc(buf))
}

func sendReader(r io.Reader) io.Reader {
	return &augReader{innerReader: r, augFunc: handleCmd}
}

func recvWriter(w io.Writer) io.Writer {
	return &augWriter{innerWriter: w, augFunc: handleCmd}
}

var (
	session *yamux.Session
)

var sigChan = make(chan os.Signal, 1)
var activeForwards []utils.Forward
var cmdSession *yamux.Session
var assemblies = map[string]bool{}

// opens the listening socket on the server side
func lfwd(fwd utils.Forward) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", fwd.LPort))
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		for {
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
	// Wait for exit signal
	for {
		select {
		case <-fwd.Quit:
			fwd.Active = false
			ln.Close()
			return
		}
	}
}

// connects to the listening port on the client side
func rfwd(fwd utils.Forward, s *yamux.Session, c net.Conn) {
	go func() {
		for {
			// accept the virtual connection initiated by the client
			proxy, err := s.Accept()
			if err != nil {
				log.Println(err)
				return
			}
			// send data to the redirect target (could be into the void..)
			fwdCon, err := net.Dial("tcp", fmt.Sprintf("%s:%s", fwd.Addr, fwd.RPort))
			if err != nil {
				log.Println(err)
				break
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

func quit() {
	time.Sleep(800 * time.Millisecond)
	fmt.Println("Bye!")
	os.Exit(0)
}

func handleCmd(buf []byte) []byte {
	cmd := strings.TrimSuffix(string(buf), "\r\n")
	cmd = strings.TrimSuffix(cmd, "\n")
	argv := strings.Split(cmd, " ")
	switch argv[0] {
	case "!exit":
		// defer exit so we can sent it to the client aswell
		go quit()
	case "!download":
		if len(argv) == 3 {
			dst := argv[2]
			go utils.DownloadListen(dst, session)
		}
	case "!lfwd":
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
				go lfwd(fwd)
				activeForwards = append(activeForwards, fwd)
			} else {
				log.Printf("Local Port %s already in use.\n", lport)
			}
		}
	case "!rfwd":
		if len(argv) == 4 {
			lport := argv[1]
			raddr := argv[2]
			rport := argv[3]
			fwd := utils.Forward{lport, rport, raddr, make(chan bool), false, true}
			go rfwd(fwd, gs, gc)
			activeForwards = append(activeForwards, fwd)
		}
	case "!rmfwd":
		if len(argv) == 2 {
			index, _ := strconv.Atoi(argv[1])
			forward := activeForwards[index]
			forward.Quit <- true
			activeForwards = append(activeForwards[:index], activeForwards[index+1:]...)
		}
	case "!vulns":
		fmt.Println("Be patient - this can take a few minutes..")
	case "!upload":
		if len(argv) != 3 {
			return buf
		}
		src := argv[1]
		go utils.UploadListen(src, session)
	case "!net":
		// same as upload for the server side, hosts the .NET assembly we want to execute
		if len(argv) < 2 {
			return buf
		}
		src := argv[1]
		// only need to upload if its a new path (downside: if you change a binary and a known path it won't get reuploaded, in this case !restart)
		if _, ok := assemblies[src]; !ok {
			assemblies[src] = true
			go utils.UploadListen(src, session)
		}
	case "!debug":
		fmt.Printf("Active Goroutines: %d\n", runtime.NumGoroutine())
	}
	return buf
}

// Run runs the main server loop
func Run(s *yamux.Session, c net.Conn) {
	gc = c
	gs = s
	session = s
	defer c.Close()
	//fmt.Printf("[xc]:")

	// open 2nd session for signals ("virtual connection")
	cmdSession, err := session.Open()
	if err != nil {
		log.Println(err)
	}

	sr := sendReader(os.Stdin)  // intercepts input that is given on stdin and then send to the network
	rw := recvWriter(os.Stdout) // intercepts output that is to received from network and then send to stdout

	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for {
			<-sigChan
			io.WriteString(cmdSession, "!sigint\n")
		}
	}()
	go io.Copy(c, sr)
	io.Copy(rw, c)
}
