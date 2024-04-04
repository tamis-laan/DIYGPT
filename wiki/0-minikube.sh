#! /bin/bash

# Start default cluster
minikube start

# Mount services
minikube mount ./services:/services

