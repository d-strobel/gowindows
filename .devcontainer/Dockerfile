FROM mcr.microsoft.com/devcontainers/base:trixie

ARG GO_VERSION
ARG GO_ARCH

RUN apt-get update && apt-get install -y \
    fish \
    python3 \
    python3-pip \
    python3-venv \
    && rm -rf /var/lib/apt/lists/*

# Install Go
RUN wget --quiet https://go.dev/dl/go${GO_VERSION}.${GO_ARCH}.tar.gz \
    && tar -C /usr/local/ -xzf go${GO_VERSION}.${GO_ARCH}.tar.gz

# Switch to user
USER vscode

# Create virtual environment
RUN python3 -m venv /home/vscode/venv

# Add venv and go to PATH
ENV PATH="/home/vscode/venv/bin:/usr/local/go/bin:${PATH}"

# Install Python packages into the venv
RUN pip install --upgrade pip && \
    pip install pre-commit
