# Installation

This guide explains how to install **DAI** using either:
- **Prebuilt binaries** from the official `dai` repository (fastest way)
- **Build from source** using Go (for those who want to compile or modify DAI)

---

## 📦 Install from prebuilt binary (recommended)

1. **Download the latest release** for your OS from:  
   [https://github.com/gorankrgovic/dai/releases/latest](https://github.com/gorankrgovic/dai/releases/latest)

2. **Unpack the archive** (if downloaded as `.zip` or `.tar.gz`):
```bash
tar -xvzf dai_<version>_linux_amd64.tar.gz
# or
unzip dai_<version>_macos_arm64.zip
```

3. **Make the binary executable**:
```bash
chmod +x dai
```

4. **Move it to your PATH**:
```bash
sudo mv dai /usr/local/bin
```

5. **Verify installation**:
```bash
dai --version
```

---

## 🛠 Build from source

You can build DAI from source if you want the latest development version or if you plan to contribute.

### 1. Install Go (if not already installed)
DAI requires **Go 1.21+**.  
Check your version:
```bash
go version
```
If you don't have Go or your version is older, download it from:  
[https://go.dev/dl/](https://go.dev/dl/)

After installation, make sure `$GOPATH/bin` is in your `PATH`.

---

### 2. Clone the repository
```bash
git clone https://github.com/gorankrgovic/dai.git
cd dai
```

---

### 3. Build the binary
```bash
go build -o dai
```

---

### 4. Move to your PATH
```bash
sudo mv dai /usr/local/bin
```

---

### 5. Verify installation
```bash
dai --version
```

---

## First-time setup

After installing (either method), run:
```bash
dai config
```
This will start the configuration wizard to set your **OpenAI API key** and select the preferred AI model.  
Configuration is stored locally at:
```text
~/.dai/config.yaml
```

---

Next: [Configuration](configuration.md)
