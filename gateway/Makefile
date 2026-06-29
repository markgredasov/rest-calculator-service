include .env
export

export PROJECT_ROOT=$(shell pwd)
export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs

up:
	@go mod tidy && \
	go run cmd/app/main.go