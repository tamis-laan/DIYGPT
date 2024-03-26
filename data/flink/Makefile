
build-local:
	docker build . --tag pipelines:latest
	kind load docker-image pipelines:latest

install:
	python -m venv .venv
