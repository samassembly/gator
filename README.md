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

## Database Setup

- Start the postgres service 
```bash 
sudo service postgresql start
```
- Open Postgres shell 
```bash 
sudo -u postgres psql
```
- Create a new database with 
```bash 
CREATE DATABASE gator;
```
- Connect to the new database with:
```bash 
\c gator
```
- Setup postgres password: 
```bash 
ALTER USER postgres PASSWORD '<password>';
```
- Leave Postgres shell: 
```bash 
exit
```

## 🐊 Gator Installation

- Navigate to desired install location
- Clone this repo with command: 
```bash 
git clone https://github.com/samassembly/gator
```
- Navigate to base of repo and enter: 
```bash 
go install .
```

## 🪿 Goose Database Migration

- This project uses Goose to maintain the database versions
- Install Goose with following: 
```bash 
go install github.com/pressly/goose/v3/cmd/goose@latest
```
- Navigate from your gator base repo to 
```bash 
cd ./sql/schema
```
- Run the Up migration: 
```bash 
goose postgres postgres://<username>:<password>@localhost:5432/gator up
```

## ⚙️ Create Configuration File

- Navigate to home directory: 
```bash 
cd ~
```
- Create a config file named 
```bash 
.gatorconfig.json
```
- Edit the file to contain the following: 
```bash 
{"db_url":"postgres://<username>:<password>@localhost:5432/gator?sslmode=disable","current_user_name":""}
```
- The field for current_user_name will be blank for this step, we can set this with a gator command