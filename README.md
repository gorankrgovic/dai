# DAI – Debug & Develop AI CLI

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gorankrgovic/dai)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/gorankrgovic/dai)](https://github.com/gorankrgovic/dai/releases/latest)
[![Docs](https://img.shields.io/badge/docs-online-blue)](https://gorankrgovic.github.io/dai/)

**DAI** (Debug & Develop AI) is a cross-platform CLI tool for AI-assisted debugging and development — right from your terminal.  
It works on Linux, macOS, and Windows, with zero complex setup, and integrates with GitHub for advanced code triage.

---

## ✨ Features
- **AI-powered code analysis** – analyze commits or local files with context.
- **GitHub integration** – open issues with AI-generated findings directly from CLI.
- **Local configuration** – no external servers, your data stays with you.
- **Multi-language support** – works with various tech stacks.
- **Fun parrot mode** – `--parrot party|insult|wise` for some extra flavor.

---

## 🚀 Installation

### 📦 Prebuilt binary (recommended)
Download the latest binary for your OS from the [Latest Release](https://github.com/gorankrgovic/dai/releases/latest).

**Linux / macOS:**
```bash
curl -L https://github.com/gorankrgovic/dai/releases/latest/download/dai_linux-amd64.tar.gz -o dai.tar.gz
tar -xzf dai.tar.gz
chmod +x dai
sudo mv dai /usr/local/bin/
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -Uri "https://github.com/gorankrgovic/dai/releases/latest/download/dai_windows-amd64.zip" -OutFile "dai.zip"
Expand-Archive dai.zip -DestinationPath .
```

---

### 🛠 Build from source
Requires **Go 1.21+**.

```bash
git clone https://github.com/gorankrgovic/dai.git
cd dai
go build -o dai
sudo mv dai /usr/local/bin
```

---

## 📚 Documentation
Full documentation is available here:  
[📄 DAI Documentation](https://gorankrgovic.github.io/dai/)

---

## 🤝 Contributing
We welcome contributions!  
See our [Contributing Guidelines](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) for details.

- Report bugs → [Bug report template](.github/ISSUE_TEMPLATE/bug_report.md)
- Request features → [Feature request template](.github/ISSUE_TEMPLATE/feature_request.md)

---

## 📦 Downloads
- [Latest Release](https://github.com/gorankrgovic/dai/releases/latest)
- [All releases](https://github.com/gorankrgovic/dai/releases)

---

## 📜 License
This project is licensed under the terms of the [GNU GPL v3](LICENSE).
