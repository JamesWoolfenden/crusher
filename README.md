# crusher

![alt text](Crusher.jpg "crusher")

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/jameswoolfenden/crusher/graphs/commit-activity)
[![Build Status](https://github.com/JamesWoolfenden/crusher/workflows/CI/badge.svg?branch=main)](https://github.com/JamesWoolfenden/crusher)
[![Latest Release](https://img.shields.io/github/release/JamesWoolfenden/crusher.svg)](https://github.com/JamesWoolfenden/crusher/releases/latest)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/JamesWoolfenden/crusher.svg?label=latest)](https://github.com/JamesWoolfenden/crusher/releases/latest)
![Terraform Version](https://img.shields.io/badge/tf-%3E%3D0.14.0-blue.svg)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
[![checkov](https://img.shields.io/badge/checkov-verified-brightgreen)](https://www.checkov.io/)
[![Github All Releases](https://img.shields.io/github/downloads/jameswoolfenden/crusher/total.svg)](https://github.com/JamesWoolfenden/crusher/releases)
[![codecov](https://codecov.io/gh/JamesWoolfenden/crusher/graph/badge.svg?token=PT5FSJF2U3)](https://codecov.io/gh/JamesWoolfenden/crusher)

Crusher is a Google Cloud BigTable maintenance utility that compacts tables by deleting old rows based on timestamp filters and key patterns.

## Table of Contents

<!--toc:start-->
- [crusher](#crusher)
  - [Table of Contents](#table-of-contents)
  - [Install](#install)
    - [MacOS](#macos)
    - [Windows](#windows)
    - [Docker](#docker)
  - [Usage](#usage)

<!--toc:end-->

## Install

Download the latest binary here:

<https://github.com/JamesWoolfenden/crusher/releases>

Install from code:

- Clone repo
- Run `go install`

Install remotely:

```shell
go install  github.com/jameswoolfenden/crusher@latest
```

### MacOS

```shell
brew tap jameswoolfenden/homebrew-tap
brew install jameswoolfenden/tap/crusher
```

### Windows

I'm now using Scoop to distribute releases,
it's much quicker to update and easier to manage than previous methods,
you can install scoop from <https://scoop.sh/>.

Add my scoop bucket:

```shell
scoop bucket add iac https://github.com/JamesWoolfenden/scoop.git
```

Then you can install a tool:

```bash
scoop install crusher
```

### Docker

```shell
docker pull jameswoolfenden/crusher
docker run --tty jameswoolfenden/crusher clip --project my-project --instance my-instance --table my-table --dry-run
```

**Note:** You'll need to mount Google Cloud credentials for authentication:

```shell
docker run --tty \
  -v ~/.config/gcloud:/root/.config/gcloud \
  jameswoolfenden/crusher clip --project my-project --instance my-instance --table my-table
```

<https://hub.docker.com/repository/docker/jameswoolfenden/crusher>

## Usage

Crusher compacts Google Cloud BigTable instances by deleting rows older than a specified number of days that match a key pattern.

### Basic Usage

Delete rows from a BigTable instance (dry-run mode by default):

```bash
crusher clip --project my-project --instance my-instance --table my-table --keyfilter ".*" --days 180 --dry-run
```

### Parameters

**Required:**
- `--project, -p`: GCP project ID
- `--instance, -i`: BigTable instance ID
- `--table, -t`: BigTable table name

**Optional:**
- `--keyfilter, -k`: Regex pattern to match row keys (default: ".*" - matches all rows)
- `--days, -d`: Delete rows older than this many days (default: 180)
- `--dry-run, -r`: Preview deletions without executing them (default: false)
- `--yes, -y`: Skip confirmation prompt (default: false)

### Examples

**Preview deletions** (recommended first step):

```bash
crusher clip --project my-project --instance my-instance --table my-table --dry-run
```

**Delete old rows with confirmation prompt**:

```bash
crusher clip --project my-project --instance my-instance --table my-table --days 90
```

This will:
1. Preview how many rows would be deleted
2. Ask for confirmation before proceeding
3. Delete the rows if confirmed

**Delete old chat history rows**:

```bash
crusher clip --project my-project --instance my-instance --table chat-data --keyfilter ".*chat_histories$" --days 90
```

**Delete all rows older than 1 year (skip confirmation)**:

```bash
crusher clip --project my-project --instance my-instance --table logs --keyfilter ".*" --days 365 --yes
```

**Note:** By default, crusher will always preview deletions and ask for confirmation unless you use `--dry-run` or `--yes` flags.

## Help

```bash
                    _
 __  _ _  _  _  ___| |_   ___  _ _
/ _|| '_|| || |(_-<| ' \ / -_)| '_|
\__||_|   \_,_|/__/|_||_|\___||_|
version: 9.9.9
NAME:
   crusher - AISB utility

USAGE:
   crusher [global options] command [command options]

VERSION:
   9.9.9

AUTHOR:
   James Woolfenden <jim.wolf@duck.com>

COMMANDS:
   clip, c     Compacts BigTable
   version, v  Outputs the application version
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

```

## Building

```go
go build
```

or

```Make
Make build
```

## Extending

Log an issue, a pr or email jim.wolf @ duck.com.
