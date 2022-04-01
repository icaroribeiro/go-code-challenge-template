#
# API-related tasks.
#
setup/api:
	go mod download

audit/api:
	go mod tidy

format/api:
	gofmt -w .;
	golint ./...

check/api:
	go run cmd/api/main.go version

run/api:
	. ./scripts/setup_env_vars.sh;
	go run cmd/api/main.go run

doc/api:
	swag init -g ./cmd/api/main.go -o ./docs/api/swagger

test/api:
	. ./scripts/setup_env_vars.test.sh;
	go test ./... -v -coverprofile=./docs/api/tests/unit/coverage.out

analyze/api:
	go tool cover -func=./docs/api/tests/unit/coverage.out > ./docs/api/tests/unit/coverage_report.out

build/mocks:
	. ./scripts/build_mocks.sh

#
# Container-related tasks.
#
startup/app:
	docker-compose up -d

test/app:
	docker exec api_container go test ./... -v -coverprofile=./docs/api/tests/unit/coverage.out

analyze/app:
	docker exec api_container go tool cover -func=./docs/api/tests/unit/coverage.out

shutdown/app:
	docker-compose down -v --rmi all

#
# Deployment-related tasks.
#
init/deploy:
	cd deployments/heroku/terraform; \
	terraform init

plan/deploy:
	. ./deployments/heroku/scripts/setup_env_vars.sh; \
	cd deployments/heroku/terraform; \
	terraform plan

apply/deploy:
	. ./deployments/heroku/scripts/build_app.sh; \
	. ./deployments/heroku/scripts/setup_env_vars.sh; \
	cd deployments/heroku/terraform; \
	terraform apply

destroy/deploy:
	. ./deployments/heroku/scripts/destroy_app.sh; \
	. ./deployments/heroku/scripts/setup_env_vars.sh; \
	cd deployments/heroku/terraform; \
	terraform destroy
