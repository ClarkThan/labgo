# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
IMAGE_TAG := "v$(shell date +%Y%m%d.%H%M%S)"

# All target
all: test build

run:
	CGO_ENABLED=0 $(GOCMD) run main.go

# Run tests
test:
	$(GOTEST) -v ./...

# set go env
setup:
	@echo ">> Setup environments"
	$(GOCMD) env -w GOPRIVATE=gitlab.meiqia.com
	$(GOCMD) env -w GOPROXY=https://goproxy.cn,direct
	$(GOCMD) env -w GOFLAGS="-buildvcs=false"

# Build the project
build: setup
	@echo ">> build automation"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)

docker-login:
	docker login harbor.meiqia.com

# build and push docker image
qa-image: build docker-login
	@echo ">> build docker image"
	docker build --platform linux/amd64 -t automation:${IMAGE_TAG} -f Dockerfile.qa .
	docker tag automation:${IMAGE_TAG} harbor.meiqia.com/backend/automation:${IMAGE_TAG}
	docker push harbor.meiqia.com/backend/automation:${IMAGE_TAG}

# Clean the build directory
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Lint the code
lint:
	@echo ">> linting code"
	golangci-lint run --max-same-issues 100 -v ./...

.PHONY: all run test setup build docker-login qa-image clean lint