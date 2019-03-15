export GOARCH=amd64 GOOS=linux
go env
go build -o "dist/gsv-linux-amd64"


export GOARCH=amd64 GOOS=darwin
go env
go build -o "dist/gsv-darwin-amd64"

export GOARCH=amd64 GOOS=netbsd
go build -o "dist/gsv-netbsd-amd64"

export GOARCH=amd64 GOOS=freebsd
go build -o "dist/gsv-freebsd-amd64"

export GOARCH=amd64 GOOS=windows
go build -o "dist/gsv-windows-amd64.exe"
