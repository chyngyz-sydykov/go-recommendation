#!/bin/bash

EXCLUDED_PACKAGES="infrastructure/config"
INCLUDED_PACKAGES=$(go list ./... | grep -v -E "$EXCLUDED_PACKAGES" | tr '\n' ',' | sed 's/,$//')

mkdir -p coverage
APP_ENV=test go test -coverpkg="$INCLUDED_PACKAGES" -coverprofile=coverage/coverage.out ./tests -v
grep -v "application/application.go" coverage/coverage.out > coverage/filtered_coverage.out