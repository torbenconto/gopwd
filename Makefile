GOFILES_BUILD             := $(shell find . -type f -name '*.go' -not -name '*_test.go')
GOPWD_OUTPUT              ?= gopwd
GOPWD_REVISION            := $(shell cat COMMIT 2>/dev/null || git rev-parse --short=8 HEAD)
BASH_COMPLETION_OUTPUT    := bash.completion
FISH_COMPLETION_OUTPUT    := fish.completion
ZSH_COMPLETION_OUTPUT     := zsh.completion
DATE                      := $(shell date -u '+%FT%T%z')
GO                        ?= GO111MODULE=on CGO_ENABLED=0 go
PREFIX                    ?= /usr
BINDIR                    ?= $(PREFIX)/bin

all: build completion

build: $(GOPWD_OUTPUT)

$(GOPWD_OUTPUT): $(GOFILES_BUILD)
	@echo -n ">> BUILD, version = $(GOPWD_VERSION)/$(GOPWD_REVISION), output = $@"
	@$(GO) build -o $@
	@echo " [OK]"

completion: $(BASH_COMPLETION_OUTPUT) $(FISH_COMPLETION_OUTPUT) $(ZSH_COMPLETION_OUTPUT)

%.completion:
	@echo ">> $* completion, output = $@"
	@./$(GOPWD_OUTPUT) completion $* > $@
	@echo " [OK]"

clean:
	@echo -n ">> CLEAN"
	@$(GO) clean -i ./...
	@rm -f $(GOPWD_OUTPUT) $(BASH_COMPLETION_OUTPUT) $(FISH_COMPLETION_OUTPUT) $(ZSH_COMPLETION_OUTPUT)
	@echo " [OK]"

install-completion:
	@sudo install -d $(PREFIX)/share/zsh/site-functions $(PREFIX)/share/bash-completion/completions $(PREFIX)/share/fish/vendor_completions.d
	@sudo install -m 0644 $(ZSH_COMPLETION_OUTPUT) $(PREFIX)/share/zsh/site-functions/_gopwd
	@sudo install -m 0644 $(BASH_COMPLETION_OUTPUT) $(PREFIX)/share/bash-completion/completions/gopwd
	@sudo install -m 0644 $(FISH_COMPLETION_OUTPUT) $(PREFIX)/share/fish/vendor_completions.d/gopwd.fish
	@printf '%s\n' '$(OK)'

install: build install-completion
	@sudo install -d $(BINDIR)
	@sudo install -m 0755 $(GOPWD_OUTPUT) $(BINDIR)/gopwd
	@printf '%s\n' '$(OK)'
