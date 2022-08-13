lambda:
	GOARCH=amd64 GOOS=linux go build main.go -ldflags="-s -w"
	zip function.zip main