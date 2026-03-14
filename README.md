# ⚡ QDrop (Quick Drop)

[![Go Version](https://img.shields.io/github/go-mod/go-version/afshintalebi/qdrop)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

**QDrop** is a blazing-fast, zero-dependency CLI tool built with Go that allows you to instantly share files and folders from your computer to your phone (or any other device) over your local Wi-Fi network using a QR Code.

Forget about emailing files to yourself, using messaging apps, or connecting USB cables.

## ✨ Features

- **Blazing Fast:** Transfers happen over your local network. No internet bandwidth is consumed.
- **Zero Configuration:** Automatically finds your local IP and selects a dynamic free port to avoid conflicts.
- **On-the-fly Zipping:** Need to send a whole folder? Just pass the folder path! QDrop will securely compress it into a `.zip` stream directly to the downloader *without* creating temporary files on your disk.
- **Auto-Shutdown:** The server gracefully shuts down the moment your file finishes downloading. You don't even need to press `Ctrl+C`.
- **Cross-Platform:** Works seamlessly on Windows, macOS, and Linux.

## 🚀 Installation

Since QDrop is built with Go, you can install it globally on your system with a single command. 

*Make sure you have [Go installed](https://go.dev/doc/install) on your system.*

```bash
go install github.com/afshintalebi/qdrop@latest
```

## 📖 Usage
Using QDrop is incredibly simple. Just open your terminal anywhere and type `qdrop` followed by the file or folder you want to share.

### 1. Share a single file:
```bash
qdrop /path/to/your/image.png
```

### 2. Share an entire folder:
```bash
qdrop /path/to/your/documents/
```

**What happens next?**
1. QDrop will instantly generate a QR Code right in your terminal.
2. Scan the QR code with your phone's camera (ensure both devices are on the same Wi-Fi).
3. The file/folder downloads immediately to your phone.
4. QDrop detects the successful transfer and automatically safely shuts down. Done! 🎉

## 🧠 Under the Hood (For Developers)
- Clean Architecture & Concurrency: Uses Go's powerful goroutines and channels to handle background server execution and graceful shutdowns without blocking the main thread.
- Memory Efficiency: When transferring large directories, QDrop utilizes archive/zip and io.Copy to stream data directly into the HTTP ResponseWriter. This means you can transfer a 10GB folder, and your RAM usage will stay near zero.
- Dynamic Port Allocation: By explicitly asking the OS to listen on port :0, QDrop avoids "address already in use" errors, allowing multiple instances to run concurrently without clashing.

## 🤝 Contributing
Pull requests are welcome! Feel free to open issues if you want to suggest new features or report bugs.

## 📄 License
This project is licensed under the MIT License.


## ⚠️ Disclaimer and Acceptable Use Policy
By downloading, installing, or using this software ("QDrop"), you acknowledge and agree to the following terms:

### 1. No Liability
This software is provided "as is" and "as available", without warranty of any kind, express or implied. In no event shall the author(s), contributors, or copyright holders be liable for any direct, indirect, incidental, special, exemplary, or consequential damages (including, but not limited to, loss of data, network breaches, or business interruption) arising in any way out of the use of this software.

### 2. Responsibility for Content
QDrop is strictly a local peer-to-peer data transfer tool. It does not store, monitor, or host any data on external servers. **You are solely responsible for the files and content you choose to transfer using this tool.** The author(s) do not condone, endorse, or take responsibility for the transfer of:
- Copyrighted material without permission.
- Malware, viruses, or malicious code.
- Any illegal or prohibited content under your local jurisdiction.

### 3. Network Security
QDrop operates by opening a temporary HTTP server on your local area network (LAN/WLAN). While the server shuts down automatically after a transfer, it is your responsibility to ensure you are connected to a **trusted and secure network** before using this tool. The author(s) are not responsible for any unauthorized access to the files you expose on your local network.

### 4. Personal Use Only
This tool is intended for personal, educational, and authorized use only. If you use this software in a corporate or enterprise environment, it is your responsibility to ensure it complies with your organization's IT and security policies.

**By using QDrop, you assume all risks associated with its operation.**

