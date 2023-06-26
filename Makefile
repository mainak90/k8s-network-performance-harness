NAME ?= perftest
GO ?= go
VERSION ?= $(shell cat VERSION)
NS := $(if $(namespace),$(namespace),default)

build-linux: # Builds the perftest cli but only for linux
	GOOS=linux GOARCH=amd64 go build -o bin/perftest main.go

build: # Builds the perftest cli but only for macos
	GOOS=darwin GOARCH=amd64 go build -o bin/perftest main.go

kubernetes: # this should output a valid Kubernetes spec for your web service, this assumes that kubectl binary is installed.
	kubectl apply -f deploy/manifests/perftest.yaml
	kubectl rollout status deployment netperf-pod
	kubectl rollout status daemonset netperf-pod

helm-vet:
	make helm-lint
	make helm-template

helm-lint:
	helm lint ./deploy/helm

helm-template:
	helm template test ./deploy/helm

helm: # Use "make helm namespace=<namespace>" to install on a custom namespace
	helm install perftest ./deploy/helm -f ./deploy/helm/values.yaml --namespace $(NS) --create-namespace --wait
