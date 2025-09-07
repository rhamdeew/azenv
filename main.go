package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

// generateCert creates a self-signed certificate and key, saving them to the specified paths
func generateCert(certPath, keyPath string) error {
	// Create directory for cert if it doesn't exist
	certDir := filepath.Dir(certPath)
	if err := os.MkdirAll(certDir, 0755); err != nil {
		return fmt.Errorf("failed to create certificate directory: %w", err)
	}

	keyDir := filepath.Dir(keyPath)
	if err := os.MkdirAll(keyDir, 0755); err != nil {
		return fmt.Errorf("failed to create key directory: %w", err)
	}

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Prepare certificate template
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // Valid for 1 year

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"AZenv Self-Signed Certificate"},
			CommonName:   "localhost",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	// Create the certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Save the certificate
	certOut, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %w", certPath, err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("failed to write certificate to file: %w", err)
	}

	// Save the private key
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %w", keyPath, err)
	}
	defer keyOut.Close()

	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(keyOut, privBlock); err != nil {
		return fmt.Errorf("failed to write private key to file: %w", err)
	}

	fmt.Printf("Generated self-signed certificate at %s and key at %s\n", certPath, keyPath)
	return nil
}

// certFilesExist checks if both certificate and key files exist
func certFilesExist(certPath, keyPath string) bool {
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func azenvHandler(w http.ResponseWriter, r *http.Request) {
	// Only respond to /azenv path
	if r.URL.Path != "/azenv" {
		http.NotFound(w, r)
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "text/html")

	// Start the HTML output
	fmt.Fprintf(w, "<html>\n<head>\n<title>AZenv Go Version</title>\n</head>\n<body>\n<pre>\n")

	// Add REMOTE_ADDR and REMOTE_PORT
	host, port, _ := strings.Cut(r.RemoteAddr, ":")
	fmt.Fprintf(w, "REMOTE_ADDR = %s\n", host)
	fmt.Fprintf(w, "REMOTE_PORT = %s\n", port)

	// Add REQUEST values
	fmt.Fprintf(w, "REQUEST_URI = %s\n", r.URL.RequestURI())
	fmt.Fprintf(w, "REQUEST_METHOD = %s\n", r.Method)

	// Process HTTP headers (format to match PHP's $_SERVER format)
	fmt.Fprintf(w, "HTTP_HOST = %s\n", r.Host)
	for name, values := range r.Header {
		// Skip Host header as we've already displayed it
		if strings.ToLower(name) == "host" {
			continue
		}

		// Format header name to match PHP's $_SERVER convention (HTTP_*)
		headerName := "HTTP_" + strings.ToUpper(strings.Replace(name, "-", "_", -1))
		for _, value := range values {
			fmt.Fprintf(w, "%s = %s\n", headerName, value)
		}
	}

	// Add REQUEST_TIME and REQUEST_TIME_FLOAT
	now := time.Now()
	nowUnix := float64(now.UnixNano()) / 1e9
	fmt.Fprintf(w, "REQUEST_TIME_FLOAT = %.4f\n", nowUnix)
	fmt.Fprintf(w, "REQUEST_TIME = %d\n", int64(nowUnix))

	// Close the HTML output
	fmt.Fprintf(w, "</pre>\n</body>\n</html>")
}

func main() {
	// Define command line flags
	httpPort := flag.Int("p", 8080, "HTTP port to listen on")
	httpsPort := flag.Int("sp", 8443, "HTTPS port to listen on")
	enableHTTPS := flag.Bool("ssl", false, "Enable HTTPS server")
	certPath := flag.String("cert", "cert/server.crt", "Path to SSL certificate file")
	keyPath := flag.String("key", "cert/server.key", "Path to SSL key file")
	genCert := flag.Bool("gen-cert", false, "Generate a self-signed certificate if none exists")
	letsEncrypt := flag.Bool("lets-encrypt", false, "Use Let's Encrypt for automatic SSL certificates")
	domain := flag.String("domain", "", "Domain name for Let's Encrypt certificate (required with -lets-encrypt)")
	cacheDir := flag.String("cache-dir", "cert-cache", "Directory to cache Let's Encrypt certificates")
	challengePort := flag.Int("challenge-port", 80, "Port for Let's Encrypt HTTP challenge (0 to disable built-in challenge server)")
	flag.Parse()

	// Validate Let's Encrypt configuration
	if *letsEncrypt && *enableHTTPS {
		if *domain == "" {
			fmt.Println("Error: -domain is required when using -lets-encrypt")
			os.Exit(1)
		}
		fmt.Printf("Using Let's Encrypt for domain: %s\n", *domain)
	} else if !*letsEncrypt && *enableHTTPS {
		// Generate certificate if requested and files don't exist
		if *genCert || !certFilesExist(*certPath, *keyPath) {
			fmt.Println("Generating self-signed certificate...")
			if err := generateCert(*certPath, *keyPath); err != nil {
				fmt.Printf("Error generating certificate: %s\n", err)
				*enableHTTPS = false // Disable HTTPS if certificate generation fails
			}
		}
	}

	// Register handler
	http.HandleFunc("/", azenvHandler)

	// Start servers
	if *enableHTTPS {
		httpsAddr := fmt.Sprintf(":%d", *httpsPort)
		
		if *letsEncrypt {
			// Let's Encrypt setup
			certManager := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(*domain),
				Cache:      autocert.DirCache(*cacheDir),
			}

			// Create TLS config with autocert
			server := &http.Server{
				Addr:      httpsAddr,
				TLSConfig: &tls.Config{
					GetCertificate: certManager.GetCertificate,
				},
			}

			// Start HTTP server for ACME challenge (if enabled)
			if *challengePort != 0 {
				go func() {
					challengeAddr := fmt.Sprintf(":%d", *challengePort)
					fmt.Printf("Starting HTTP challenge server on %s for Let's Encrypt validation\n", challengeAddr)
					err := http.ListenAndServe(challengeAddr, certManager.HTTPHandler(nil))
					if err != nil {
						fmt.Printf("HTTP challenge server error: %s\n", err)
					}
				}()
			} else {
				fmt.Println("Built-in challenge server disabled. Ensure reverse proxy handles /.well-known/acme-challenge/")
			}

			// Start HTTPS server with Let's Encrypt
			go func() {
				fmt.Printf("HTTPS server starting on %s with Let's Encrypt\n", httpsAddr)
				fmt.Printf("Access environment variables at https://%s%s/azenv\n", *domain, httpsAddr)
				err := server.ListenAndServeTLS("", "")
				if err != nil {
					fmt.Printf("HTTPS server error: %s\n", err)
				}
			}()
		} else {
			// Traditional certificate setup
			go func() {
				fmt.Printf("HTTPS server starting on %s\n", httpsAddr)
				fmt.Printf("Access environment variables at https://localhost%s/azenv\n", httpsAddr)
				err := http.ListenAndServeTLS(httpsAddr, *certPath, *keyPath, nil)
				if err != nil {
					fmt.Printf("HTTPS server error: %s\n", err)
				}
			}()
		}
	}

	// Always start HTTP server (main thread)
	httpAddr := fmt.Sprintf(":%d", *httpPort)
	fmt.Printf("HTTP server starting on %s\n", httpAddr)
	fmt.Printf("Access environment variables at http://localhost%s/azenv\n", httpAddr)
	http.ListenAndServe(httpAddr, nil)
}
