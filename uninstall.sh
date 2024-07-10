#!/bin/bash

# Step 1: Remove the Binary
sudo rm /usr/local/bin/gopwd

# Step 2: Remove Autocomplete Script
SHELL_TYPE=$(basename "$SHELL")

case $SHELL_TYPE in
bash)
    sed -i '/gopwd_autocomplete.sh/d' ~/.bashrc
    ;;
zsh)
    sed -i '/gopwd_autocomplete.sh/d' ~/.zshrc
    ;;
fish)
    sed -i '/gopwd_autocomplete.fish/d' ~/.config/fish/config.fish
    ;;
*)
    echo "Autocomplete cleanup not supported for $SHELL_TYPE shell automatically. Please manually remove any gopwd autocomplete setup."
    ;;
esac

# Optional: Step 3: Cleanup additional files/directories
# Replace <path_to_additional_files> with the actual path if applicable
# rm -rf <path_to_additional_files>

echo "gopwd uninstalled successfully."