all:
	go fmt ./...
	go vet ./...
	golint ./...

mime:
	wget https://raw.githubusercontent.com/jshttp/mime-db/master/db.json -O mime.json
