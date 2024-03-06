#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER docker WITH PASSWORD 'secret';
    CREATE DATABASE komrifa_dev;
    GRANT ALL PRIVILEGES ON DATABASE komrifa_dev TO docker;
    CREATE DATABASE komrifa_test;
    GRANT ALL PRIVILEGES ON DATABASE komrifa_test TO docker;
EOSQL
