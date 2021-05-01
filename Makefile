BUILD=GO111MODULE=off go build -ldflags="-s -w" -buildmode=pie -trimpath 
GENERATE=GO111MODULE=off go generate
OUT_LINUX=xc
OUT_WINDOWS=xc.exe
SRC=xc.go
LINUX_COMPRESS=upx xc -o xcc; rm xc && mv xcc xc
XC_POSTGEN=cp xc.bak xc.go; rm xc.bak
XC_WIN_POSTGEN=cp client/client_windows.bak client/client_windows.go; rm client/client_windows.bak

all: clean generate linux64 windows64 postgen  

generate:
	${GENERATE}	

postgen:
	${XC_POSTGEN}
	${XC_WIN_POSTGEN}

linux64:
	GOOS=linux GOARCH=amd64 ${BUILD} -o ${OUT_LINUX} ${SRC}
	${LINUX_COMPRESS}	

windows64:	
	GOOS=windows GOARCH=amd64 ${BUILD} -o ${OUT_WINDOWS} ${SRC}	

clean:
	- rm files/keys/host*
	- rm files/keys/key* 
	mkdir -p files/keys
	yes 'y' | ssh-keygen -t ed25519 -f files/keys/key -q -N ""
	yes 'y' | ssh-keygen -f host_dsa -N '' -t dsa -f files/keys/host_dsa -q -N ""
	yes 'y' | ssh-keygen -f host_rsa -N '' -t rsa -f files/keys/host_rsa -q -N ""
	rm -f ${OUT_LINUX} ${OUT_WINDOWS} shell/keys.go meter/sc.go