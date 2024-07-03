yq "select(.spec.template.spec | has (\"containers\")) | .spec.template.spec.containers[0].args += \"--log.custom-level=piecestore=WARN,collector=WARN,reputation:service=WARN\""
