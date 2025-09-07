# AZenv Go Version

AZenv is a lightweight Go web server that displays environment variables and HTTP request information.

## Features

- Displays client information (IP address, port)
- Shows HTTP request details (URI, method, headers)
- Provides request timing information
- Runs as a standalone web server
- Supports HTTPS with self-signed certificates
- Let's Encrypt automatic certificate management
- Cross-platform support (Linux, macOS, Windows)
- Docker support with lightweight Alpine-based images

## Installation

### Using Prebuilt Binaries

1. Download the appropriate binary for your platform from the [Releases](https://github.com/rhamdeew/azenv/releases) page.
2. Extract the archive:
   ```bash
   # For Linux/macOS
   tar -xzf azenv-{os}-{arch}.tar.gz

   # For Windows
   unzip azenv-windows-amd64.exe.zip
   ```
3. Move the binary to a location in your PATH:
   ```bash
   # Linux/macOS
   sudo mv azenv-{os}-{arch} /usr/local/bin/azenv
   chmod +x /usr/local/bin/azenv

   # Windows
   # Move to a directory in your PATH
   ```

### Building from Source

1. Ensure you have Go 1.21 or later installed.
2. Clone the repository:
   ```bash
   git clone https://github.com/rhamdeew/azenv.git
   cd azenv
   ```
3. Build the binary:
   ```bash
   go build -o azenv
   ```
4. Optionally, move to a location in your PATH:
   ```bash
   sudo mv azenv /usr/local/bin/
   ```

### Using Docker

1. Using Docker Compose (recommended):
   ```bash
   # Clone the repository
   git clone https://github.com/rhamdeew/azenv.git
   cd azenv
   
   # Start the service
   docker-compose up -d
   
   # View logs
   docker-compose logs -f azenv
   ```

2. Using Docker directly:
   ```bash
   # Build the image
   docker build -t azenv .
   
   # Run HTTP only
   docker run -d -p 8080:8080 --name azenv azenv
   
   # Run with HTTPS and self-signed certificates
   docker run -d -p 8080:8080 -p 8443:8443 -v ./cert:/app/cert --name azenv azenv ./azenv -ssl -gen-cert
   ```

3. Docker configuration options:
   ```bash
   # Custom ports
   docker-compose run azenv ./azenv -p 8080 -sp 8443
   
   # With Let's Encrypt (requires domain)
   docker-compose run azenv ./azenv -ssl -lets-encrypt -domain example.com
   ```

### Installing as a Systemd Service (Linux)

1. Copy the provided service file to the systemd directory:
   ```bash
   sudo cp azenv.service /etc/systemd/system/
   ```
2. Ensure the binary is in the correct location:
   ```bash
   sudo cp azenv /usr/local/bin/
   sudo chmod +x /usr/local/bin/azenv
   ```
3. Start and enable the service:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl start azenv
   sudo systemctl enable azenv
   ```

## Usage

### Running the Server

```bash
# Run HTTP on default port (8080)
azenv

# Run HTTP on a specific port
azenv -p 3000

# Run both HTTP (8080) and HTTPS (8443) with auto-generated self-signed certificate
azenv -ssl -gen-cert

# Run HTTP and HTTPS on custom ports
azenv -p 3000 -sp 4430 -ssl -gen-cert

# Use existing certificate and key files
azenv -ssl -cert /path/to/certificate.crt -key /path/to/private.key
```

### Command Line Options

- `-p`: HTTP port (default: 8080)
- `-sp`: HTTPS port (default: 8443)
- `-ssl`: Enable HTTPS server
- `-gen-cert`: Auto-generate a self-signed certificate if it doesn't exist
- `-cert`: Path to certificate file (default: cert/server.crt)
- `-key`: Path to private key file (default: cert/server.key)

### Accessing the Environment Information

Once the server is running, open your web browser and navigate to:

```
# For HTTP
http://localhost:8080/azenv

# For HTTPS (when enabled)
https://localhost:8443/azenv
```

When accessing via HTTPS with a self-signed certificate, you'll need to accept the browser security warning.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
