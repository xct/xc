BUILD=go build
GENERATE=go generate
OUT_LINUX=xc
OUT_WINDOWS=xc.exe
OUT_WINDOWS_ARM=xc_arm.exe
SRC=xc.go

AUTH_KEY = 

all: clean linux64 windows64 #windowsArm64    

linux64:
	GOOS=linux GOARCH=amd64 ${GENERATE}
	GOOS=linux GOARCH=amd64 ${BUILD} -o ${OUT_LINUX} ${SRC}

windows64:
	GOOS=linux GOARCH=amd64 ${GENERATE}
	GOOS=windows GOARCH=amd64 ${BUILD} -o ${OUT_WINDOWS} ${SRC}

windowsArm64:
	GOOS=linux GOARCH=arm ${GENERATE}
	GOOS=windows GOARCH=arm ${BUILD} -o ${OUT_WINDOWS_ARM} ${SRC}

clean:
	mkdir -p files/keys
	yes 'y' | ssh-keygen -t ed25519 -f files/keys/key -q -N ""
	yes 'y' | ssh-keygen -f host_dsa -N '' -t dsa -f files/keys/host_dsa -q -N ""
	yes 'y' | ssh-keygen -f host_rsa -N '' -t rsa -f files/keys/host_rsa -q -N ""
	rm -f ${OUT_LINUX} ${OUT_WINDOWS} ${OUT_WINDOWS_ARM} shell/keys.go meter/sc.go