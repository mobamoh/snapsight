dock_up:
    docker compose up
docker_down:
    docker compose down
modd:
    modd
goose_status:
    goose -dir "migrations" postgres "host=localhost port=5432 user=admin password=nimda dbname=snapsight sslmode=disable" status
goose_up:
    goose -dir "migrations" postgres "host=localhost port=5432 user=admin password=nimda dbname=snapsight sslmode=disable" up
goose_down:
    goose -dir "migrations" postgres "host=localhost port=5432 user=admin password=nimda dbname=snapsight sslmode=disable" down


# The --remove-orphans flag removes containers we may not have defined in our
# docker-compose, but are still present from old configs.
docker compose down --remove-orphans
docker compose up -d