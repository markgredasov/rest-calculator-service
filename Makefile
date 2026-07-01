export PROJECT_ROOT=$(pwd)

deploy-up:
	@docker compose up --build gateway computing-core

deploy-down:
	@docker compose down gateway computing-core
