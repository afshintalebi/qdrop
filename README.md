# ⚡ QDrop (Quick Drop)

[![Go Version](https://img.shields.io/github/go-mod/go-version/afshintalebi/qdrop)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

**QDrop** is a blazing-fast, zero-dependency CLI tool built with Go that allows you to instantly share files and folders from your computer to your phone (or any other device) over your local Wi-Fi network using a QR Code.

Forget about emailing files to yourself, using messaging apps, or connecting USB cables.

## ✨ Features

- **Two-Way Transfer:** Send data to your phone, or receive files and folders from your phone directly to your computer!
- **Intelligent Folder Support:** 
  - **Sending:** Securely compresses massive folders into a `.zip` stream on-the-fly *without* creating temporary files on your disk.
  - **Receiving:** Upload entire folders from your phone. QDrop recreates your exact directory tree and subfolders perfectly on your computer without needing client-side zipping!
- **Interactive UI with Progress & Loaders:** When receiving files, you get a clean web UI where you can select items. It features a **Processing Loader** to prevent browser freezes when selecting heavy folders, a summary of your selection, and a manual upload button with a live progress bar.
- **Blazing Fast:** Transfers happen over your local network. No internet bandwidth is consumed.
- **Zero Configuration:** Automatically finds your local IP and selects a dynamic free port to avoid conflicts.
- **Auto-Shutdown:** The server gracefully shuts down the moment your transfer finishes.
- **Cross-Platform:** Works seamlessly on Windows, macOS, and Linux.

## 🚀 Installation

Since QDrop is built with Go, you can install it globally on your system with a single command. 

*Make sure you have [Go installed](https://go.dev/doc/install) on your system.*

```bash
go install github.com/afshintalebi/qdrop@latest
```

## 📖 Usage
Using QDrop is incredibly simple. Just open your terminal anywhere and type qdrop followed by the file, folder, or command you want to execute.

**1. Share a single file (Computer ➔ Phone):**
```bash
qdrop /path/to/your/image.png
```

**2. Share an entire folder (Computer ➔ Phone):**
```bash
qdrop /path/to/your/folder/
```

What happens next?
1. QDrop will instantly generate a QR Code right in your terminal.
2. Scan the QR code with your phone's camera (ensure both devices are on the same Wi-Fi).
3. The file/folder downloads immediately to your phone.
4. QDrop detects the successful transfer and safely shuts down. Done! 🎉

**3. Receive files/folders (Phone ➔ Computer):**
Want to send photos, videos, or entire folders from your phone to your computer? Navigate to your desired folder in the terminal and type:

```bash
qdrop receive
```
What happens next?
1. Scan the generated QR Code with your phone.
2. A clean, minimal web interface opens on your phone's browser.
3. Select individual files or an entire folder.
4. Review your selection (e.g., "15 items selected (45 MB)").
5. Tap Start Upload and watch the progress bar.
6. The files and their directory structure are instantly saved to your computer, and QDrop shuts down automatically!

## 🧠 Under the Hood (For Developers)
- **Clean Architecture & Concurrency**: Uses Go's powerful goroutines and channels to handle background server execution and graceful shutdowns without blocking the main thread.
- **Memory Efficiency (Download)**: When transferring large directories, QDrop utilizes archive/zip and io.Copy to stream data directly into the HTTP ResponseWriter. This means you can transfer a 10GB folder, and your RAM usage will stay near zero.
- **Client-Optimized Uploads**: Instead of forcing mobile browsers to zip heavy folders (which causes crashes), QDrop uploads raw files using the webkitdirectory path. The Go backend handles the heavy lifting by perfectly reconstructing the nested directory structure on your local disk.
- **Non-Blocking JS UI**: To prevent the mobile browser from freezing when processing folders with thousands of files, QDrop utilizes setTimeout to yield the main thread. This allows a smooth CSS spinner to render while JavaScript handles the heavy counting and FormData generation in the background.
- **Dynamic Port Allocation**: By explicitly asking the OS to listen on port :0, QDrop avoids "address already in use" errors, allowing multiple instances to run concurrently without clashing.
- **Security**: Incoming file paths are strictly sanitized using filepath.Clean to prevent Directory Traversal (../) attacks.

## 🤝 Contributing
Pull requests are welcome! Feel free to open issues if you want to suggest new features or report bugs.

## 📄 License
This project is licensed under the MIT License.
