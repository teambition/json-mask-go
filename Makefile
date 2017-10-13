test:
	go test -v

cover:
	@rm -rf *.coverprofile
	go test -coverprofile=json-mask.coverprofile -v
	gover
	go tool cover -html=json-mask.coverprofile
