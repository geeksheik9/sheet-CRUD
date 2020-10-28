#!/usr/bin/env bash
echo "Generating swagger"
swagger generate spec -o ./swagger-ui/swagger.yaml -m
#to view swagger docs locally run
#swagger serve ./swagger-ui/swagger.yaml