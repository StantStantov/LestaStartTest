run:
  ./tmp/app

test:
  go test ./internal/... -count=1

build:
  go build -C ./cmd -o ../tmp/app

templ:
  #!/usr/bin/env sh
  go tool templ generate --path "web/html"

  for directory in `find web/html -type d`; do
    if `find "${directory}" -maxdepth 1 -name "*_templ.go" | read v`; then
      copy_to=${directory/web\/html/internal/views};
      if ! `find ${copy_to} -type d | read v`; then
        mkdir ${copy_to};
      fi;
      mv "${directory}"/*.go "${copy_to}"/;
    fi;
  done
