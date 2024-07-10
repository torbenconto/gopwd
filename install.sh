#!/bin/bash

# Define the repository URL and a temporary directory for cloning
REPO_URL="https://github.com/torbenconto/gopwd.git"
TEMP_DIR=$(mktemp -d)

# Clone the repository
git clone $REPO_URL $TEMP_DIR

# Change to the repository directory
cd $TEMP_DIR

# Run make install
make install

# Cleanup: Go back to the original directory and remove the temporary directory
cd -
rm -rf $TEMP_DIR

# Completion message
echo "Installation completed successfully."