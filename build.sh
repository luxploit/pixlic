#!/bin/bash

set -e

PROJECT_NAME="pixlic"

# Usage:
# ./build.sh [-a <arch> -p <platform>]

ARCH="native"            
PLATFORM="native"

while [[ $# -gt 0 ]]; do
  case "$1" in
    -a)
      ARCH="$2"
      shift 2
      ;;
    -p)
      PLATFORM="$2"
      shift 2
      ;;
    -h)
      echo "Usage: $0 [-a <arch> -p <platform>]"
      exit 1
      ;;
    *)
      shift
      ;;
  esac
done

# Determine file extension based on OS
EXT=".lxb"
if uname | grep -iq 'mingw\|msys\|cygwin'; then
  EXT=".exe"
fi

# Set GOOS and GOARCH if not native
ENV_VARS=""
if [ "$ARCH" != "native" ]; then
  ENV_VARS+="GOARCH=$ARCH "
fi
if [ "$PLATFORM" != "native" ]; then
  ENV_VARS+="GOOS=$PLATFORM "
fi

if [ "$PLATFORM" == "linux" ]; then
  EXT=".lxb" # manual override for cross-compilation
fi

# Output binary info
BUILD_DIR="bin"
mkdir -p "$BUILD_DIR"

BIN_NAME="${PROJECT_NAME}${EXT}"
BIN_PATH="$BUILD_DIR/$BIN_NAME"

# Build
echo "Building $PROJECT_NAME..."
eval "${ENV_VARS}go build -o \"$BIN_PATH\""
echo "Build complete: $BIN_PATH"