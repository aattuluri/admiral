#!/bin/bash

[ $# -lt 1 ] && { echo "Usage: $0 <k8s_version>" ; exit 1; }

k8s_version=$1

echo "Creating K8s cluster with version: $k8s_version"

if [[ $IS_LOCAL == "false" ]]; then
  sudo -E minikube start --driver=none --cpus 4 --memory 4096 --kubernetes-version=$k8s_version &> $HOME/minikube.log 2>&1 < /dev/null

else
  minikube start --memory=4096 --cpus=4 --kubernetes-version=$k8s_version --vm-driver "virtualbox"
  #label node for locality load balancing
  kubectl label nodes minikube --overwrite failure-domain.beta.kubernetes.io/region=us-west-2
fi