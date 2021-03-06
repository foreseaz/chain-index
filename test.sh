#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

DEFAULT_MIGRATION_FOLDER="$(pwd)/migrations"
export TEST_POSTGRES_MIGRATIONS_FOLDER=${TEST_POSTGRES_MIGRATIONS_FOLDER:-${DEFAULT_MIGRATION_FOLDER}}
export TEST_POSTGRES_USERNAME=${TEST_POSTGRES_USERNAME:-postgres}
export TEST_POSTGRES_PASSWORD=${TEST_POSTGRES_PASSWORD:-postgres}
export TEST_POSTGRES_HOST=${TEST_POSTGRES_HOST:-localhost}
export TEST_POSTGRES_PORT=${TEST_POSTGRES_PORT:-5433}
export TEST_POSTGRES_DATABASE=${TEST_POSTGRES_DATABASE:-postgres_test}
export TEST_POSTGRES_SCHEMA=${TEST_POSTGRES_SCHEMA:-public}

GOPATH=$(go env GOPATH)
DOCKER_COMPOSE_PROJECT="chain_indexing_test"

# main
RET_VALUE=0

cd "$(dirname "${BASH_SOURCE[0]}")"

POSTGRES_DRIVER_URL="postgres://${TEST_POSTGRES_USERNAME}:${TEST_POSTGRES_PASSWORD}@${TEST_POSTGRES_HOST}:${TEST_POSTGRES_PORT}/${TEST_POSTGRES_DATABASE}?options=--search_path%3d=${TEST_POSTGRES_SCHEMA}"

echoerr() { echo "$@" 1>&2; }

check_ginkgo() {
    set +e
    command -v "${GOPATH}/bin/ginkgo" -r > /dev/null
    RET_VALUE=$?
    set -e
}

trap handle_ctrlc INT
handle_ctrlc() {
    echo "Captured Ctrl-C. Stopping test environment ..."
    teardown

    echo
    echo "Test Interrupted"
    exit 1
}

teardown() {
    if [[ "${TEST_DB}" == 1 ]]; then
        docker-compose -f ./docker/docker-compose.test.yml -p "${DOCKER_COMPOSE_PROJECT}" down
    fi
}

setup() {
    if [[ "${TEST_DB}" == 1 ]]; then
        docker-compose -f ./docker/docker-compose.test.yml -p "${DOCKER_COMPOSE_PROJECT}" up -d
        wait_postgres_ready
    fi
}

wait_postgres_ready() {
    CONTAINER_ID=$(docker ps -aqf "name=${DOCKER_COMPOSE_PROJECT}_postgres")
    
    while :
    do
        is_postgres_ready "${CONTAINER_ID}"
        if [[ "${RET_VALUE}" == 0 ]]; then
            echo "Postgres is ready to accept connections"
            return
        fi
        echo "Waiting for Postgres container(${CONTAINER_ID}) to be ready ..."
        sleep 1
    done
}

is_postgres_ready() {
    set +e
    docker exec -it "$1" pg_isready > /dev/null
    RET_VALUE=$?
    set -e
}

run_test() {
    set +e
    "${GOPATH}/bin/ginkgo" -v -r --trace
    RET_VALUE=$?
    set -e
}

run_test_watch() {
    set +e
    "${GOPATH}/bin/ginkgo" watch -v -r
    RET_VALUE=$?
    set -e
}

INSTALL_DEPENDENCY=0
WATCH_MODE=0
TEST_DB=1
while [[ $# > 0 ]]; do
    case "$1" in
        --install-dependency)
            INSTALL_DEPENDENCY=1
            shift 1
        ;;
        --no-db)
            TEST_DB=0
            shift 1
        ;;
        --no-ssl)
            POSTGRES_DRIVER_URL="${POSTGRES_DRIVER_URL}&sslmode=disable"
            shift 1
        ;;
        --watch)
            WATCH_MODE=1
            shift 1
        ;;
        *)
            echoerr "Unknown argument: $1"
            usage
        ;;
    esac
done

check_ginkgo
if [[ "${RET_VALUE}" != 0 ]]; then
    if [[ "${INSTALL_DEPENDENCY}" = 0 ]]; then
        echoerr "Cannot find ginkgo. Provide \`--install-dependency\` to attempt to install it."
        exit 1
    fi
    
    go get -u github.com/onsi/ginkgo/ginkgo
    check_ginkgo
    if [[ "${RET_VALUE}" != 0 ]]; then
        echoerr "Failed to install ginkgo. Please install it manually."
        exit 1
    fi
fi

echo "Stopping any previous test envionrment ..."
teardown
setup

export TEST_DATABASE_URL="${POSTGRES_DRIVER_URL}"
echo
if [[ "${WATCH_MODE}" == 1 ]]; then
    run_test_watch
else
    run_test
fi
TEST_RESULT="${RET_VALUE}"
echo

echo "Test finished. Stopping test envionrment ..."
teardown

echo
if [[ "${TEST_RESULT}" != 0 ]]; then
    echo "Test Failed. Check the trace above"
else
    echo "Test Passed"
fi
exit "${TEST_RESULT}"
