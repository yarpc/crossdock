project := crossdock

export GO15VENDOREXPERIMENT=1


.PHONY: install
install:
	glide --version || go get github.com/Masterminds/glide
	glide install
	go build `glide novendor`


.PHONY: test
test:
	go test `glide novendor`
	./tests/succeed.sh
	./tests/fail.sh


.PHONY: clean
clean:
	docker-compose kill
	docker-compose rm -f


.PHONY: crossdock
crossdock:
	docker-compose run crossdock


.PHONY: crossdock-fresh
crossdock-fresh: clean
	docker-compose pull
	docker-compose build
	docker-compose run crossdock


.PHONY: run
run:
	docker-compose build crossdock
	docker-compose run crossdock


.PHONY: scratch
scratch:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
	docker build -f Dockerfile.scratch -t scratch-crossdock .
	docker run scratch-crossdock

