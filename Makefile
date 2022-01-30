# API-related tasks.
setup/api:
	go mod download

audit/api:
	go mod tidy

format/api:
	gofmt -w .
	golint ./...

check/api:
	go run cmd/api/main.go version

run/api:
	. ./scripts/setup_env.sh; \
	go run cmd/api/main.go run

doc/api:
	swag init -g ./cmd/api/main.go -o ./docs/api/swagger

test/api:
	. ./scripts/setup_env.test.sh; \
	go test ./... -v -coverprofile=./docs/api/tests/unit/coverage_report.out

analyze/api:
	go tool cover -func=./docs/api/tests/unit/coverage_report.out

build/mocks:
	. ./scripts/build_mocks.sh

# Container-related tasks.
startup/docker:
	docker-compose --env-file ./.env up -d

test/docker:
	docker exec --env-file ./.env.test api_container go test ./... -v -coverprofile=./docs/tests/api/coverage.out

analyze/docker:
	docker-compose exec api go tool cover -func=./docs/tests/api/coverage.out

shutdown/docker:
	docker-compose down -v --rmi all

# Deployment-related tasks.
init/deploy:
	cd deployments/heroku/terraform; \
	terraform init

plan/deploy:
	. ./deployments/heroku/scripts/setup_env.sh; \
	cd deployments/heroku/terraform; \
	terraform plan

apply/deploy:
	. ./deployments/heroku/scripts/copy_code.sh; \
	. ./deployments/heroku/scripts/setup_env.sh; \
	cd deployments/heroku/terraform; \
	terraform apply

destroy/deploy:
	. ./deployments/heroku/scripts/delete_code.sh; \
	. ./deployments/heroku/scripts/setup_env.sh; \
	cd deployments/heroku/terraform; \
	terraform destroy
