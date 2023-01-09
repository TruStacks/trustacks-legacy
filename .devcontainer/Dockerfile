ARG VARIANT="1.18-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

RUN apt update && apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    libarchive-tools \
    lsb-release

# install sops.
RUN curl -LO https://github.com/mozilla/sops/releases/download/v3.7.3/sops-v3.7.3.linux.amd64 && \
    sudo install sops-v3.7.3.linux.amd64 /usr/local/bin/sops && \
    rm -f sops-v3.7.3.linux.amd64

# install k3s
RUN echo -e 'token: "secret"\ndebug: true' && \
    curl -sfL https://get.k3s.io | INSTALL_K3S_SKIP_ENABLE=true sh - && \
    apt install fuse-overlayfs -y

# configure k3s iptables
RUN apt remove iptables nftables -y
ENV PATH="/var/lib/rancher/k3s/data/current/bin/:/var/lib/rancher/k3s/data/current/bin/aux:$PATH"

# copy k3s traefik config
COPY traefik-config.yaml /var/lib/rancher/k3s/server/manifests/traefik-config.yaml

# install kubectl & helm / kelm-s3
RUN curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash && \
    sudo ln -s $(which bsdtar) /usr/local/bin/tar
ENV KUBECONFIG=/etc/rancher/k3s/k3s.yaml

# install okteto
RUN curl https://get.okteto.com -sSfL | sh

COPY post-start.sh /usr/local/bin/post-start

USER vscode

# install the helm s3 plugin
RUN alias tar=bsdtar; helm plugin install https://github.com/hypnoglow/helm-s3.git