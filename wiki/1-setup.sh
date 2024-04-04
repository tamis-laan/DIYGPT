#! /bin/bash

# Make sure system is running properly
kubectl wait --for=condition=Ready pods --all --namespace=kube-system --timeout=300s

# Add kafka operator
helm repo add strimzi http://strimzi.io/charts

# Add the Apache Flink k8s operator helm repository
helm repo add flink-operator-repo https://downloads.apache.org/flink/flink-kubernetes-operator-1.8.0/

# Add the Apache Spark k8s operator helm repository
helm repo add spark-operator https://kubeflow.github.io/spark-operator

# Update 
helm repo update strimzi flink-operator-repo spark-operator

# Install certificate manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.4/cert-manager.yaml

# Install kafka operator into cluster
helm install kafka-operator strimzi/strimzi-kafka-operator \
	--namespace operators --create-namespace --version 0.40.0\
	--set watchAnyNamespace=true

# Wait for cert manager
kubectl wait --for=condition=Ready pods --all --namespace=cert-manager --timeout=300s

# Install Flink k8s operator
helm install flink-kubernetes-operator flink-operator-repo/flink-kubernetes-operator \
	--namespace operators --create-namespace --version 1.8.0

# Install spark operator into cluster
helm install spark spark-operator/spark-operator \
	--namespace operators --create-namespace --version v1beta2-1.3.8-3.1.1

# Install Minio
kubectl apply -f k8s/minio-dev.yaml

# Install postgresql
kubectl apply -f k8s/postgresql-dev.yaml

# Install kafka data bus cluster
kubectl apply -f k8s/kafka-clusters.yaml

# Wait for cert manager
kubectl wait --for=condition=Ready pods --all --namespace=kafka --timeout=300s

# Install kafka topics
kubectl apply -f k8s/kafka-topics.yaml
