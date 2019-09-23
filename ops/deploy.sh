#!/usr/bin/env bash

# deploy.sh --release=apps --grafana-url=192.168.99.100:30540 --grafana-auth=eyJrIjoiQWxaRXJJN1NUMHFzUzZaMHlNQUtwck0xVWdmN3hFblEiLCJuIjoiRGVwbG95bWVudCBBbm5vdGF0aW9ucyIsImlkIjoxfQ==

function print_usage() {
  echo "usage: $0 [-gar] [-f infile] [-o outfile]"
  echo "  -g   --grafana-url     Grafana URL"
  echo "  -a   --grafana-auth    Grafana auth token"
  echo "  -r   --release         Helm release"
}

set -e

for i in "$@"; do
  case $i in
  -g=* | --grafana-url=*)
    GRAFANA_URL="${i#*=}"
    shift
    ;;
  -a=* | --grafana-auth=*)
    GRAFANA_AUTH="${i#*=}"
    shift
    ;;
  -r=* | --release=*)
    RELEASE="${i#*=}"
    shift
    ;;
  -h=* | --help=*)
    print_usage
    exit 1
    shift
    ;;
  *) ;;

  esac
done

if [ -z "$GRAFANA_URL" ] || [ -z "$GRAFANA_AUTH" ] || [ -z "$RELEASE" ]; then
  print_usage
  exit 1
fi

cd ../
VERSION="$(git rev-parse --short=7 HEAD)"

docker build -t nsokil/gotourl:"$VERSION" .
docker push nsokil/gotourl:"$VERSION"

helm upgrade "$RELEASE" ./ops/helm/gotourl --set image.tag="$VERSION"

curl -XPOST "$GRAFANA_URL"/api/annotations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $GRAFANA_AUTH" \
  --data @- <<EOF
  {
    "text": "Deployment 'gotourl:$VERSION'\n\n
      <a href=\"https://github.com/nestorsokil/goto-url/commit/$VERSION\">GitHub ($VERSION)</a>",
    "tags": [
      "deployment"
    ]
  }
EOF
