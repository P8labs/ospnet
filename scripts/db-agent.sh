#!/bin/bash

DB=agent.db

goose -dir internal/agent/db/migrations sqlite3 $DB up
cd internal/agent/db && sqlc generate