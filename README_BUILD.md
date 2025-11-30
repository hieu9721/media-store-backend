# Media Store Backend - Quick Build Guide

## 🚀 Quick Start

### Windows
```cmd
build.bat all
```

### Linux/macOS
```bash
chmod +x build.sh
./build.sh all
```

### Cross-platform (với Make)
```bash
make build-all
```

### Single Platform Build
```bash
# Current platform
go build -o media-store-backend

# Windows
GOOS=windows GOARCH=amd64 go build -o media-store-backend.exe

# Linux
GOOS=linux GOARCH=amd64 go build -o media-store-backend

# macOS
GOOS=darwin GOARCH=amd64 go build -o media-store-backend
```

---

## 📦 Build Files Created

Dự án bao gồm các build scripts sau:

1. **build.sh** - Bash script cho Linux/macOS/Git Bash
2. **build.bat** - Batch script cho Windows CMD/PowerShell
3. **Makefile** - Universal build system (khuyến nghị)
4. **Dockerfile** - Container build
5. **docker-compose.yml** - Full stack với MongoDB

---

## 🎯 Supported Platforms

Build scripts tự động build cho:

### Windows
- AMD64 (64-bit)
- 386 (32-bit)
- ARM64

### Linux
- AMD64 (64-bit)
- 386 (32-bit)
- ARM64
- ARM (32-bit)

### macOS
- AMD64 (Intel)
- ARM64 (Apple Silicon)

---

## 📋 Build Commands Reference

### Using build.sh (Linux/macOS)
```bash
./build.sh all          # All platforms
./build.sh windows      # Windows only
./build.sh linux        # Linux only
./build.sh mac          # macOS only
./build.sh arm          # ARM platforms only
```

### Using build.bat (Windows)
```cmd
build.bat all           # All platforms
build.bat windows       # Windows only
build.bat linux         # Linux only
build.bat mac           # macOS only
build.bat arm           # ARM platforms only
```

### Using Makefile
```bash
make build              # Current platform
make build-all          # All platforms
make windows            # Windows builds
make linux              # Linux builds
make mac                # macOS builds
make arm                # ARM builds
make clean              # Clean builds
make help               # Show all commands
```

---

## 🐳 Docker Build

### Build and Run
```bash
docker-compose up --build
```

### Manual Docker Build
```bash
docker build -t media-store-backend .
docker run -p 8080:8080 media-store-backend
```

---

## 📊 Build Output

Artifacts được tạo trong thư mục `build/`:

```
build/
├── media-store-backend-windows-amd64.exe
├── media-store-backend-windows-386.exe
├── media-store-backend-windows-arm64.exe
├── media-store-backend-linux-amd64
├── media-store-backend-linux-386
├── media-store-backend-linux-arm64
├── media-store-backend-linux-arm
├── media-store-backend-darwin-amd64
└── media-store-backend-darwin-arm64
```

---

## ⚙️ Environment Setup

### Required
- Go 1.21 or higher
- Git (for build.sh on Windows)

### Optional
- Make (for Makefile)
- Docker (for containerized builds)

### Install Make on Windows
```powershell
# Using Chocolatey
choco install make

# Using Scoop
scoop install make
```

---

## 🔧 Development Tools

### Generate JWT Secret
```bash
# Using Makefile
make gen-secret

# Using script
go run scripts/generate-jwt-secret.go
```

### Install Dev Tools
```bash
make install-tools
```

### Run with Auto-reload
```bash
make dev  # Requires air
```

---

## 📚 Full Documentation

- **[BUILD_GUIDE.md](BUILD_GUIDE.md)** - Comprehensive build documentation
- **[API_SUMMARY.md](API_SUMMARY.md)** - API reference
- **[UPLOAD_API_GUIDE.md](UPLOAD_API_GUIDE.md)** - Upload endpoints
- **[CHANGELOG.md](CHANGELOG.md)** - Version history

---

## 🎓 Examples

### Build for Production
```bash
# Build all platforms with version
VERSION=1.0.0 make build-all

# Create release archives
make release
```

### Development Build
```bash
# Quick build for testing
make build

# Run immediately
make run
```

### Docker Production
```bash
# Build and run full stack
docker-compose up -d

# View logs
docker-compose logs -f backend

# Stop
docker-compose down
```

---

## ⚡ Quick Tips

1. **Use Makefile when possible** - Most flexible and feature-rich
2. **build.sh for Linux/macOS** - Native shell scripting
3. **build.bat for Windows** - When Make is not available
4. **Docker for deployment** - Consistent environment

---

## 🐛 Troubleshooting

### Can't execute build.sh
```bash
chmod +x build.sh
```

### Make not found (Windows)
```cmd
# Use build.bat instead
build.bat all
```

### Build fails
```bash
# Clean and rebuild
make clean
make deps
make build
```
