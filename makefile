mac:
	go env -w GOARCH=amd64
	go env -w GOOS=darwin
	go env -w CGO_ENABLED=0
	go env -w GO111MODULE=on
	go env -w GOPROXY=https://goproxy.cn,direct
	go mod  tidy

linux:
	go env -w GOARCH=amd64
	go env -w GOOS=linux
	go env -w CGO_ENABLED=0
	go env -w GO111MODULE=on
	go env -w GOPROXY=https://goproxy.cn,direct
	go mod  tidy

conf-local:
	mv utils/conf/conf.yaml utils/conf/conf.yaml.bak
	mv utils/conf/conf.yaml.local utils/conf/conf.yaml

conf-rollback:
	mv utils/conf/conf.yaml utils/conf/conf.yaml.local
	mv utils/conf/conf.yaml.bak utils/conf/conf.yaml

build-file:
	go build -o chat -ldflags "-w -s"  -trimpath main.go

build-mac: mac conf-local build-file conf-rollback

build-linux: linux conf-local build-file conf-rollback

run-mac: mac
	go run main.go

update:
	go mody tidy
