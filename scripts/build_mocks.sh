#!/bin/bash

#
# Internal
#
# Generate a mock object related to healthcheck's service.
HEALTHCHECK_SERVICE_PATH="internal/core/ports/application/service/healthcheck"
MOCK_HEALTHCHECK_SERVICE_PATH="internal/core/ports/application/mockservice/healthcheck"
mockery --dir "$HEALTHCHECK_SERVICE_PATH" --name IService --outpkg healthcheck --structname Service --output "$MOCK_HEALTHCHECK_SERVICE_PATH" --filename mock_service.go

#
# PKG
#
# Generate a mock object related to auth.
AUTH_PATH="pkg/auth"
MOCK_AUTH_PATH="tests/mocks/pkg/mockauth"
mockery --dir "$AUTH_PATH" --name IAuth --outpkg mockauth --structname Auth --output "$MOCK_AUTH_PATH" --filename mock_auth.go

# Generate a mock object related to security.
SECURITY_PATH="pkg/security"
MOCK_SECURITY_PATH="tests/mocks/pkg/mocksecurity"
mockery --dir "$SECURITY_PATH" --name ISecurity --outpkg mocksecurity --structname Security --output "$MOCK_SECURITY_PATH" --filename mock_security.go

# Generate a mock object related to validator.
VALIDATOR_PATH="pkg/validator"
MOCK_VALIDATOR_PATH="tests/mocks/pkg/mockvalidator"
mockery --dir "$VALIDATOR_PATH" --name IValidator --outpkg mockvalidator --structname Validator --output "$MOCK_VALIDATOR_PATH" --filename mock_validator.go