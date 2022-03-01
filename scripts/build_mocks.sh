#!/bin/bash

# PKG
# Generate a mock object related to auth.
AUTH_PATH="pkg/auth"
MOCK_AUTH_PATH="tests/pkg/mockauth"
mockery -dir "$AUTH_PATH" -name IAuth -outpkg mockauth -structname MockAuth -output "$MOCK_AUTH_PATH" -filename mock_auth.go

# Generate a mock object related to security.
SECURITY_PATH="pkg/security"
MOCK_SECURITY_PATH="tests/pkg/mocksecurity"
mockery -dir "$SECURITY_PATH" -name ISecurity -outpkg mocksecurity -structname MockSecurity -output "$MOCK_SECURITY_PATH" -filename mock_security.go

# Generate a mock object related to validator.
VALIDATOR_PATH="pkg/validator"
MOCK_VALIDATOR_PATH="tests/pkg/mockvalidator"
mockery -dir "$VALIDATOR_PATH" -name IValidator -outpkg mockvalidator -structname MockValidator -output "$MOCK_VALIDATOR_PATH" -filename mock_validator.go