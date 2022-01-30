#!/bin/bash

# Database Repository
# ----------------------------------------------------------------------------------------------------
# Generate a mock object related to auth datastore repository.
# AUTH_REPOSITORY_PATH="internal/core/ports/datastore/repository/auth"
# AUTH_REPOSITORY_MOCK_PATH="internal/core/ports/datastore/repositorymock/auth"
# mockery -dir "$AUTH_REPOSITORY_PATH" -name IRepository -outpkg auth -structname RepositoryMock -output "$AUTH_REPOSITORY_MOCK_PATH" -filename repository_mock.go

# Generate a mock object related to login datastore repository.
# LOGIN_REPOSITORY_PATH="internal/core/ports/datastore/repository/login"
# LOGIN_REPOSITORY_MOCK_PATH="internal/core/ports/datastore/repositorymock/login"
# mockery -dir "$LOGIN_REPOSITORY_PATH" -name IRepository -outpkg login -structname RepositoryMock -output "$LOGIN_REPOSITORY_MOCK_PATH" -filename repository_mock.go

# Generate a mock object related to user datastore repository.
USER_REPOSITORY_PATH="internal/core/ports/infrastructure/persistence/datastore/repository/user"
USER_REPOSITORY_MOCK_PATH="internal/core/ports/datastore/repositorymock/user"
mockery -dir "$USER_REPOSITORY_PATH" -name IRepository -outpkg user -structname RepositoryMock -output "$USER_REPOSITORY_MOCK_PATH" -filename repository_mock.go

# Service
# ----------------------------------------------------------------------------------------------------
# Generate a mock object related to auth service.
# AUTH_SERVICE_PATH="internal/core/ports/application/service/auth"
# AUTH_SERVICE_MOCK_PATH="internal/core/ports/application/servicemock/auth"
# mockery -dir "$AUTH_SERVICE_PATH" -name IService -outpkg auth -structname ServiceMock -output "$AUTH_SERVICE_MOCK_PATH" -filename service_mock.go

# Generate a mock object related to healthcheck service.
HEALTHCHECK_SERVICE_PATH="internal/core/ports/application/service/healthcheck"
HEALTHCHECK_SERVICE_MOCK_PATH="internal/core/ports/application/servicemock/healthcheck"
mockery -dir "$HEALTHCHECK_SERVICE_PATH" -name IService -outpkg healthcheck -structname ServiceMock -output "$HEALTHCHECK_SERVICE_MOCK_PATH" -filename service_mock.go

# Generate a mock object related to user service.
USER_SERVICE_PATH="internal/core/ports/application/service/user"
USER_SERVICE_MOCK_PATH="internal/core/ports/application/servicemock/user"
mockery -dir "$USER_SERVICE_PATH" -name IService -outpkg user -structname ServiceMock -output "$USER_SERVICE_MOCK_PATH" -filename service_mock.go

# Application
# ----------------------------------------------------------------------------------------------------
# Generate a mock object related to validator.
# VALIDATOR_PATH="pkg/validator"
# VALIDATOR_MOCK_PATH="pkg/validatormock"
# mockery -dir "$VALIDATOR_PATH" -name IValidator -outpkg validatormock -structname ValidatorMock -output "$VALIDATOR_MOCK_PATH" -filename validator_mock.go

# Infrastructure
# ----------------------------------------------------------------------------------------------------
# Generate a mock object related to authentication.
# AUTH_PATH="internal/infrastructure/auth"
# AUTH_MOCK_PATH="internal/infrastructure/authmock"
# mockery -dir "$AUTH_PATH" -name IAuth -outpkg authmock -structname AuthMock -output "$AUTH_MOCK_PATH" -filename auth_mock.go

# Generate a mock object related to security.
# SECURITY_PATH="internal/infrastructure/security"
# SECURITY_MOCK_PATH="internal/infrastructure/securitymock"
# mockery -dir "$SECURITY_PATH" -name ISecurity -outpkg securitymock -structname SecurityMock -output "$SECURITY_MOCK_PATH" -filename security_mock.go
