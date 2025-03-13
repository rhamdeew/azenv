package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
)

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
	// Define command line flag for port
	port := flag.Int("p", 8080, "port to listen on")
	flag.Parse()

	// Register handler only for /azenv path
	http.HandleFunc("/", azenvHandler)
	
	// Start the server
	serverAddr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Server starting on %s\n", serverAddr)
	fmt.Printf("Access environment variables at http://localhost%s/azenv\n", serverAddr)
	http.ListenAndServe(serverAddr, nil)
}
