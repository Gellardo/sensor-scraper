#!/bin/bash

set -euo pipefail

# Define the database file and table name
DB_FILE="db.sqlite"
TABLE_NAME="value_table"

require() { hash "$@" || return 127; }
require sqlite-utils

# Function to insert random data into the database
generate_random_entry() {
    local current_timestamp=$(date +%s)
    local two_weeks_ago=$((current_timestamp - 1209600)) # 2 weeks in seconds (60s * 60s * 24s * 14 days)
    local timestamp=$((two_weeks_ago + (RANDOM * 1337) % (current_timestamp - two_weeks_ago)))


    local value=$(awk -v min=1 -v max=10 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    local sensor=$(( RANDOM % 3 ))

    # Use sqlite-utils to insert data into the specified table
    echo "{\"sensorid\": $sensor, \"timestamp\": $timestamp, \"value\": $value}"
}

# Insert 10 random data points (adjust the number as needed)
for ((i=0; i<100; i++)); do
    generate_random_entry
done | sqlite-utils insert --nl "$DB_FILE" "$TABLE_NAME" -
