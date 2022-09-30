#
# API building and running locally
# Set of tasks related to API building and running.
#
setup-api:
	go mod download

audit-api:
	go mod tidy

format-api:
	gofmt -w .; \
	golint ./...

run-api:
	. ./scripts/setup_env_vars.sh; \
	go run cmd/api/main.go

doc-api:
	swag init -g ./cmd/api/main.go -o ./docs/api/swagger

#
# API test
# Set of tasks related to API testing locally.
#
build-mocks:
	. ./scripts/build_mocks.sh

test-api:
	. ./scripts/setup_env_vars.test.sh; \
	go test ./internal/... -v -coverprofile=./docs/api/tests/unit/coverage.out && go tool cover -func=./docs/api/tests/integration/coverage.out > ./docs/api/tests/unit/coverage_report.out; \
	go test ./tests/api/... -v -coverprofile=./docs/api/tests/integration/coverage.out && go tool cover -func=./docs/api/tests/integration/coverage.out > ./docs/api/tests/integration/coverage_report.out

#
# APP test container
# Set of tasks related to APP testing container.
#
start-deps:
	docker network create testapp_network; \
	docker build -t postgrestestdb -f ./database/postgres/Dockerfile .; \
	docker run --name postgrestestdb_container --env-file ./database/postgres/.env.test -d -p 5434:5432 -v postgrestestdb-data:/var/lib/postgresql/data --restart on-failure postgrestestdb; \
	docker network connect testapp_network postgrestestdb_container

test-app:
	docker build -t apitest -f ./tests/Dockerfile.test .; \
	docker run --name apitest_container --env-file ./tests/.env.test -d -p 5001:5001 --restart on-failure apitest; \
	docker network connect testapp_network apitest_container; \
	docker exec --env-file ./tests/.env.test apitest_container poetry run pytest; \
	docker network disconnect testapp_network apitest_container; \
	docker stop apitest_container; \
 	docker rm apitest_container; \
 	docker rmi apitest

finish-deps:
	docker network disconnect testapp_network postgrestestdb_container; \
	docker stop postgrestestdb_container; \
	docker rm postgrestestdb_container; \
	docker rmi postgrestestdb; \
	docker network rm testapp_network

#
# APP production container
# Set of tasks related to APP production container starting up and shutting down.
#
startup-app:
	docker-compose up -d --build api

shutdown-app:
	docker-compose down -v --rmi all

# #
# # Set of tasks related to APP container testing
# #
# start-deps:
# 	docker-compose up -d --build postgrestestdb

# finish-deps:
# 	docker-compose rm --force --stop -v postgrestestdb

# test-app:
# 	docker exec --env-file ./.env.test api_container go test ./... -v -coverprofile=./docs/api/tests/unit/coverage.out

# analyze-app:
# 	docker-compose exec api go tool cover -func=./docs/api/tests/unit/coverage.out

# #
# # Set of tasks related to APP deployment.
# #
# init-deploy:
# 	cd deployments/heroku/terraform; \
# 	terraform init

# plan-deploy:
# 	. ./deployments/heroku/scripts/setup_env_vars.sh; \
# 	cd deployments/heroku/terraform; \
# 	terraform plan

# apply-deploy:
# 	. ./deployments/heroku/scripts/build_app.sh; \
# 	. ./deployments/heroku/scripts/setup_env_vars.sh; \
# 	cd deployments/heroku/terraform; \
# 	terraform apply

# destroy-deploy:
# 	. ./deployments/heroku/scripts/destroy_app.sh; \
# 	. ./deployments/heroku/scripts/setup_env_vars.sh; \
# 	cd deployments/heroku/terraform; \
# 	terraform destroy
