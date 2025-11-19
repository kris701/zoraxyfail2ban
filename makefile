.PHONY: all

all: zoraxyfail2ban_linux_386 zoraxyfail2ban_linux_amd64 zoraxyfail2ban_linux_arm zoraxyfail2ban_linux_arm64 zoraxyfail2ban_linux_mipsle zoraxyfail2ban_linux_riscv64 zoraxyfail2ban_windows_amd64.exe

zoraxyfail2ban_linux_386:
	mkdir -p ./build
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o zoraxyfail2ban_linux_386
	mv zoraxyfail2ban_linux_386 ./build/

zoraxyfail2ban_linux_amd64:
	mkdir -p ./build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o zoraxyfail2ban_linux_amd64
	mv zoraxyfail2ban_linux_amd64 ./build/

zoraxyfail2ban_linux_arm:
	mkdir -p ./build
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o zoraxyfail2ban_linux_arm
	mv zoraxyfail2ban_linux_arm ./build/

zoraxyfail2ban_linux_arm64:
	mkdir -p ./build
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o zoraxyfail2ban_linux_arm64
	mv zoraxyfail2ban_linux_arm64 ./build/

zoraxyfail2ban_linux_mipsle:
	mkdir -p ./build
	CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -o zoraxyfail2ban_linux_mipsle
	mv zoraxyfail2ban_linux_mipsle ./build/

zoraxyfail2ban_linux_riscv64:
	mkdir -p ./build
	CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build -o zoraxyfail2ban_linux_riscv64
	mv zoraxyfail2ban_linux_riscv64 ./build/

zoraxyfail2ban_windows_amd64.exe:
	mkdir -p ./build
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o zoraxyfail2ban_windows_amd64.exe
	mv zoraxyfail2ban_windows_amd64.exe ./build/

.PHONY: all zoraxyfail2ban_linux_386 zoraxyfail2ban_linux_amd64 zoraxyfail2ban_linux_arm zoraxyfail2ban_linux_arm64 zoraxyfail2ban_linux_mipsle zoraxyfail2ban_linux_riscv64 zoraxyfail2ban_windows_amd64.exe
