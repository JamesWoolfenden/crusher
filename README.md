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

crusher manages labels in Dockerfiles and their layers

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
docker run --tty --volume /local/path/to/tf:/tf jameswoolfenden/crusher scan -d /tf
```

<https://hub.docker.com/repository/docker/jameswoolfenden/crusher>

## Usage

### Directory scan

This will look for the .github/workflow folder and update all the files it finds
there, and display a diff of the changes made to each file:

```bash
$crusher label -d .
```

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
