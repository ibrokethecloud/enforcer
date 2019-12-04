#!/bin/sh
set -e

echo "Going to update Trivy cache before booting the webhook!"
trivy  --download-db-only

exec "$@"