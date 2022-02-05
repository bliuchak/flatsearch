current_dir = $(shell pwd)

# based on https://gist.github.com/prwhite/8168133
help:              ## show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

run:               ## runs flatsearch and all necessary components
	docker-compose up --build flatsearch

stop:              ## stops flatsaerch and all necessary components
	docker-compose down --remove-orphans

migration-status:  ## runs migration status
	docker-compose up --build goose

migration-up:      ## runs migration up
	GOOSE_CMD=up docker-compose up --build goose

migration-down:    ## runs migration down
	GOOSE_CMD=down docker-compose up --build goose

lint: golangci hadolint

golangci:          ## runs golangci-lint
	docker run --rm -v $(current_dir):/app -w /app golangci/golangci-lint:v1.32-alpine golangci-lint run -v

hadolint:          ## runs hadolint
	docker run --rm -i hadolint/hadolint:latest-alpine < Dockerfile
