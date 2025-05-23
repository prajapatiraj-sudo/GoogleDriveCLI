# 🚀 gdrivecli

A blazing-fast, terminal-based CLI tool to upload files to Google Drive with real-time progress bars, resumable uploads, and support for large files.

Built in Go using the Google Drive API and Cobra CLI framework.

---

## ⚙️ Features

- ✅ Upload any file to your Google Drive
- 🔁 Resumable uploads (no loss if interrupted)
- 📊 Real-time progress bar with percentage & ETA
- 🧠 Smart token management (OAuth2 + caching)
- 💬 Clean CLI interface using [Cobra](https://github.com/spf13/cobra)
- 🪄 Cross-platform (Linux, macOS, Windows)

---

## 📦 Installation

```bash
git clone https://github.com/prajapatiraj-sudo/gdrivecli.git
cd gdrivecli
go build -o gdrivecli ./gdrivecli
