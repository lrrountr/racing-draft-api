#!/bin/bash
set -u

CONTAINER=$1
i=0
while true; do
    i=$(($i+1))
    if [ $i -gt 20 ]; then
        echo "Postgres not ready in time..."
        exit 1
    fi
    out=$(docker exec ${CONTAINER} psql -t -d postgres -U postgres -h localhost -p 5432 -c "SELECT count(datname) FROM pg_database where datname = 'postgres'")
    ## Trim our whitespace
    out=$(echo ${out} | xargs)
    if [ "${out}" = "1" ]; then
        exit 0
    fi
    sleep 0.5
done
