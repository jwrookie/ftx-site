#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")
PROJECT_ROOT="$SCRIPT_ROOT/.."

pushd "$PROJECT_ROOT"/dao
mockgen -source=./lucky_draw.go -package=mock > ./mock/lucky_draw.go -self_package github.com/foxdex/ftx-site/dao/mock
popd
