project := xlang


.PHONY: install
install:
	glide --version || go get github.com/Masterminds/glide
	GO15VENDOREXPERIMENT=1 glide install
	GO15VENDOREXPERIMENT=1 go build `glide novendor`


.PHONY: test
test: lint
	GO15VENDOREXPERIMENT=1 go test `glide novendor`


.PHONY: lint
lint:
	go list . | xargs golint
	go list . | xargs go vet


.PHONY: xlang
xlang:
	docker-compose run xlang

.PHONY: xlang-fresh
xlang-fresh:
	docker-compose kill
	docker-compose rm -f
	docker-compose pull
	docker-compose build
	docker-compose run xlang


.PHONY: run
run:
	docker-compose build xlang
	docker-compose run xlang


.PHONY: scratch
scratch:
	GO15VENDOREXPERIMENT=1 CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main `GO15VENDOREXPERIMENT=1 glide novendor`
	docker build -f Dockerfile.scratch -t scratch-xlang .
	docker run scratch-xlang

