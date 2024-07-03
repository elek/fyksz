#!/usr/bin/env bash
yq '.' resources/*.yaml | \
fyksz apply ./log-level.sh | \
fyksz save --name "yq '.metadata.name + \"-\" + (.kind | downcase) + \".yaml\"'"
