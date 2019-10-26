#!/usr/bin/env bash

# rollback.sh --release=apps --revision=0 --grafana-url=192.168.99.100:30540 --grafana-auth=eyJrIjoiQWxaRXJJN1NUMHFzUzZaMHlNQUtwck0xVWdmN3hFblEiLCJuIjoiRGVwbG95bWVudCBBbm5vdGF0aW9ucyIsImlkIjoxfQ==

function print_usage() {
  echo "usage: $0 [-gar]"
  echo "  -g   --grafana-url     Grafana URL"
  echo "  -a   --grafana-auth    Grafana auth token"
  echo "  -r   --release         Helm release"
  echo "  -v   --revision        Helm release revision to rollback to"
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
  -v=* | --revision=*)
    REVISION="${i#*=}"
    shift
    ;;
  *) ;;

  esac
done

if [ -z "$GRAFANA_URL" ] || [ -z "$GRAFANA_AUTH" ] || [ -z "$RELEASE" ]; then
  print_usage
  exit 1
fi

if [ -z "$REVISION" ]; then
  REVISION=0 # last successful revision
fi

# todo nsokil kinda sloppy, there may be a few "tag"s
ROLLED_VERSION="$(helm get values "$RELEASE" | grep 'tag' | awk '{print $2}')"

helm rollback "$RELEASE" "$REVISION"
curl -XPOST "$GRAFANA_URL"/api/annotations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $GRAFANA_AUTH" \
  --data @- <<EOF
  {
    "text": "Rollback deployment 'gotourl:$ROLLED_VERSION'\n\n
      <a href=\"https://github.com/nestorsokil/goto-url/commit/$ROLLED_VERSION\">GitHub ($ROLLED_VERSION)</a>",
    "tags": [
      "rollback"
    ]
  }
EOF
