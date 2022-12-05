#!/bin/sh

echo "starting k3s..."
cd /tmp && nohup bash -c 'sudo k3s server --snapshotter=fuse-overlayfs --write-kubeconfig-mode 644 &' > k3s.log