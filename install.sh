#!/bin/bash
# This script automatically downloads the latest release of sysx from GitHub,
# determines the appropriate binary for your CPU architecture,
# and installs it to /usr/local/bin so it's available in your PATH.

API_URL="https://api.github.com/repos/krau/sysx/releases/latest"

# Fetch the latest release tag (e.g., "v0.1.4") from GitHub.
echo "Fetching latest release info from $API_URL ..."
TAG=$(curl -s "$API_URL" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$TAG" ]; then
    echo "Error: Unable to fetch the latest release tag."
    exit 1
fi
echo "Latest release tag: $TAG"

# Determine the current system's CPU architecture.
ARCH=$(uname -m)
# Map the architecture to the naming convention used in the release files.
case "$ARCH" in
x86_64)
    ARCH="amd64"
    ;;
aarch64)
    ARCH="arm64"
    ;;
*)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac
echo "Detected architecture: $ARCH"

# Construct the download URL for the tarball using the tag and architecture.
# The file name format is: sysx-<tag>-linux-<arch>.tar.gz (e.g., sysx-v0.1.4-linux-arm64.tar.gz)
DOWNLOAD_URL="https://github.com/krau/sysx/releases/download/${TAG}/sysx-${TAG}-linux-${ARCH}.tar.gz"
echo "Download URL: $DOWNLOAD_URL"

# Create a temporary directory for downloading and extracting the tarball.
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR" || {
    echo "Error: Could not change directory to $TMP_DIR"
    exit 1
}

# Download the tarball using curl. The -L flag follows any redirects.
echo "Downloading sysx tarball..."
if ! curl -L -o sysx.tar.gz "$DOWNLOAD_URL"; then
    echo "Error: Download failed."
    exit 1
fi

# Extract the downloaded tarball.
echo "Extracting sysx tarball..."
if ! tar -xzf sysx.tar.gz; then
    echo "Error: Extraction failed."
    exit 1
fi

# Verify that the sysx binary exists in the extracted contents.
if [ ! -f sysx ]; then
    echo "Error: sysx binary not found after extraction."
    exit 1
fi

# Ensure the sysx binary is executable.
chmod +x sysx

# Install the sysx binary to /usr/local/bin, which should be in the user's PATH.
# This step requires sudo privileges.
echo "Installing sysx to /usr/local/bin ..."
if ! sudo mv sysx /usr/local/bin/; then
    echo "Error: Installation failed."
    exit 1
fi

# Clean up by removing the temporary directory.
cd /
rm -rf "$TMP_DIR"

echo "sysx has been installed successfully and is available in your PATH."
