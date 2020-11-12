GO111MODULE=on

.PHONY: production
production: vue prin
prin:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/prin ./cmd/main.go

vue:
	cd web && npm install && npm run build

.PHONY: dev
dev: prin-dev vue-dev
prin-dev:
	-rm ./bin/prin
	go fmt ./...
	go build -o ./bin/prin ./cmd/main.go
	./bin/prin

vue-dev:
	cd web && npm run serve

.PHONY: clean
clean:
	-rm ./bin/*