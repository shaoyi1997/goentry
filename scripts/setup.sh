#!/usr/bin/env bash

# Source the common.sh script
# shellcheck source=./common.sh
. "$(git rev-parse --show-toplevel || echo ".")/scripts/common.sh"

cd "$PROJECT_DIR"

# If on MacOs with Homebrew available, use it to install golangci-lint
if ! has golangci-lint; then
    if has brew; then
        echo_info "Install golangci-lint for static code analysis (via Homebrew)"
        brew install golangci/tap/golangci-lint
        brew upgrade golangci/tap/golangci-lint
    else
        echo_info "Install golangci-lint for static code analysis (via curl)"
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh |
          sh -s -- -b $(go env GOPATH)/bin v1.39.0

    fi
fi

if ! has goimports; then
    echo_info "Install goimports"
    go get -v -u golang.org/x/tools/cmd/goimports
fi

#if ! has mockgen; then
#    echo_info "Install mockgen"
#    go get -v -u github.com/golang/mock/gomock
#    go install -v -i github.com/golang/mock/mockgen
#fi

echo_info "Config git hooks pash"
git config core.hooksPath "$PROJECT_DIR/scripts/git-hooks"

cd "$WORKING_DIR" || exit 1

go version