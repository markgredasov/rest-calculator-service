export PROJECT_ROOT=$(pwd)

deploy-up:
	@docker compose up --build gateway computing-core

deploy-down:
	@docker compose down gateway computing-core

tests:
	@cd ./gateway && echo "testing gateway..." && go test -v ./internal/features/... && cd .. && \
	cd ./computing-core && echo "testing computing-core..." && go test -v ./internal/features/... && cd ..
