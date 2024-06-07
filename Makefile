# lamba will build to go project and create the zip file for uploading to AWS.
lambda:
	GOARCH=amd64 GOOS=linux go build -o main main.go
	zip function.zip main

test:
	go test ./...