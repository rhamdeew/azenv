# AZenv Go Version

AZenv is a lightweight Go web server that displays environment variables and HTTP request information.

## Features

- Displays client information (IP address, port)
- Shows HTTP request details (URI, method, headers)
- Provides request timing information
- Runs as a standalone web server
- Cross-platform support (Linux, macOS, Windows)

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
# Run on default port (8080)
azenv

# Run on a specific port
azenv -p 3000
```

### Accessing the Environment Information

Once the server is running, open your web browser and navigate to:

```
http://localhost:8080/azenv
```

This will display all the environment variables and request information in your browser.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
