docker-build:
	docker build -f ./deployments/build/Dockerfile -t urlshortener .

docker-deploy:
	docker build -f ./deployments/build/Dockerfile -t urlshortener . && \
	docker-compose -p urlshortener -f ./deployments/environment/docker-compose.yaml up -d

docker-deploy-down:
	docker-compose -p urlshortener -f ./deployments/environment/docker-compose.yaml down

env:
	docker-compose -p url_shortener -f ./deployments/environment/docker-compose.dev.yaml up -d

env-down:
	docker-compose -p url_shortener -f ./deployments/environment/docker-compose.dev.yaml down

mocks:
	mockery --all --with-expecter --dir ./pkg/app/urlshortener --output ./pkg/app/urlshortener/mocks


goose:
	goose -dir ./deployments/migrate/ -v postgres "user=url_shortener_service dbname=url_shortener host=127.0.0.1 port=5432 password=pqV7EJ8bYJpFDXXJtw66s6JKG4xpZb4v sslmode=disable" up
