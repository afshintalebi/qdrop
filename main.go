package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mdp/qrterminal/v3"
)

const uploadHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>QDrop Receive</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; background: #f0f2f5; margin: 0; }
        .card { background: white; padding: 30px; border-radius: 12px; box-shadow: 0 4px 12px rgba(0,0,0,0.1); text-align: center; width: 80%; max-width: 400px; }
        h2 { color: #333; margin-top: 0; }
        input[type="file"] { margin: 20px 0; width: 100%; }
        button { background: #007bff; color: white; border: none; padding: 12px 24px; border-radius: 6px; font-size: 16px; cursor: pointer; width: 100%; transition: background 0.3s; }
        button:hover { background: #0056b3; }
        button:disabled { background: #cccccc; cursor: not-allowed; }
        
        /* Progress Bar Styles */
        .progress-container { display: none; margin-top: 20px; text-align: left; }
        .progress-bar-bg { width: 100%; background-color: #e9ecef; border-radius: 8px; overflow: hidden; margin-top: 8px; height: 20px; }
        .progress-bar-fill { height: 100%; background-color: #28a745; width: 0%; transition: width 0.2s; }
        .status-text { font-size: 14px; color: #666; display: flex; justify-content: space-between; }
    </style>
</head>
<body>
    <div class="card">
        <h2>Drop File Here 📲</h2>
        <form id="uploadForm">
            <input type="file" id="fileInput" name="file" required>
            <button type="submit" id="submitBtn">Upload to Computer</button>
        </form>

        <div class="progress-container" id="progressContainer">
            <div class="status-text">
                <span id="statusLabel">Uploading...</span>
                <span id="percentLabel">0%</span>
            </div>
            <div class="progress-bar-bg">
                <div class="progress-bar-fill" id="progressBar"></div>
            </div>
        </div>
    </div>

    <script>
        document.getElementById('uploadForm').addEventListener('submit', function(e) {
            e.preventDefault(); // Prevent default form submission

            var fileInput = document.getElementById('fileInput');
            if (fileInput.files.length === 0) return;

            var file = fileInput.files[0];
            var formData = new FormData();
            formData.append('file', file);

            // UI Elements
            var submitBtn = document.getElementById('submitBtn');
            var progressContainer = document.getElementById('progressContainer');
            var progressBar = document.getElementById('progressBar');
            var percentLabel = document.getElementById('percentLabel');
            var statusLabel = document.getElementById('statusLabel');

            // Setup UI for uploading
            submitBtn.disabled = true;
            submitBtn.innerText = "Uploading...";
            progressContainer.style.display = "block";

            // Format file size
            var fileSizeMB = (file.size / (1024 * 1024)).toFixed(2);
            statusLabel.innerText = "Sending " + fileSizeMB + " MB...";

            // Create AJAX request
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/receive', true);

            // Listen to progress event
            xhr.upload.onprogress = function(event) {
                if (event.lengthComputable) {
                    var percentComplete = Math.round((event.loaded / event.total) * 100);
                    progressBar.style.width = percentComplete + '%';
                    percentLabel.innerText = percentComplete + '%';
                }
            };

            // Listen to success/failure
            xhr.onload = function() {
                if (xhr.status === 200) {
                    // Replace the whole page with the success HTML from Go server
                    document.body.innerHTML = xhr.responseText;
                } else {
                    statusLabel.innerText = "Error uploading file!";
                    statusLabel.style.color = "red";
                    submitBtn.disabled = false;
                    submitBtn.innerText = "Try Again";
                }
            };

            xhr.onerror = function() {
                statusLabel.innerText = "Network Error!";
                statusLabel.style.color = "red";
                submitBtn.disabled = false;
                submitBtn.innerText = "Try Again";
            };

            // Send the file
            xhr.send(formData);
        });
    </script>
</body>
</html>
`

const successHTML = `
<!DOCTYPE html>
<html>
<body style="font-family: sans-serif; text-align: center; padding-top: 50px; background: #e6ffed;">
    <h1 style="color: #28a745;">✅ Success!</h1>
    <p>File has been sent to the computer.</p>
    <p>You can close this page.</p>
</body>
</html>
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  Send a file:    qdrop <file_or_folder_path>")
		fmt.Println("  Receive a file: qdrop receive")
		os.Exit(1)
	}
	targetPath := os.Args[1]

	// Check if user wants to receive a file
	if targetPath == "receive" {
		if err := runReceive(); err != nil {
			log.Fatalf("Fatal error: %v\n", err)
		}
		return
	}

	// Delegate main logic to 'run' function to handle errors cleanly
	if err := run(targetPath); err != nil {
		log.Fatalf("Fatal error: %v\n", err)
	}
}

// run orchestrates the setup, starting the server, and graceful shutdown
func run(targetPath string) error {
	// 1. Validate the path
	fileInfo, err := os.Stat(targetPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("path '%s' does not exist", targetPath)
	}

	// 2. Network setup
	ip, err := getLocalIP()
	if err != nil {
		return fmt.Errorf("finding local IP: %v", err)
	}

	listener, err := net.Listen("tcp", ":0") // Dynamic port
	if err != nil {
		return fmt.Errorf("finding a free port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	// 3. Prepare URL and Server
	baseName := filepath.Base(targetPath)
	isDir := fileInfo.IsDir()
	
	// If it's a directory, append .zip to the URL endpoint for clarity
	endpointName := baseName
	if isDir {
		endpointName += ".zip"
	}
	downloadURL := fmt.Sprintf("http://%s:%d/%s", ip, port, endpointName)

	done := make(chan struct{})
	srv := &http.Server{}

	// Register the handler
	http.HandleFunc(fmt.Sprintf("/%s", endpointName), createDownloadHandler(targetPath, isDir, done))

	// 4. UI output
	printUI(targetPath, downloadURL, isDir)

	// 5. Start Server
	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	// 6. Wait for download to finish
	<-done
	fmt.Println("\n✅ Transfer complete! Shutting down gracefully...")

	// 7. Graceful Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

// createDownloadHandler returns an http.HandlerFunc that handles either a single file or a directory zip
func createDownloadHandler(targetPath string, isDir bool, done chan<- struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("📥 Transfer started to %s\n", r.RemoteAddr)

		if isDir {
			// Serve directory as an on-the-fly zip stream
			err := streamZipFolder(w, targetPath)
			if err != nil {
				log.Printf("Error creating zip: %v\n", err)
			}
		} else {
			// Serve a single file
			http.ServeFile(w, r, targetPath)
		}

		// Signal the main goroutine to shut down the server
		go func() {
			done <- struct{}{}
		}()
	}
}

// streamZipFolder zips a directory directly to the http.ResponseWriter
// This is memory efficient: it streams data without creating a physical .zip file on disk
func streamZipFolder(w http.ResponseWriter, sourceDir string) error {
	zipFilename := filepath.Base(sourceDir) + ".zip"

	// Set headers for file download
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, zipFilename))

	// Create a zip writer that writes directly to the HTTP response stream
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close() // Ensure the zip footer is written

	// Walk through the directory and add files to the zip
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a proper internal zip path (relative to the source folder)
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Preserve directory structure inside the zip
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		
		// Ensure correct slash direction for zip formats
		header.Name = filepath.ToSlash(relPath)
		if info.IsDir() {
			header.Name += "/"
			return nil // Don't try to write content for directories
		}

		header.Method = zip.Deflate // Use compression

		// Create a writer for this specific file inside the zip
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Open the actual file on disk
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Stream the file content into the zip writer
		_, err = io.Copy(writer, file)
		return err
	})
}

// getLocalIP finds the local IP address of the machine on the Wi-Fi/LAN network
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && isPrivateIP(ipnet.IP) {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("could not find local IP address")
}

// isPrivateIP checks if the IP is a local network IP (e.g., 192.168.x.x, 10.x.x.x)
func isPrivateIP(ip net.IP) bool {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// printUI displays the terminal interface with the QR code
func printUI(targetPath, downloadURL string, isDir bool) {
	fmt.Println(strings.Repeat("=", 45))
	if isDir {
		fmt.Printf("📁 Directory: %s (Will be zipped on-the-fly)\n", filepath.Base(targetPath))
	} else {
		fmt.Printf("📄 File: %s\n", filepath.Base(targetPath))
	}
	fmt.Println(strings.Repeat("=", 45))

	// Print compact QR code
	qrterminal.GenerateHalfBlock(downloadURL, qrterminal.L, os.Stdout)

	fmt.Printf("\n🔗 Or open manually: %s\n", downloadURL)
	fmt.Println("⏳ Waiting for download... (Auto-closes when finished)")
}

// runReceive orchestrates the setup for receiving files from phone to computer
func runReceive() error {
	// 1. Network setup (reused logic to avoid modifying original run function)
	ip, err := getLocalIP()
	if err != nil {
		return fmt.Errorf("finding local IP: %v", err)
	}

	listener, err := net.Listen("tcp", ":0") // Dynamic port
	if err != nil {
		return fmt.Errorf("finding a free port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	// 2. Prepare URL and Server
	downloadURL := fmt.Sprintf("http://%s:%d/receive", ip, port)
	done := make(chan struct{})
	srv := &http.Server{}

	// Register the receive handler
	http.HandleFunc("/receive", createUploadHandler(done))

	// 3. UI output
	printReceiveUI(downloadURL)

	// 4. Start Server
	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	// 5. Wait for upload to finish
	<-done
	fmt.Println("\n✅ File received successfully! Shutting down gracefully...")

	// 6. Graceful Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

// createUploadHandler handles both GET (showing form) and POST (receiving file) requests
func createUploadHandler(done chan<- struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Serve the HTML upload form to the phone
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(uploadHTML))
			return
		}

		if r.Method == http.MethodPost {
			// Parse the multipart form (Max 10 MB in RAM, rest streamed to disk)
			err := r.ParseMultipartForm(10 << 20)
			if err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			// Get the file from the request
			file, header, err := r.FormFile("file")
			if err != nil {
				http.Error(w, "Failed to get file", http.StatusBadRequest)
				return
			}
			defer file.Close()

			// Create a new file in the current directory on the computer
			dst, err := os.Create(header.Filename)
			if err != nil {
				http.Error(w, "Failed to save file on server", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			// Stream the uploaded file to disk
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, "Failed to save file", http.StatusInternalServerError)
				return
			}

			log.Printf("📥 Saved file: %s (%d bytes)\n", header.Filename, header.Size)

			// Send success HTML to phone
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(successHTML))

			// Signal the main goroutine to shut down the server
			go func() {
				done <- struct{}{}
			}()
		}
	}
}

// printReceiveUI displays the terminal interface for the receive mode
func printReceiveUI(downloadURL string) {
	fmt.Println(strings.Repeat("=", 45))
	fmt.Println("📥 RECEIVE MODE")
	fmt.Println("Files will be saved in the current directory.")
	fmt.Println(strings.Repeat("=", 45))
	
	// Print compact QR code
	qrterminal.GenerateHalfBlock(downloadURL, qrterminal.L, os.Stdout)
	
	fmt.Printf("\n🔗 Or open manually on your phone: %s\n", downloadURL)
	fmt.Println("⏳ Waiting for you to upload a file... (Auto-closes when finished)")
}