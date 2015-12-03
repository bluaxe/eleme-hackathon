all:main

main:
	rm -rf main
	@cp -R mem $(GOPATH)/src/
	@cp -R cache $(GOPATH)/src/
	@cp -R server $(GOPATH)/src/
	@cp -R persist $(GOPATH)/src/
	@cp -R common $(GOPATH)/src/
	@cp -R service $(GOPATH)/src/
	@cp -R thirdparty/* $(GOPATH)/src/
	go build main.go

clean:
	rm -rf main
