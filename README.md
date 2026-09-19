# web game 

Test multiple files

    go test ./...

## run

    go run cmd/server/main.go

Browser:

    http://localhost:8080

## swagger

    http://localhost:8080/swagger/index.html

### Regenerate swagger docs

Install swag if not already installed:

    go install github.com/swaggo/swag/cmd/swag@latest

Generate docs (run from project root):

    swag init -g cmd/server/main.go

In my case

    ~/go/bin/swag init -g cmd/server/main.go

# pitfalls

The Go json.Unmarshal works case-insensitive, while JS decoding json is case-sensitive.

# TESTS for frontend

Assume you have node installed in to your machine. The required version is 18 or later.

    node --test test/js/**/*.test.js

