#!/usr/bin/env bash
set -euo pipefail

# If no args are provided, run a sane default in GPU mode.
if [ "$#" -eq 0 ]; then
  exec /usr/local/bin/eth-vanity -compatible -mode 2 -engine gpu
fi

exec /usr/local/bin/eth-vanity "$@"
