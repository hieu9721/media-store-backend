#!/bin/bash

# Media Store Backend - Multi-platform Build Script
# Supports: Windows, Linux, macOS, ARM

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# App info
APP_NAME="media-store-backend"
VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="build"
LDFLAGS="-s -w -X main.Version=${VERSION}"

echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     Media Store Backend - Multi-platform Build Script         ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Clean build directory
echo -e "${YELLOW}🧹 Cleaning build directory...${NC}"
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Function to build for a platform
build_platform() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT_NAME=$3

    echo -e "${BLUE}📦 Building for ${GOOS}/${GOARCH}...${NC}"

    GOOS=${GOOS} GOARCH=${GOARCH} go build \
        -ldflags="${LDFLAGS}" \
        -o ${BUILD_DIR}/${OUTPUT_NAME} \
        .

    if [ $? -eq 0 ]; then
        SIZE=$(ls -lh ${BUILD_DIR}/${OUTPUT_NAME} | awk '{print $5}')
        echo -e "${GREEN}✓ Successfully built ${OUTPUT_NAME} (${SIZE})${NC}"
    else
        echo -e "${RED}✗ Failed to build ${OUTPUT_NAME}${NC}"
        return 1
    fi
}

# Parse arguments
BUILD_ALL=false
BUILD_WINDOWS=false
BUILD_LINUX=false
BUILD_MAC=false
BUILD_ARM=false

if [ $# -eq 0 ]; then
    BUILD_ALL=true
else
    for arg in "$@"; do
        case $arg in
            all)
                BUILD_ALL=true
                ;;
            windows)
                BUILD_WINDOWS=true
                ;;
            linux)
                BUILD_LINUX=true
                ;;
            mac|darwin)
                BUILD_MAC=true
                ;;
            arm)
                BUILD_ARM=true
                ;;
            *)
                echo -e "${RED}Unknown platform: $arg${NC}"
                echo "Usage: ./build.sh [all|windows|linux|mac|arm]"
                exit 1
                ;;
        esac
    done
fi

echo -e "${YELLOW}Version: ${VERSION}${NC}"
echo ""

# Build for all platforms or specific ones
if [ "$BUILD_ALL" = true ] || [ "$BUILD_WINDOWS" = true ]; then
    echo -e "${BLUE}=== Windows Builds ===${NC}"
    build_platform "windows" "amd64" "${APP_NAME}-windows-amd64.exe"
    build_platform "windows" "386" "${APP_NAME}-windows-386.exe"
    build_platform "windows" "arm64" "${APP_NAME}-windows-arm64.exe"
    echo ""
fi

if [ "$BUILD_ALL" = true ] || [ "$BUILD_LINUX" = true ]; then
    echo -e "${BLUE}=== Linux Builds ===${NC}"
    build_platform "linux" "amd64" "${APP_NAME}-linux-amd64"
    build_platform "linux" "386" "${APP_NAME}-linux-386"
    build_platform "linux" "arm64" "${APP_NAME}-linux-arm64"
    build_platform "linux" "arm" "${APP_NAME}-linux-arm"
    echo ""
fi

if [ "$BUILD_ALL" = true ] || [ "$BUILD_MAC" = true ]; then
    echo -e "${BLUE}=== macOS Builds ===${NC}"
    build_platform "darwin" "amd64" "${APP_NAME}-darwin-amd64"
    build_platform "darwin" "arm64" "${APP_NAME}-darwin-arm64"
    echo ""
fi

if [ "$BUILD_ARM" = true ]; then
    echo -e "${BLUE}=== ARM Builds ===${NC}"
    build_platform "linux" "arm" "${APP_NAME}-linux-arm"
    build_platform "linux" "arm64" "${APP_NAME}-linux-arm64"
    echo ""
fi

# Summary
echo -e "${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                    Build Complete!                             ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}Build artifacts location: ${BUILD_DIR}/${NC}"
echo ""
ls -lh ${BUILD_DIR}
echo ""
echo -e "${BLUE}Usage examples:${NC}"
echo -e "  Windows:  ${BUILD_DIR}/${APP_NAME}-windows-amd64.exe"
echo -e "  Linux:    ${BUILD_DIR}/${APP_NAME}-linux-amd64"
echo -e "  macOS:    ${BUILD_DIR}/${APP_NAME}-darwin-amd64"
echo ""
