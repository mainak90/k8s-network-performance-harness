NAME ?= perftest
GO ?= go
VERSION ?= $(shell cat VERSION)

build-linux: # Builds the perftest cli but only for linux
	GOOS=linux GOARCH=amd64 go build -o bin/perftest main.go

build: # Builds the perftest cli but only for macos
	GOOS=darwin GOARCH=amd64 go build -o bin/perftest main.go

kubernetes: # this should output a valid Kubernetes spec for your web service, this assumes that kubectl binary is installed.
	kubectl apply -f deploy/manifests/perftest.yaml
	kubectl rollout status deployment netperf-pod
	kubectl rollout status daemonset netperf-pod