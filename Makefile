tests:
	@go test -v -cover -race github.com/brian-nunez/bconfig/... github.com/brian-nunez/bconfig/drivers/file/... github.com/brian-nunez/bconfig/drivers/env/... -coverprofile=coverage.out
	@go tool cover -func=coverage.out
	
coverage: tests
	@go tool cover -html=coverage.out -o coverage.html
	@open coverage.html

clean:
	@rm -f coverage.out coverage.html
