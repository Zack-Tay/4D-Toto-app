DB_NAME=sg_lottery
DB_USER=postgres

setup:
	psql -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);"
	psql -U $(DB_USER) -d $(DB_NAME) -f migrations/001_init.up.sql
	psql -U $(DB_USER) -d $(DB_NAME) -f scripts/setup.sql

# Reset database
reset:
	psql -U $(DB_USER) -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	make setup

# Connect to database
connect:
	psql -U $(DB_USER) -d $(DB_NAME)

.PHONY: setup reset connect