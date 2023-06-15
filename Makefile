postgres:
	docker run --name postgresdb --env POSTGRES_PASSWORD="1234512345" --publish "5436:5432" --detach --rm postgres