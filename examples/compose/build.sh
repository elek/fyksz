#!/usr/bin/env bash
set -euo pipefail
yq '.' templates/base.yaml templates/redis.yaml templates/cockroach.yaml | \
   fyksz apply --filter "yq '.services.*.labels | select(.app=\"redis\")'"  yq '.services.*.image="xxx"' | \
   yq eval-all '. as $item ireduce ({}; . * $item )' > docker-compose.yaml


#
@compose filter redis | yq '..services.*.image.=xxx'
yq compose output
