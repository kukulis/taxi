# web game 

Test multiple files

    go test ./...

## run

    go run cmd/server/main.go

Browser:

    http://localhost:8880

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

## Testing WS client

    go run cmd/wsclient/main.go -type driver_accepts_offer -data '{"passenger_id":"p1"}'


## certificates

Cmd line go to tls directory and run command:

    openssl genrsa -out server.key 2048
    
    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

Example to run without questions prompt:

    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -nodes -subj "/C=LT/ST=Lietuva/L=Kaunas/O=Darbelis/OU=taxi/CN=taxi"

For testing from a phone, the certificate needs a `subjectAltName` matching the address the phone will
actually connect to (your dev machine's LAN IP, e.g. `192.168.1.23`, found via `ip addr` / `ifconfig`) —
without it, mobile browsers reject the certificate outright instead of offering to proceed anyway:

    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -nodes -subj "/C=LT/ST=Lietuva/L=Kaunas/O=Darbelis/OU=taxi/CN=taxi" -addext "subjectAltName=IP:192.168.1.23,IP:127.0.0.1,DNS:localhost"
