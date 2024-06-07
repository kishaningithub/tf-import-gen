# tf-import-gen (Terraform import generator)

[![Build Status](https://github.com/kishaningithub/tf-import-gen/actions/workflows/build.yml/badge.svg)](https://github.com/kishaningithub/tf-import-gen/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kishaningithub/tf-import-gen)](https://goreportcard.com/report/github.com/kishaningithub/tf-import-gen)
[![Latest release](https://img.shields.io/github/release/kishaningithub/tf-import-gen.svg)](https://github.com/kishaningithub/tf-import-gen/releases)

Tool to generate terraform import statements to simplify state migrations from one terraform code base to another.

<!-- TOC -->
* [tf-import-gen (Terraform import generator)](#tf-import-gen-terraform-import-generator)
  * [Installation](#installation)
    * [Using Homebrew (Mac and linux)](#using-homebrew-mac-and-linux)
    * [Using docker](#using-docker)
    * [Using pkgx](#using-pkgx)
    * [Others](#others)
  * [Examples](#examples)
    * [Generating import statements by module](#generating-import-statements-by-module)
    * [Generating import statements by resource](#generating-import-statements-by-resource)
    * [Generating import statements by multiple resource](#generating-import-statements-by-multiple-resource)
    * [Generating import statements for all resources](#generating-import-statements-for-all-resources)
  * [Usage](#usage)
  * [Contributing](#contributing)
<!-- TOC -->

## Installation

### Using Homebrew (Mac and linux)

```bash
brew install kishaningithub/tap/tf-import-gen
```

### Using docker

```bash
alias tf-import-gen="docker run -i ghcr.io/kishaningithub/tf-import-gen:latest"

terraform show -json | tf-import-gen
```

### Using pkgx

```bash
terraform show -json | pkgx tf-import-gen@latest
```

or if you prefer the env style which only adds tf-import-gen to your current shell session

```bash
env +tf-import-gen@latest

terraform show -json | tf-import-gen --version
```

### Others

Head over to the [releases page](https://github.com/kishaningithub/tf-import-gen/releases) and download a binary for your platform

## Examples

### Generating import statements by module

```bash
$ terraform show -json | tf-import-gen module.example

import {
  to = module.example.aws_glue_catalog_database.example_db
  id = "123456789012:example_db"
}

import {
  to = module.example.aws_iam_instance_profile.example_instance_profile
  id = "example_instance_profile"
}
```

### Generating import statements by resource

```bash
$ terraform show -json | tf-import-gen aws_instance.example

import {
  to = aws_instance.example
  id = "i-123456789012"
}
```

### Generating import statements by multiple resource

```bash
$ terraform show -json | tf-import-gen aws_instance.example module.example

import {
  to = aws_instance.example
  id = "i-123456789012"
}

import {
  to = module.example.aws_glue_catalog_database.example_db
  id = "123456789012:example_db"
}

import {
  to = module.example.aws_iam_instance_profile.example_instance_profile
  id = "example_instance_profile"
}
```

### Generating import statements for all resources

```bash
$ terraform show -json | tf-import-gen

import {
  to = module.example.aws_glue_catalog_database.example_db
  id = "123456789012:example_db"
}

import {
  to = module.example.aws_iam_instance_profile.example_instance_profile
  id = "example_instance_profile"
}

import {
  to = aws_instance.example
  id = "i-123456789012"
}
```

## Usage

```bash
$ tf-import-gen --help

Generate terraform import statements to simplify state migrations from one terraform code base to another.

The address argument can be used to filter the instances by resource or module. If
no pattern is given, import statements are generated for all the resources.

The addresses must either be module addresses or absolute resource
addresses, such as:
  aws_instance.example
  module.example
  module.example.module.child
  module.example.aws_instance.example

Usage:
  tf-import-gen [flags] address...

Examples:

## Generating import statements by module
terraform show -json | tf-import-gen module.example

## Generating import statements by resource
terraform show -json | tf-import-gen aws_instance.example

## Generating import statements by multiple resources
terraform show -json | tf-import-gen aws_instance.example module.example

## Generating import statements for all resources
terraform show -json | tf-import-gen


Flags:
  -h, --help      help for tf-import-gen
  -v, --version   version for tf-import-gen
```


## Contributing

PRs are always welcome!. Refer [CONTRIBUTING.md](./CONTRIBUTING.md) for more information



