lambda:
	GOARCH=amd64 GOOS=linux go build main.go
	zip function.zip main