.PHONY: producer


export SHELL := /bin/sh


.ONESHELL:
dev:
	# Use minikube cluster docker environment
	eval $$(minikube -p minikube docker-env); \
	# fd -e go | entr -rc docker build . --tag "services:latest"; echo "hello world"
	fd -e go | entr -rc sh -c '\
			docker build . --tag "services:latest"; \
			kubectl delete -f producer.yaml \
			kubectl apply -f producer.yaml \
	'

producer:
	fd -e go | entr -rc go run  cmd/producer/main.go --migrate
