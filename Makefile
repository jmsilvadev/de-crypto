test:
	go test -race -v ./...

test-cov:
	go test -race -cover ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

lint:
	go vet ./...

clean:
	go clean -testcache

up:
	docker compose up -d

up-build:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f

kafka-consumer:
	docker exec -it kafka /opt/kafka/bin/kafka-console-consumer.sh \
		--topic de-crypto-events \
		--bootstrap-server kafka:9092 \
		--from-beginning

