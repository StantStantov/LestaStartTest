set dotenv-load := true
set dotenv-required := false

export COMPOSE_BAKE := "true"
export GOCACHE := `go env GOCACHE`
export GOMODCACHE := `go env GOMODCACHE`

remote-context := "lesta-start"
dev-compose := "deployments/dev/compose.dev.yaml"
test-compose := "deployments/test/compose.test.yaml"
prod-compose := "deployments/prod/compose.prod.yaml"

default:
  @just --list

build-dev:
  docker compose -f {{dev-compose}} build

build-test:
  docker compose -f {{test-compose}} build

build-prod:
  docker compose -f {{prod-compose}} build

test-unit:
  docker compose -f {{test-compose}} \
    run --rm --no-deps app go test -v -count=1 -parallel=16 -tags=unit ./internal/...

test-integration:
  -docker compose -f {{test-compose}} \
    run --rm app go test -v -count=1 -parallel=16 -tags=integration ./internal/...
  docker compose -f {{test-compose}} stop db

up-dev:
  docker compose -f {{dev-compose}} up

up-prod:
  docker compose -f {{prod-compose}} up

deploy:
  docker --context {{remote-context}} stack deploy -c {{prod-compose}} lesta-start
