# Echo-go
Simple echo TCP server written in [Go](https://go.dev).

## How to use
Run the server (listens on port 8000 by default):
```bash
go run .
```

Connect to server with netcat or telnet:
```bash
telnet 127.0.0.1 8000
# Or
nc 127.0.0.1 8000
```

