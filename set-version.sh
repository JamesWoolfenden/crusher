#!/bin/sh
set -e
latestTag=$(git describe --tags)
{
  echo "Updating version file with new tag: $latestTag"
  echo "package version"
  echo ""
  echo "// Version controls the applications version."
  echo "const Version = \"$latestTag\""
} >> src/version/version.go
