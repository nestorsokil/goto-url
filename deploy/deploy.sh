#!/usr/bin/env bash

set -e

cd ../

version="$(git rev-parse --short=7 HEAD)"

docker build -t nsokil/gotourl:"$version" .
docker push nsokil/gotourl:"$version"

helm upgrade apps ./helm/gotourl --set gotourl.image.tag="$version"

# todo curl grafana annotations

