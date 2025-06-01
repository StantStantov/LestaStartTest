set dotenv-load := true
set dotenv-required := true

export GOCACHE := `go env GOCACHE`
export GOMODCACHE := `go env GOMODCACHE`

remote-context := "lesta-start"
dev-compose := "deployments/dev/compose.dev.yaml"
test-compose := "deployments/test/compose.test.yaml"
prod-compose := "deployments/prod/compose.prod.yaml"

default:
  @just --list

build-dev:
  COMPOSE_BAKE=true docker compose -f {{dev-compose}} build

build-test:
  COMPOSE_BAKE=true docker compose -f {{test-compose}} build

build-prod:
  COMPOSE_BAKE=true docker compose -f {{prod-compose}} build

test:
  python scripts/run-tests.py

up-dev:
  docker compose -f {{dev-compose}} up

up-prod:
  docker compose -f {{prod-compose}} up

deploy:
  docker --context {{remote-context}} stack deploy -c {{prod-compose}} lesta-start

templ:
  #!/usr/bin/env sh
  go tool templ generate --path "web/html"

  for directory in `find web/html -type d`; do
    if `find "${directory}" -maxdepth 1 -name "*_templ.go" | read v`; then
      copy_to=${directory/web\/html/internal/api/web/views};
      if ! `find ${copy_to} -type d | read v`; then
        mkdir ${copy_to};
      fi;
      mv "${directory}"/*.go "${copy_to}"/;
    fi;
  done
