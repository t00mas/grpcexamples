#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

for dir in $(find ${DIR}/proto -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq); do
  files=$(find "${dir}" -name '*.proto')

  protoc -I ${DIR} --go_out=${DIR}/proto/gen/go --go-grpc_out=${DIR}/proto/gen/go ${files}
done