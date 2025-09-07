# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AZenv is a lightweight Go web server that displays environment variables and HTTP request information. The project is a single-file Go application (`main.go`) that creates an HTTP/HTTPS server to display client and request information in a format similar to PHP's `$_SERVER` variable.

## Architecture

- **Single binary application**: The entire server is contained in `main.go`
- **Handler function**: `azenvHandler()` processes requests to `/azenv` and generates HTML responses
- **Certificate management**: Built-in functionality to generate self-signed TLS certificates
- **Dual server support**: Can run both HTTP and HTTPS simultaneously

## Development Commands

### Building
```bash
# Build the binary
go build -o azenv

# Build for production (optimized)
go build -ldflags="-s -w" -o azenv

# Cross-compilation example (Linux AMD64)
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o azenv-linux-amd64
```

### Running
```bash
# Run HTTP server on default port (8080)
go run main.go

# Run with custom HTTP port
go run main.go -p 3000

# Run with HTTPS enabled and auto-generate certificate
go run main.go -ssl -gen-cert

# Run both HTTP and HTTPS on custom ports
go run main.go -p 3000 -sp 4430 -ssl -gen-cert

# Run with Let's Encrypt on custom HTTPS port (requires domain)
go run main.go -ssl -lets-encrypt -domain example.com -sp 8443

# Run with Let's Encrypt and custom cache directory
go run main.go -ssl -lets-encrypt -domain example.com -cache-dir /path/to/certs
```

### Testing
```bash
# Run Go tests (if any exist)
go test

# Format code
go fmt

# Vet code for issues
go vet

# Run the application and test endpoint
curl http://localhost:8080/azenv
```

### Dependencies
- Go 1.23+ required (specified in `go.mod`)
- `golang.org/x/crypto` - for Let's Encrypt autocert functionality

## Command Line Flags

- `-p`: HTTP port (default: 8080)
- `-sp`: HTTPS port (default: 8443) 
- `-ssl`: Enable HTTPS server
- `-gen-cert`: Auto-generate self-signed certificate
- `-cert`: Path to certificate file (default: cert/server.crt)
- `-key`: Path to private key file (default: cert/server.key)
- `-lets-encrypt`: Use Let's Encrypt for automatic SSL certificates
- `-domain`: Domain name for Let's Encrypt certificate (required with -lets-encrypt)
- `-cache-dir`: Directory to cache Let's Encrypt certificates (default: cert-cache)
- `-challenge-port`: Port for Let's Encrypt HTTP challenge (default: 80, set to 0 to disable)

## Deployment

### Systemd Service
A systemd service file (`azenv.service`) is provided with multiple configuration options:

```bash
# Copy and edit the service file
sudo cp azenv.service /etc/systemd/system/
sudo nano /etc/systemd/system/azenv.service

# Uncomment and modify the appropriate ExecStart line:
# - Option 1: Self-signed certificates (default)
# - Option 2: Let's Encrypt standalone
# - Option 3: Let's Encrypt with reverse proxy 
# - Option 4: Custom certificates

sudo systemctl daemon-reload
sudo systemctl start azenv
sudo systemctl enable azenv
```

### Reverse Proxy Setup with Nginx

For production deployments behind Nginx, use the provided `nginx-example.conf`:

1. **AZenv configuration** (systemd service):
   ```bash
   ExecStart=/usr/local/bin/azenv -p 9080 -sp 9443 -ssl -lets-encrypt -domain yourdomain.com -challenge-port 0 -cache-dir /var/lib/azenv/certs
   ```

2. **Nginx configuration**:
   ```bash
   sudo cp nginx-example.conf /etc/nginx/sites-available/azenv
   sudo ln -s /etc/nginx/sites-available/azenv /etc/nginx/sites-enabled/
   sudo nginx -t && sudo systemctl reload nginx
   ```

This setup allows:
- Nginx handles port 80/443 and SSL termination
- Let's Encrypt ACME challenges proxied to AZenv
- AZenv runs on custom ports (9080/9443)
- Certificate management handled by AZenv's autocert

### GitHub Actions
The project uses GitHub Actions (`.github/workflows/release.yml`) for:
- Cross-platform builds (Linux, macOS, Windows for AMD64 and ARM64)
- Automated releases when tags are pushed
- Binary artifact generation with compression

## Key Functions

- `generateCert()`: Creates self-signed TLS certificates and keys
- `certFilesExist()`: Checks for existing certificate files
- `azenvHandler()`: Main HTTP handler that formats request/environment data
- `main()`: Entry point handling CLI flags and server setup