#!/usr/bin/env bash
yq '.' resources/*.yaml | \
fyksz apply DEPLOYMENT_FILTER="yq '. | select(.metadata.annotations.component == \"storagenode\")'" ./log-level.sh | \
fyksz save --name "yq '.metadata.name + \"-\" + (.kind | downcase) + \".yaml\"'"
