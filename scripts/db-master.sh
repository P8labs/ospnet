#!/bin/bash

DB=master.db

goose -dir internal/master/db/migrations sqlite3 $DB up
cd internal/master/db && sqlc generate