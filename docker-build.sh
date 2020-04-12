GOOS=linux
go build -ldflags "-s -w" -o gstp-linux
docker build -t gstp:latest .