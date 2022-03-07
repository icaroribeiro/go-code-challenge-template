#!/bin/bash

# Delete the temporary directory used as the context for building the Docker container.
DESTINATION="deployments/heroku/app/.temp"

rm -rf "$DESTINATION"