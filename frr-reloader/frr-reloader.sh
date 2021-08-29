#!/bin/bash
cleanup() {
  echo "Caught an exit signal.."
  rm -f "$PIDFILE"
  rm -f "$LOCKFILE"
  kill_sleep
  exit
}

reload_frr() {
  flock 200
  echo "Caught SIGHUP and acquired lock! Reloading FRR.."
  kill_sleep

  echo "Checking the configuration file syntax"
  if ! python3 /usr/lib/frr/frr-reload.py --test --stdout "$FILE_TO_RELOAD" ; then
    echo "Syntax error spotted: aborting.."
    return
  fi

  echo "Applying the configuration file"
  if ! python3 /usr/lib/frr/frr-reload.py --reload --overwrite --stdout "$FILE_TO_RELOAD" ; then
    echo "Failed to fully apply configuration file"
    return
  fi
  
  echo "FRR reloaded successfully!"
} 200<"$LOCKFILE"

kill_sleep() {
  kill "$sleep_pid"
}

trap cleanup SIGTERM SIGQUIT SIGINT
trap 'reload_frr &' HUP

SHARED_VOLUME="${SHARED_VOLUME:-/etc/frr_reloader}"
PIDFILE="$SHARED_VOLUME/reloader.pid"
FILE_TO_RELOAD="$SHARED_VOLUME/frr.conf"
LOCKFILE="$SHARED_VOLUME/lock"

echo "PID is: $$, writing to $PIDFILE"
printf "$$" > "$PIDFILE"
touch "$LOCKFILE"

while true
do
    sleep infinity &
    sleep_pid=$!
    wait $sleep_pid 2>/dev/null
done
