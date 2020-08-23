.PHONY: deploy
deploy: compile-js docker-build docker-deploy

.PHONY: docker-build
docker-build:
	docker build -f backend/Dockerfile -t docker.io/chuckdha/kubeyaml-backend:latest backend
	docker build -f frontend/Dockerfile -t docker.io/chuckdha/kubeyaml-frontend:latest frontend

.PHONY: docker-deploy
docker-deploy:
	docker push docker.io/chuckdha/kubeyaml-backend:latest
	docker push docker.io/chuckdha/kubeyaml-frontend:latest

.PHONY: compile-js
compile-js:
	npm run --prefix frontend build

