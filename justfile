[doc('List all available commands')]
help:
  @just -l --list-heading $'Available commands:\n'

[doc('(re)Generate project mocks')]
gen-mocks:
  docker run --rm -v "$PWD":/src -w /src vektra/mockery:3

[doc('Run golangci-lint')]
lint:
  docker run --rm -v $(pwd):/app -v ~/.cache/golangci-lint/:/root/.cache -w /app golangci/golangci-lint:v2.8.0 golangci-lint run --timeout 3m --verbose
