#!/bin/bash

# Leverage the default env variables as described in:
# https://docs.github.com/en/actions/reference/environment-variables#default-environment-variables
if [[ $GITHUB_ACTIONS != "true" ]]
then
  /usr/bin/crusher "$@"
  exit $?
fi

flags=""

echo "running command:"
echo crusher label -f "$INPUT_FILE" "$flags"

/usr/bin/crusher label -f "$INPUT_FILE" "$flags"
export crusher_EXIT_CODE=$?
