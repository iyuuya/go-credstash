# go-credstash [![Go](https://github.com/iyuuya/go-credstash/actions/workflows/go.yml/badge.svg)](https://github.com/iyuuya/go-credstash/actions/workflows/go.yml)

## Install

```
$ go install github.com/iyuuya/go-credstash/cmd/credstash@latest
```

## Usage

```
$ credstash
credstash cli

Usage:
  credstash [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      delete a key
  get         Show a value for key name
  help        Help about any command
  list        Show all stored keys
  put         Put a value for key name
  setup       Setup credstash repository on DynamoDB

Flags:
  -e, --endpoint string   Endpoint for dynamodb-local
  -h, --help              help for credstash

Use "credstash [command] --help" for more information about a command.
```

### LocalStack Example
```
$ aws configure --profile localstack
AWS Access Key ID [None]: dummy
AWS Secret Access Key [None]: dummy
Default region name [None]: ap-northeast-1
Default output format [None]: json

$ export AWS_REGION=ap-northeast-1
$ export AWS_ENDPOINT=https://localhost.localstack.cloud:4566

# create ksm key and alias
$ KEY_ID=$(aws --profile localstack --endpoint-url $AWS_ENDPOINT kms create-key | jq -r .KeyMetadata.KeyId); echo $KEY_ID
$ aws --profile localstack --endpoint-url $AWS_ENDPOINT kms create-alias --alias-name alias/credstash --target-key-id $KEY_ID

# setup (create table)
$ credstash setup

# save to dynamodb
$ credstash put hello
secret value> world v1

# versioning
$ credstash put hello
secret value> world v2

# show list
$ credstash list
hello
hello

# show list with version
$ credstash list -v
hello --version: 0000000000000000001
hello --version: 0000000000000000002

# get a value
$ credstash get hello
world v2

# get a value by specific version
$ credstash get hello -v 0000000000000000001
world v1

# deletes
$ credstash delete hello -v 0000000000000000001
$ credstash list -v
hello --version: 0000000000000000002

$ credstash delete hello
$ credstash list
```
