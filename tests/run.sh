#!/bin/bash

set -o errexit
set -o pipefail

[ $# -lt 3 ] && { echo "Usage: $0 <k8s_version> <istio_version> <admiral_install_dir>" ; exit 1; }

k8s_version=$1
istio_version=$2
install_dir=$3
os=""
source ./create_cluster.sh $k8s_version
# Uncomment below line if setup fails due to KUBECONFIG not set
export KUBECONFIG=~/.kube/config
#label node for locality load balancing
kubectl label nodes minikube failure-domain.beta.kubernetes.io/region=us-west-2
# change $os from "linux" to "osx" when running on local computer
if [[ "$OSTYPE" == "darwin"* ]]; then
  os="osx"
else
  if [[ $istio_version == "1.5"* ]]; then
    os="linux"
  else
    os="linux-amd64"
  fi
fi
./install_istio.sh $istio_version $os
./dns_setup.sh $install_dir
$install_dir/scripts/install_admiral.sh $install_dir
$install_dir/scripts/install_rollouts.sh
$install_dir/scripts/cluster-secret.sh $KUBECONFIG  $KUBECONFIG admiral
$install_dir/scripts/install_sample_services.sh $install_dir

#allow for stabilization of DNS and Istio SEs before running tests
sleep 10
./test1.sh "webapp" "sample" "greeting"
./test2.sh "webapp" "sample-rollout-bluegreen" "greeting"
#cleanup to fee up the pipeline minkube resources
if [[ $IS_LOCAL == "false" ]]; then
  kubectl delete ns sample-rollout-bluegreen
fi
./test3.sh "grpc-client" "sample" "grpc-server" $install_dir
./test4.sh "webapp" "sample"

./cleanup.sh $istio_version