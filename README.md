# Prerequisites

Before running this application, ensure the following dependencies are installed on your system:

## ✅ Requirements

- Go (version 1.20 or higher recommended)
- PostgreSQL (version 13 or higher recommended)

---

## 🛠️ Installation Steps

### 1. Install Go

#### macOS (using Homebrew)

```bash
brew install go
```

#### Ubuntu / Debian

```bash
sudo apt update
sudo apt install golang-go
```

#### Windows

- Download the installer from the official Go website
- Run the installer and follow the setup instructions

> After installation, verify with:
```bash
go version
```

---

### 2. Install PostgreSQL

#### macOS (using Homebrew)

```bash
brew install postgresql
brew services start postgresql
```

#### Ubuntu / Debian

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

#### Windows

- Download the installer from the official PostgreSQL website
- Run the installer and follow the setup instructions

> After installation, verify with:
```bash
psql --version
```

---

## 🧪 Verify Setup

To confirm everything is working:

```bash
go version       # Should print your Go version
psql --version   # Should print your PostgreSQL version
```
