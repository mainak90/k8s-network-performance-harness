NAME ?= perftest
GO ?= go
VERSION ?= $(shell cat VERSION)
NS := $(if $(namespace),$(namespace),default)

build-linux: # Builds the perftest cli but only for linux
	GOOS=linux GOARCH=amd64 go build -o bin/perftest main.go

build: # Builds the perftest cli but only for macos
	GOOS=darwin GOARCH=arm64 go build -o bin/perftest main.go

kubernetes: # Assumes that kubectl binary is installed. Please pass in make kubernetes namespace=<namespace to use custom ns>
	kubectl apply -f deploy/manifests/perftest.yaml -n $(NS)
	kubectl rollout status deployment netperf-pod -n $(NS)
	kubectl rollout status daemonset netperf-pod -n $(NS)

helm-vet:
	make helm-lint
	make helm-template

helm-lint:
	helm lint ./deploy/helm

helm-template:
	helm template test ./deploy/helm

helm: # Use "make helm namespace=<namespace>" to install on a custom namespace
	helm install perftest ./deploy/helm -f ./deploy/helm/values.yaml --namespace $(NS) --create-namespace --wait
