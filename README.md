# ⚡ QDrop (Quick Drop)

[![Go Version](https://img.shields.io/github/go-mod/go-version/afshintalebi/qdrop)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

**QDrop** is a blazing-fast, zero-dependency CLI tool built with Go that allows you to instantly share files and folders from your computer to your phone (or any other device) over your local Wi-Fi network using a QR Code.

Forget about emailing files to yourself, using messaging apps, or connecting USB cables.

## ✨ Features

- **Two-Way Transfer:** Send files to your phone, or receive files from your phone directly to your computer's current directory!
- **Blazing Fast:** Transfers happen over your local network. No internet bandwidth is consumed.
- **Zero Configuration:** Automatically finds your local IP and selects a dynamic free port to avoid conflicts.
- **On-the-fly Zipping:** Need to send a whole folder? Just pass the folder path! QDrop will securely compress it into a `.zip` stream directly to the downloader *without* creating temporary files on your disk.
- **Auto-Shutdown:** The server gracefully shuts down the moment your file finishes downloading or uploading. You don't even need to press `Ctrl+C`.
- **Cross-Platform:** Works seamlessly on Windows, macOS, and Linux.

## 🚀 Installation

Since QDrop is built with Go, you can install it globally on your system with a single command. 

*Make sure you have [Go installed](https://go.dev/doc/install) on your system.*

```bash
go install github.com/afshintalebi/qdrop@latest
```

## 📖 Usage
Using QDrop is incredibly simple. Just open your terminal anywhere and type qdrop followed by the file, folder, or command you want to execute.

###  1. Share a single file (Computer ➔ Phone):
```bash
qdrop /path/to/your/image.png
```

### 2. Share an entire folder (Computer ➔ Phone):
```bash
qdrop /path/to/your/folder/
```

What happens next?
1. QDrop will instantly generate a QR Code right in your terminal.
2. Scan the QR code with your phone's camera (ensure both devices are on the same Wi-Fi).
3. The file/folder downloads immediately to your phone.
4. QDrop detects the successful transfer and automatically safely shuts down. Done! 🎉

### 3. Receive a file (Phone ➔ Computer):
Want to send a photo or video from your phone to your computer? Navigate to your desired folder in the terminal and type:
```bash
qdrop receive
```

What happens next?
1. Scan the generated QR Code with your phone.
2. A clean, minimal web page opens on your phone's browser.
3. Select the file you want to upload and tap "Upload".
4. The file is instantly saved to your computer's current terminal directory, and QDrop shuts down automatically.

## 🧠 Under the Hood (For Developers)
- Clean Architecture & Concurrency: Uses Go's powerful goroutines and channels to handle background server execution and graceful shutdowns without blocking the main thread.
- Memory Efficiency (Download): When transferring large directories, QDrop utilizes archive/zip and io.Copy to stream data directly into the HTTP ResponseWriter. This means you can transfer a 10GB folder, and your RAM usage will stay near zero.
- Memory Efficiency (Upload): During 'receive' mode, QDrop safely handles massive file uploads by streaming multipart form data directly to the disk after a small RAM buffer, ensuring your system doesn't crash from out-of-memory errors.
- Dynamic Port Allocation: By explicitly asking the OS to listen on port :0, QDrop avoids "address already in use" errors, allowing multiple instances to run concurrently without clashing.

## 🤝 Contributing
Pull requests are welcome! Feel free to open issues if you want to suggest new features or report bugs.

## 📄 License
This project is licensed under the MIT License.
