run:
	source config/local.sh && \
	go run cmd/go-auth/main.go

live:
	source config/local.sh && \
	APP_PORT=9029 \
	gin -i -p 9030 -a 9029 -d cmd/go-auth/ -t ./ run

db-migrate:
	source config/local.sh && \
	go run cmd/db-migrate/main.go

go-test:
	source config/local.sh && \
	go test -v ./...