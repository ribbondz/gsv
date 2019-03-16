set GOARCH=amd64
set GOOS=linux
go build -o dist/gsv-linux-amd64

set GOARCH=amd64
set GOOS=darwin
go build -o dist/gsv-darwin-amd64

set GOARCH=amd64
set GOOS=netbsd
go build -o dist/gsv-netbsd-amd64

set GOARCH=amd64
set GOOS=freebsd
go build -o dist/gsv-freebsd-amd64

set GOARCH=amd64
set GOOS=windows
go build -o dist/gsv.exe

REM cmd /k
REM cmd /k
