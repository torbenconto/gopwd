#!/bin/bash

if ! command -v go &> /dev/null
then
    echo "Go could not be found. Please install Go and retry."
    exit 1
fi

git clone https://github.com/torbenconto/gopwd.git
if [ $? -ne 0 ]; then
    echo "Failed to clone the repository."
    exit 1
fi

cd gopwd || exit

go build -o gopwd
if [ $? -ne 0 ]; then
    echo "Failed to build the CLI."
    exit 1
fi

sudo mv gopwd /usr/local/bin/
if [ $? -ne 0 ]; then
    echo "Failed to move the binary."
    exit 1
fi

sudo chmod +x /usr/local/bin/gopwd

SHELL_TYPE=$(basename "$SHELL")

case $SHELL_TYPE in
bash)
    gopwd completion bash > gopwd_autocomplete.sh
    echo "source $(pwd)/gopwd_autocomplete.sh" >> ~/.bashrc
    echo "Autocomplete script added to .bashrc. Please restart your shell or source your .bashrc to activate."
    ;;
zsh)
    gopwd completion zsh > gopwd_autocomplete.sh
    echo "source $(pwd)/gopwd_autocomplete.sh" >> ~/.zshrc
    echo "Autocomplete script added to .zshrc. Please restart your shell or source your .zshrc to activate."
    ;;
fish)
    gopwd completion fish > gopwd_autocomplete.fish
    echo "source $(pwd)/gopwd_autocomplete.fish" >> ~/.config/fish/config.fish
    echo "Autocomplete script added to config.fish. Please restart your shell or source your config.fish to activate."
    ;;
*)
    echo "Autocomplete setup is not supported for $SHELL_TYPE shell automatically. Please refer to the documentation for manual setup."
    ;;
esac