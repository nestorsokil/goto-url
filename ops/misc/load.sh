#!/usr/bin/env bash

set -e

function print_usage() {
  echo "usage: $0 [-ulhs]"
  echo "  -u   --url          URL to load"
  echo "  -l   --sleep-low    Time bound for sleep (lower)"
  echo "  -h   --sleep-high   Time bound for sleep (higher)"
  echo "  -s   --silent       Don't print output"
}

for i in "$@"; do
  case $i in
  -u=* | --url=*)
    URL="${i#*=}"
    shift
    ;;
  -l=* | --sleep-low=*)
    SLEEP_LO="${i#*=}"
    shift
    ;;
  -h=* | --sleep-high=*)
    SLEEP_HI="${i#*=}"
    shift
    ;;
  -s=* | --silent=*)
    SILENT="true"
    shift
    ;;
  --help=*)
    print_usage
    exit 1
    shift
    ;;
  *) ;;

  esac
done

if [[ -z "$SILENT" ]]; then
    echo "Putting some load on $URL..."
fi

while true
do
  curl -s "$URL" > /dev/null
  SLEEP_TIME=$(echo "scale=2; ($[($RANDOM%$SLEEP_HI)+$SLEEP_LO]) / 1000.0" | bc -l)
  echo $SILENT
  if [[ -z "$SILENT" ]]; then
    echo Sleeping for $SLEEP_TIME
  fi
  sleep $SLEEP_TIME
done