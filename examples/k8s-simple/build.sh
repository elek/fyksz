#!/usr/bin/env bash
yq '.' resources/*.yaml | \
fyksz save --name "yq '.metadata.name + \"-\" + (.kind | downcase) + \".yaml\"'"
