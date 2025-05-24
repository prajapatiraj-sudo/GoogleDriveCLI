
# 🛠 Build & Deploy Instructions for gdrivecli

This guide walks you through how to build, compile, and distribute the `gdrivecli` tool — a Go-based CLI for uploading files to Google Drive, with real-time progress, resumable uploads, and trash-cleaning.

---

## ✅ Prerequisites

- Go 1.17+ installed on your system
- Git (for source control)
- A Google Cloud project with Drive API enabled
- `credentials.json` (OAuth 2.0 client credentials from Google)

---

## 📁 Project Structure

```
gdrivecli/
├── main.go
├── cmd/
│   ├── root.go
│   ├── upload.go
│   ├── config.go
│   └── clean.go
├── utils/
│   ├── auth.go
│   └── (optional) clean.go
└── credentials.json (used once to authorize)
```

---

## 🔐 Step 1: Configure OAuth Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/).
2. Enable **Google Drive API**.
3. Create OAuth2.0 credentials for a **Desktop app**.
4. Download `credentials.json`.

Place it in:
```
~/.config/gdrivecli/credentials.json
```

The CLI will also store:
```
~/.config/gdrivecli/token.json
```

---

## ⚙️ Step 2: Build for All Linux Users (Static Binary)

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gdrivecli ./gdrivecli
```

✅ This builds a fully static binary:
- No glibc version issues
- Works on any Linux (Ubuntu 18.04, 20.04, 22.04, etc.)

Verify it’s static:

```bash
ldd gdrivecli
# should say: not a dynamic executable
```

---

## 📦 Step 3: Install the CLI System-Wide (Optional)

```bash
sudo mv gdrivecli /usr/local/bin/
chmod +x /usr/local/bin/gdrivecli
```

Now usable as:

```bash
gdrivecli upload myfile.zip
gdrivecli clean-trash
```

---

## 🚀 Commands Available

- `gdrivecli config` → Authorize with Google Drive
- `gdrivecli upload file.zip` → Upload file with progress bar + ETA
- `gdrivecli clean-trash` → Permanently delete all trashed files

---

## 🛡️ Security Tips

- Never upload `credentials.json` or `token.json` to GitHub
- Add to `.gitignore`:
```
gdrivecli
*.json
```

---

## 🧪 Run without building (for dev testing)

```bash
go run ./gdrivecli
```

## 📄 License

MIT License.
