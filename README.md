# doc-integration

## Purpose
Repository for defining DMS (Document Management System) interoperation API.

## Contents
* Registry - project written in Go to act like a centralized registry for storing all of the intercommunicating organization info, including public keys and URLs.
* openapi.yml (root dir) - proposal for a unified API for receiving documents among different DMS vendors.

## Registry deployment
### Requirements
1. PostgreSQL >12
2. Golang (only for compilation)

### Steps
1. Create user, DB in Postgres.
2. Change DB related values in config.json
3. Populate DB as necessary with organizations.
4. Run respective build script for your OS.
5. Execute the binary.
