# Create a new migration
db/migrations/new:
	migrate create -seq -ext=.sql -dir=./migrations $(name)

# Run migrations up
db/migrations/up:
	migrate -path=./migrations -database=$(DB_DSN) up

# Run migrations down
db/migrations/down:
	migrate -path=./migrations -database=$(DB_DSN) down

# Open postgres
db/psql:
	psql $(DB_DSN)