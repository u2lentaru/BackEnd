  
.PHONY: clean
clean:
	$(eval IMAGES := $(shell docker-compose -f docker-compose.yml images | awk '{ print $$2" "$$4 }' | grep 'geekbrainsqueues_' | awk '{ print $$2 }'))
	docker-compose -f docker-compose.yml rm --stop --force -v
	docker rmi $(IMAGES) --force 2>/dev/null; true

.PHONY: env
env: clean
	docker-compose -f docker-compose.yml up -d

.PHONY: server
server:
	CGO_ENABLED=0 go build -v -mod=vendor -ldflags "-s" -o ./bin/api ./api
	CGO_ENABLED=0 go build -v -mod=vendor -ldflags "-s" -o ./bin/process ./process

.PHONY: test
test:
	go test -v -race ./...