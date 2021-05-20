#!/usr/bin/env bash

golangci-lint run

EXIT_CODE=$?
exit $EXIT_CODE
