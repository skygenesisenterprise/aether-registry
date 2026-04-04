<div align="center">

# 💰 Aether Bank CLI

[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](https://github.com/skygenesisenterprise/aether-bank/blob/main/LICENSE) [![Go](https://img.shields.io/badge/Go-1.25+-blue?style=for-the-badge&logo=go)](https://golang.org/) [![Cobra](https://img.shields.io/badge/Cobra-1.9+-lightgrey?style=for-the-badge&logo=go)](https://github.com/spf13/cobra) [![Viper](https://img.shields.io/badge/Viper-1.19+-blue?style=for-the-badge&logo=go)](https://github.com/spf13/viper)

**🔥 Internal Banking CLI - Secure and Scriptable Command-Line Interface for Aether Bank**

A comprehensive internal CLI tool designed for **developers and employees** to interact with the Aether Bank backend in a **stable**, **secure**, and **scriptable** way. Built with Go + Cobra + Viper for maximum reliability and extensibility.

[🚀 Quick Start](#-quick-start) • [📋 Commands](#-commands) • [🛠️ Tech Stack](#️-tech-stack) • [📁 Architecture](#-architecture) • [🔧 Development](#-development) • [🤝 Contributing](#-contributing)

[![GitHub stars](https://img.shields.io/github/stars/skygenesisenterprise/aether-bank?style=social)](https://github.com/skygenesisenterprise/aether-bank/stargazers) [![GitHub forks](https://img.shields.io/github/forks/skygenesisenterprise/aether-bank?style=social)](https://github.com/skygenesisenterprise/aether-bank/network) [![GitHub issues](https://img.shields.io/github/issues/skygenesisenterprise/aether-bank)](https://github.com/skygenesisenterprise/aether-bank/issues)

</div>

---

## 🌟 What is Aether Bank CLI?

**Aether Bank CLI** is an enterprise-grade command-line interface for internal banking operations. It provides **secure**, **scriptable**, and **auditable** access to banking backend services including user management, accounts, transactions, transfers, and ledger operations.

### 🎯 Our Vision

- **🔐 Enterprise-Grade Security** - Secure token management, audit logging, and role-based access
- **📦 Modular Architecture** - Separation between CLI layer and API client for testability
- **⚡ High Performance** - Built with Go for fast execution and minimal resource usage
- **🎯 Developer-Focused** - Scriptable commands with JSON output for automation
- **🏗️ Extensible Design** - Well-structured packages (cmd + internal) for easy extension
- **🧪 Testable** - Mock API client for comprehensive unit testing
- **🌐 Multi-Environment** - Support for staging and production environments

---

## 📊 Current Status

> **✅ V1 Available**: Fully functional CLI with authentication, user management, accounts, transactions, transfers, ledger, and debug commands.

### ✅ **Currently Implemented**

#### 🏗️ **Core Foundation**

- ✅ **Authentication System** - Login/logout with secure token storage
- ✅ **User Management** - List and get user information
- ✅ **Account Operations** - List and get account details
- ✅ **Transaction Management** - List, get, and simulate transactions
- ✅ **Transfer Operations** - Create transfers with dry-run support
- ✅ **Ledger Commands** - Audit and ledger entries
- ✅ **Debug Tools** - Log tailing and transaction debugging
- ✅ **Multi-Environment** - Staging and production support
- ✅ **JSON & Table Output** - Human-readable and JSON output modes
- ✅ **Mock API Client** - Complete mock implementation for testing

#### 🛠️ **Development Infrastructure**

- ✅ **Go CLI Framework** - Cobra + Viper implementation
- ✅ **Unit Tests** - Test coverage for key commands
- ✅ **Multi-Platform Build** - Linux and macOS (x86_64 and ARM)
- ✅ **Installation Script** - Simple local install script
- ✅ **Configuration Management** - YAML-based configuration with Viper

### 📋 **Planned Features**

- **Real API Client Integration** - Connect to actual banking backend
- **Advanced Filtering** - Complex query parameters for lists
- **Batch Operations** - Bulk transaction processing
- **Audit Logging** - Comprehensive operation logging
- **Shell Completion** - Bash/Zsh/Fish autocompletion
- **Interactive Mode** - REPL for exploratory operations

---

## 🚀 Quick Start

### 📋 Prerequisites

- **Go** 1.25.0 or higher
- **Git** for cloning
- **Make** (optional, for build shortcuts)

### 🔧 Installation & Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/skygenesisenterprise/aether-bank.git
   cd aether-bank
   ```

2. **Build the CLI**

   ```bash
   cd package/cli
   go build -o bank .
   ```

3. **Install locally**

   ```bash
   chmod +x install.sh
   ./install.sh
   ```

4. **Configure environment**

   ```bash
   # Copy sample config
   cp config.sample.yaml ~/.bank/config.yaml

   # Or use defaults (staging environment)
   bank --help
   ```

### 🎯 **Quick Commands**

```bash
# Authentication
bank auth login --email admin@aetherbank.com --password secret
bank auth whoami
bank auth logout

# Users
bank user list
bank user get user_001

# Accounts
bank account list --user user_001
bank account get acc_001

# Transactions
bank tx list --user user_001
bank tx get tx_001
bank tx simulate

# Transfers
bank transfer create --from acc_001 --to acc_002 --amount 10000 --dry-run
bank transfer status transfer_001

# Ledger
bank ledger audit
bank ledger entries

# Debug
bank logs tail --lines 50
bank debug tx tx_001

# Output formats
bank user list --json
bank tx list --user user_001 | jq '.'
```

---

## 📋 Commands

### 🔐 **Authentication**

| Command                                                 | Description                     |
| ------------------------------------------------------- | ------------------------------- |
| `bank auth login --email <email> --password <password>` | Login and store token           |
| `bank auth logout`                                      | Clear stored credentials        |
| `bank auth whoami`                                      | Show current authenticated user |

### 👤 **Users**

| Command              | Description    |
| -------------------- | -------------- |
| `bank user list`     | List all users |
| `bank user get <id>` | Get user by ID |

### 💳 **Accounts**

| Command                         | Description              |
| ------------------------------- | ------------------------ |
| `bank account list --user <id>` | List accounts for a user |
| `bank account get <id>`         | Get account by ID        |

### 💸 **Transactions**

| Command                             | Description            |
| ----------------------------------- | ---------------------- |
| `bank tx list --user <id> [--json]` | List user transactions |
| `bank tx get <id>`                  | Get transaction by ID  |
| `bank tx simulate`                  | Simulate a transaction |

### 🔄 **Transfers**

| Command                                                                   | Description         |
| ------------------------------------------------------------------------- | ------------------- |
| `bank transfer create --from <id> --to <id> --amount <value> [--dry-run]` | Create transfer     |
| `bank transfer status <id>`                                               | Get transfer status |

### 📒 **Ledger**

| Command               | Description                |
| --------------------- | -------------------------- |
| `bank ledger audit`   | Run ledger integrity audit |
| `bank ledger entries` | List ledger entries        |

### 🔍 **Debug**

| Command                        | Description                  |
| ------------------------------ | ---------------------------- |
| `bank logs tail [--lines <n>]` | Tail recent system logs      |
| `bank debug tx <id>`           | Debug transaction processing |

### ℹ️ **Other**

| Command          | Description           |
| ---------------- | --------------------- |
| `bank --version` | Show CLI version      |
| `bank --help`    | Show help information |

---

## 🛠️ Tech Stack

### ⚙️ **Core Framework**

```
Go 1.25+ + Cobra + Viper
├── 🔐 Cobra (CLI Framework)
│   ├── Command Management
│   ├── Flag Parsing
│   └── Help Generation
├── 📦 Viper (Configuration)
│   ├── YAML Config Files
│   ├── Environment Variables
│   └── Key-Value Store
├── 🧪 Testify (Testing)
│   ├── Mocking
│   └── Assertions
└── 📊 go-pretty (Table Output)
    ├── Formatted Tables
    └── Human-Readable Output
```

### 🏗️ **Architecture**

```
CLI Package Structure
├── cmd/                    # Command definitions (Cobra)
│   ├── auth/             # Authentication commands
│   ├── user/            # User management commands
│   ├── account/         # Account commands
│   ├── tx/              # Transaction commands
│   ├── transfer/        # Transfer commands
│   ├── ledger/          # Ledger commands
│   └── debug/           # Debug commands
├── internal/              # Internal packages
│   ├── config/          # Configuration management
│   ├── authstore/       # Token storage
│   ├── output/          # Output formatting (JSON/Table)
│   └── client/          # API client (Mock + Real)
├── main.go                # Entry point
└── go.mod                # Go modules
```

---

## 📁 Architecture

### 🏗️ **Package Structure**

```
package/cli/
├── cmd/                        # Cobra Commands
│   ├── root.go                 # Root command
│   ├── version.go              # Version command
│   ├── auth/
│   │   └── auth.go             # auth login/logout/whoami
│   ├── user/
│   │   └── user.go             # user list/get
│   ├── account/
│   │   └── account.go          # account list/get
│   ├── tx/
│   │   └── tx.go               # tx list/get/simulate
│   ├── transfer/
│   │   └── transfer.go         # transfer create/status
│   ├── ledger/
│   │   └── ledger.go           # ledger audit/entries
│   └── debug/
│       └── debug.go            # logs tail, debug tx
├── internal/                    # Internal Packages
│   ├── config/
│   │   └── config.go           # Viper configuration
│   ├── authstore/
│   │   └── auth.go             # Token management
│   ├── output/
│   │   └── output.go           # JSON/Table output
│   └── client/
│       └── client.go           # API client interface + Mock
├── scripts/
│   └── build.sh                # Multi-platform build script
├── config.sample.yaml          # Sample configuration
├── install.sh                  # Local installation script
└── README.md                   # This file
```

### 🔄 **Data Flow**

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   User Input     │───▶│   Cobra CLI      │───▶│   API Client    │
│   (Commands)    │    │   (Validation)   │    │   (Mock/Real)   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │                       │
                              ▼                       ▼
                       ┌──────────────────┐    ┌─────────────────┐
                       │   Output Formatter│    │   Banking API    │
                       │   (JSON/Table)    │◀───│   (Backend)     │
                       └──────────────────┘    └─────────────────┘
```

---

## 🔧 Development

### 🛠️ **Build Commands**

```bash
# Build binary
go build -o bank .

# Multi-platform build
chmod +x scripts/build.sh
./scripts/build.sh v1.0.0

# Run tests
go test ./... -v

# Run with mock data
./bank user list
./bank tx list --user user_001 --json
```

### 📋 **Development Guidelines**

- **Cobra Commands** - All commands in `cmd/` package
- **Internal Packages** - Business logic in `internal/` package
- **Interface-Based Client** - API client implements `Client` interface
- **Mock Client** - Use `MockClient` for testing without backend
- **Output Abstraction** - Use `Output` struct for formatting

### 🔐 **Security Considerations**

- Token storage in `~/.bank/token.json` with 0600 permissions
- No hardcoded credentials
- Environment-based configuration
- Audit logging for sensitive operations

---

## 🗺️ Development Roadmap

### 🎯 **Phase 1: Core CLI (✅ Complete - V1)**

- ✅ **Authentication** - Login/logout with token management
- ✅ **Basic Commands** - User, account, transaction commands
- ✅ **Transfer Support** - Create transfers with dry-run
- ✅ **Ledger Commands** - Audit and entries
- ✅ **Debug Tools** - Log tailing and transaction debugging
- ✅ **Configuration** - Multi-environment support
- ✅ **Output Formats** - JSON and table output
- ✅ **Unit Tests** - Mock client testing

### 🚀 **Phase 2: Integration (🔄 In Progress)**

- 🔄 **Real API Client** - Connect to actual banking backend
- 🔄 **Authentication Enhancement** - OAuth2/JWT integration
- 🔄 **Advanced Filtering** - Complex query parameters
- 🔄 **Shell Completion** - Bash/Zsh/Fish support

### 🌟 **Phase 3: Enterprise Features (📋 Planned)**

- 📋 **Batch Operations** - Bulk transaction processing
- 📋 **Audit Logging** - Comprehensive operation logging
- 📋 **Interactive Mode** - REPL for exploratory operations
- 📋 **Plugin System** - Extensible command architecture

---

## 🤝 Contributing

We welcome contributions from internal developers and employees!

### 🎯 **How to Get Started**

1. **Fork and clone** the repository
2. **Create a branch** for your feature or fix
3. **Implement changes** following the architecture guidelines
4. **Add tests** for new functionality
5. **Submit a pull request** with clear description

### 🏗️ **Areas Needing Help**

- **Real API Client** - Backend integration
- **Authentication** - OAuth2/JWT implementation
- **Testing** - Comprehensive test coverage
- **Documentation** - API docs and usage guides
- **Shell Completion** - Autocompletion support

---

## 📊 Project Status

| Component                | Status      | Technology         | Notes                       |
| ------------------------ | ----------- | ------------------ | --------------------------- |
| **CLI Framework**        | ✅ Complete | Go + Cobra + Viper | Full command implementation |
| **Authentication**       | ✅ Complete | Token-based        | Login/logout/whoami         |
| **User Management**      | ✅ Complete | Mock Client        | List and get operations     |
| **Account Operations**   | ✅ Complete | Mock Client        | List and get operations     |
| **Transactions**         | ✅ Complete | Mock Client        | List, get, simulate         |
| **Transfers**            | ✅ Complete | Mock Client        | Create with dry-run         |
| **Ledger**               | ✅ Complete | Mock Client        | Audit and entries           |
| **Debug Tools**          | ✅ Complete | Mock Client        | Logs and transaction debug  |
| **Output Formatting**    | ✅ Complete | go-pretty          | JSON and table formats      |
| **Multi-Environment**    | ✅ Complete | Viper              | Staging and production      |
| **Unit Tests**           | ✅ Complete | Testify            | Mock client testing         |
| **Real API Integration** | 📋 Planned  | Go                 | Backend connection          |
| **Shell Completion**     | 📋 Planned  | Cobra              | Autocompletion support      |
| **Batch Operations**     | 📋 Planned  | Go                 | Bulk processing             |

---

## 🏆 Sponsors & Partners

**Development led by [Sky Genesis Enterprise](https://skygenesisenterprise.com)**

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2025 Sky Genesis Enterprise

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

---

## 🙏 Acknowledgments

- **Sky Genesis Enterprise** - Project leadership
- **Go Community** - High-performance programming language
- **Cobra Team** - Excellent CLI framework
- **Viper Team** - Configuration management
- **Testify Team** - Testing utilities

---

<div align="center">

### 💰 **Aether Bank CLI - Secure Banking Operations at Your Fingertips**

[⭐ Star This Repo](https://github.com/skygenesisenterprise/aether-bank) • [🐛 Report Issues](https://github.com/skygenesisenterprise/aether-bank/issues) • [📖 Documentation](https://github.com/skygenesisenterprise/aether-bank/wiki)

---

**Built with ❤️ by the [Sky Genesis Enterprise](https://skygenesisenterprise.com) team**

_Enterprise-grade CLI for secure and scriptable banking operations_

</div>
