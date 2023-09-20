#!/bin/bash

# Define the database file and table name
DB_FILE="db.sqlite"
TABLE_NAME="value_table"

# Function to insert random data into the database
insert_random_data() {
    local current_timestamp=$(date +%s)
    local two_weeks_ago=$((current_timestamp - 1209600)) # 2 weeks in seconds (60s * 60s * 24s * 14 days)
    local timestamp=$((two_weeks_ago + (RANDOM * 1337) % (current_timestamp - two_weeks_ago)))


    local value=$(awk -v min=1 -v max=10 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')

    # Use sqlite-utils to insert data into the specified table
    echo "{\"timestamp\": $timestamp, \"value\": $value}" | sqlite-utils insert "$DB_FILE" "$TABLE_NAME" -
}

# Insert 10 random data points (adjust the number as needed)
for ((i=0; i<10; i++)); do
    insert_random_data
done
