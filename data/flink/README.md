# Flink on Kubernetes

The `Dockerfile` is the main docker image that we use to run any pipeline we want. 

The `pipelines` folder contains all pipelines and is copied into the `Dockerfile` when it is build. The makefile contains handy commands for development. 

The `setup.sh` script helps with installing the Flink kubernetes operator.

## Setup
Start a cluster using kind:
``` bash
kind create cluster
```

Install the flink operator
``` bash
./setup.sh
```

Build and load the docker image containing the pipelines
``` bash
make build-local
```

Apply your pipeline
``` bash
kubectl apply -f pipelines/python/demo/deploy.yaml
```

To stop the pipeline from running
``` bash
kubectl delete -f pipelines/python/demo/deploy.yaml
```

To stop and wipe the local cluster
``` bash
kind delete cluster
```
