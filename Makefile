test:
	go test -v --race

cover:
	@rm -rf *.coverprofile
	go test -coverprofile=json-mask.coverprofile -v --race
	gover
	go tool cover -html=json-mask.coverprofile
