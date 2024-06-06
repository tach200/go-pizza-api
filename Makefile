# lamba will build to go project and create the zip file for uploading to AWS.
lambda:
	GOARCH=amd64 GOOS=linux go build main.go
	zip function.zip main