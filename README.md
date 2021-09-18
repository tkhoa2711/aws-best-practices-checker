# aws-best-practices-checker

[![CI](https://github.com/tkhoa2711/aws-best-practices-checker/actions/workflows/ci.yml/badge.svg)](https://github.com/tkhoa2711/aws-best-practices-checker/actions/workflows/ci.yml)

A utility tool to check a given AWS infrastructure against common best practices.

## Requirements

* Go 1.17+
* AWS credentials are properly setup

## Installation

```sh
go install github.com/tkhoa2711/aws-best-practices-checker
```

## Usage

If installed with `go install`, the tool can be invoked with the following command:

```sh
aws-best-practices-checker
```

Or it can also be invoked by cloning this repository and running manually with
`go run`:

```sh
git clone github.com/tkhoa2711/aws-best-practices-checker
cd aws-best-practices-checker
go run main.go
```

## License

MIT License
