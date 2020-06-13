lint:
	@golangci-lint run --deadline=5m

test:
	@echo "scale up zookeeper cluster"
	docker-compose -f "test/docker-compose/docker-compose-zk-cluster.yml" up -d --build
	@sleep 3
	go test -count=1 -v -p 1 $(shell go list ./...)
	@echo "shutdown zookeeper cluster"
	docker-compose -f "test/docker-compose/docker-compose-zk-cluster.yml" down
