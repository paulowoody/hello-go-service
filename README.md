# Hello Go Service Deployment Guide

This repository contains a small Go HTTP service that can be run locally with Podman or deployed to a local Kubernetes cluster with Kind.

## Architecture

The project has three main parts:

1. **Go HTTP service**  
   A lightweight web server that exposes:
    * `/` for the main response
    * `/health` for liveness checks
    * `/ready` for readiness checks
    * `/fail` to simulate an unhealthy state

2. **Container image**  
   The Go application is packaged into a container image so it can run consistently in different environments.

3. **Kubernetes deployment**  
   The app can be deployed to a local Kind cluster using a Deployment and a NodePort Service, with health probes configured for Kubernetes.

## Repository structure

The repository is organized to separate application code, deployment manifests, and build files:

* `cmd/` — application entry points
* `internal/` — internal application packages
* `deploy/` — Kubernetes manifests and Kind configuration
* `build/` — container build files
* `go.mod` — Go module definition

## Build from the repository root

The container build expects the repository root to be the build context. That way, the Dockerfile in `build/` can copy the files it needs from the project root:

* `podman build -f build/Containerfile -t localhost/hello-go-service .`

## Prerequisites

Before you start, make sure the required tools are installed:

* podman
* kind
* kubectl

If you're using Kind with Podman, you may also need Docker/Podman compatibility settings:

* `sudo apt update && sudo apt install podman-docker`
* `export DOCKER_HOST="unix://$XDG_RUNTIME_DIR/podman/podman.sock"`


## Workflow

### 1. Build the container image

Build the Go application into a container image from the repository root:

* `podman build -f build/Containerfile -t localhost/hello-go-service .`

### 2. Run the service locally with Podman

This is the quickest way to test the service without Kubernetes. It maps port `8082` on your machine to port `8080` inside the container:

* `podman run -d -p 8082:8080 hello-go-service`

### 3. Deploy to Kubernetes with Kind

Use this option if you want to test the app in a Kubernetes environment. Kind runs a local cluster inside containers.

Create the cluster:

* `KIND_EXPERIMENTAL_PROVIDER=podman kind create cluster --name go-learning --config deploy/kind-config.yaml`

Load the image into the cluster so Kubernetes can use it:

* `kind load docker-image localhost/hello-go-service:latest --name go-learning`

Deploy the application and service:

* `kubectl apply -f deploy/deployment.yaml`
* `kubectl apply -f deploy/service.yaml`

If you update the image but keep the same tag, restart the deployment so Kubernetes runs the newer image:

* `kubectl rollout restart deployment hello-go-deployment`
* `kubectl rollout status deployment hello-go-deployment`

The `rollout status` command is useful because it tells you when the new pods are ready and the update has completed successfully.

### 4. Test the endpoints

These checks help confirm the app is responding correctly:

* `curl -v localhost:8082/`  
  Returns the main response for the app.

* `curl -v localhost:8082/health`  
  Checks whether the service is healthy.

* `curl -v localhost:8082/ready`  
  Checks whether the service is ready to receive traffic.

* `curl -v localhost:8082/fail`  
  Forces the service into an unhealthy state for testing.

* `curl -v localhost:8082/health`  
  Should now return `500` after calling `/fail`.

### 5. Update the image after making changes

If you change the Go code, rebuild the image and redeploy it.

#### Option 1: Reuse the same image tag

Use this when you want to keep the same tag, such as `latest`, and refresh the running workload manually:

* Rebuild the image
* Reload it into Kind:
    * `kind load docker-image localhost/hello-go-service:latest --name go-learning`
* Restart the deployment:
    * `kubectl rollout restart deployment hello-go-deployment`
* Wait for the update to finish:
    * `kubectl rollout status deployment hello-go-deployment`

This works well for local development, but it requires a manual restart because Kubernetes does not automatically detect changes when the image tag stays the same.

#### Option 2: Use a new image tag

Use this when you want deployments to clearly track a specific image version:

* Build a versioned image, for example:
    * `podman build -t localhost/hello-go-service:v2 .`
* Update `deploy/deployment.yaml` to use the new tag
* Apply the updated manifest:
    * `kubectl apply -f deploy/deployment.yaml`
* Check rollout progress:
    * `kubectl rollout status deployment hello-go-deployment`

This is usually the cleaner approach because each version is explicit and easier to trace.
