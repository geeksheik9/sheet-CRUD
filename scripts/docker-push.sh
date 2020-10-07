#!/usr/bin/env bash

. ./scripts/version.sh

name="sheet-crud"
version="$sheet_crud_version"
repo="repository"
namespace="docker"
username="geeksheik9"
versionedImageName="$name:$version"
taggedImage="$username/$versionedImageName"

echo "version $version"

echo "Building ${taggedImage}"
docker build --no-cache -t ${taggedImage} --build-arg VERSION=${version} .

echo "Please login to docker hub"
docker logout
docker login --username=$username

echo "pushing image ${taggedImage}"
docker push ${taggedImage}