COMMIT_HASH = $(shell git log -n 1 --pretty=format:"%H")

build:
	@DOCKER_BUILDKIT=1 docker build \
		--build-arg GITHUB_TOKEN=$(GITHUB_TOKEN) \
		--build-arg COMMIT_HASH=$(COMMIT_HASH) \
		--target release \
		-f Dockerfile \
		-t $(IMAGE_TAG) .

run:
	go run ./cmd/server.go