#!/usr/bin/env bash
. ./scripts/version.sh

echo "Please enter Mongo url:"
read mongourl

docker run -d -t -i -p 3002:3000 -e LOCAL_MONGO="$mongourl" geeksheik9/sheet-crud:$sheet_crud_version