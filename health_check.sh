#!/bin/sh

DB_HOST=$POSTGRES_HOST
DB_USER=$POSTGRES_USERNAME
DB_PASSWORD=$POSTGRES_PASSWORD
DB_NAME=$POSTGRES_NAME

# Attempt to connect to the database using psql
if psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1" >/dev/null; then
    echo "Database connection successful"
    exit 0
else
    echo "Failed to connect to database"
    exit 1
fi
