release:
	-@mkdir release 2> /dev/null || true
	GOOS=linux GOARCH=amd64 go generate ./...
	GOOS=linux GOARCH=amd64 go build -o release/server ./cmd/server
	GOOS=linux GOARCH=amd64 go build -o release/migrate ./cmd/migrate

.PHONY: release

deploy: release
	ssh covve@128.199.217.21 "sudo systemctl stop special-needs"
	scp release/server covve@128.199.217.21:~/app/server
	scp release/migrate covve@128.199.217.21:~/app/migrate
	ssh covve@128.199.217.21 "sudo systemctl start special-needs"

.PHONY: deploy

dev:
	go generate ./...
	go run cmd/server/main.go

.PHONY: dev

migrate:
	go generate ./...
	go run cmd/migrate/main.go $(ARGS)

.PHONY: migrate

seed:
	go generate ./...
	go run cmd/seed/main.go $(ARGS)

.PHONY: seed
