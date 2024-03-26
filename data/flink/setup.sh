#! /bin/bash

# Add the Apache Flink k8s operator helm repository
helm repo add flink-operator-repo https://downloads.apache.org/flink/flink-kubernetes-operator-1.8.0/

# Install certificate manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.4/cert-manager.yaml

# Wait for cert manager
kubectl wait --for=condition=Ready pods --all --namespace=cert-manager --timeout=300s

# Install Flink k8s operator
helm install flink-kubernetes-operator flink-operator-repo/flink-kubernetes-operator
