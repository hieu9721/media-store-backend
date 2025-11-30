@echo off
REM Media Store Backend - Multi-platform Build Script for Windows
REM Supports: Windows, Linux, macOS, ARM

setlocal enabledelayedexpansion

set APP_NAME=media-store-backend
set VERSION=1.0.0
set BUILD_DIR=build
set LDFLAGS=-s -w

echo ================================================================
echo      Media Store Backend - Multi-platform Build Script
echo ================================================================
echo.

REM Clean build directory
echo Cleaning build directory...
if exist %BUILD_DIR% rmdir /s /q %BUILD_DIR%
mkdir %BUILD_DIR%

REM Parse arguments
set BUILD_ALL=false
set BUILD_WINDOWS=false
set BUILD_LINUX=false
set BUILD_MAC=false
set BUILD_ARM=false

if "%~1"=="" (
    set BUILD_ALL=true
) else (
    for %%a in (%*) do (
        if "%%a"=="all" set BUILD_ALL=true
        if "%%a"=="windows" set BUILD_WINDOWS=true
        if "%%a"=="linux" set BUILD_LINUX=true
        if "%%a"=="mac" set BUILD_MAC=true
        if "%%a"=="darwin" set BUILD_MAC=true
        if "%%a"=="arm" set BUILD_ARM=true
    )
)

echo Version: %VERSION%
echo.

REM Windows builds
if "%BUILD_ALL%"=="true" set BUILD_WINDOWS=true
if "%BUILD_WINDOWS%"=="true" (
    echo === Windows Builds ===
    call :build_platform windows amd64 %APP_NAME%-windows-amd64.exe
    call :build_platform windows 386 %APP_NAME%-windows-386.exe
    call :build_platform windows arm64 %APP_NAME%-windows-arm64.exe
    echo.
)

REM Linux builds
if "%BUILD_ALL%"=="true" set BUILD_LINUX=true
if "%BUILD_LINUX%"=="true" (
    echo === Linux Builds ===
    call :build_platform linux amd64 %APP_NAME%-linux-amd64
    call :build_platform linux 386 %APP_NAME%-linux-386
    call :build_platform linux arm64 %APP_NAME%-linux-arm64
    call :build_platform linux arm %APP_NAME%-linux-arm
    echo.
)

REM macOS builds
if "%BUILD_ALL%"=="true" set BUILD_MAC=true
if "%BUILD_MAC%"=="true" (
    echo === macOS Builds ===
    call :build_platform darwin amd64 %APP_NAME%-darwin-amd64
    call :build_platform darwin arm64 %APP_NAME%-darwin-arm64
    echo.
)

REM ARM builds
if "%BUILD_ARM%"=="true" (
    echo === ARM Builds ===
    call :build_platform linux arm %APP_NAME%-linux-arm
    call :build_platform linux arm64 %APP_NAME%-linux-arm64
    echo.
)

echo ================================================================
echo                     Build Complete!
echo ================================================================
echo.
echo Build artifacts location: %BUILD_DIR%\
echo.
dir /b %BUILD_DIR%
echo.
echo Usage examples:
echo   Windows:  %BUILD_DIR%\%APP_NAME%-windows-amd64.exe
echo   Linux:    %BUILD_DIR%\%APP_NAME%-linux-amd64
echo   macOS:    %BUILD_DIR%\%APP_NAME%-darwin-amd64
echo.
goto :eof

:build_platform
set GOOS=%~1
set GOARCH=%~2
set OUTPUT_NAME=%~3

echo Building for %GOOS%/%GOARCH%...

set CGO_ENABLED=0
go build -ldflags="%LDFLAGS%" -o %BUILD_DIR%\%OUTPUT_NAME% .

if %errorlevel% equ 0 (
    echo [OK] Successfully built %OUTPUT_NAME%
) else (
    echo [ERROR] Failed to build %OUTPUT_NAME%
)

goto :eof
