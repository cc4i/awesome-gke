#!/bin/bash

# Remark: orginally fork from https://raw.githubusercontent.com/kubernetes-retired/contrib/master/startup-script/manage-startup-script.sh
# and modify as per request.

set -o errexit
set -o nounset
set -o pipefail

CHECKPOINT_PATH="${CHECKPOINT_PATH:-/tmp/startup-script.kubernetes.io_$(md5sum <<<"${STARTUP_SCRIPT}" | cut -c-32)}"
CHECK_INTERVAL_SECONDS="30"
EXEC=(nsenter -t 1 -m -u -i -n -p --)

do_startup_script() {
  local err=0;

  echo "${EXEC[@]}" bash -c "${STARTUP_SCRIPT}"
  "${EXEC[@]}" bash -c "${STARTUP_SCRIPT}" && err=0 || err=$?
  if [[ ${err} != 0 ]]; then
    echo "!!! startup-script failed! exit code '${err}'" 1>&2
    return 1
  fi

  "${EXEC[@]}" touch "${CHECKPOINT_PATH}"
  echo "!!! startup-script succeeded!" 1>&2
  return 0
}

while :; do
  "${EXEC[@]}" stat "${CHECKPOINT_PATH}" > /dev/null 2>&1 && err=0 || err=$?
  if [[ ${err} != 0 ]]; then
    do_startup_script
  fi

  sleep "${CHECK_INTERVAL_SECONDS}"
done