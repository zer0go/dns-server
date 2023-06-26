# Simple DNS Server

## Install
```shell
$ make install
```

## Usage
```shell
$ export RECORDS_JSON='{"example.local":"1.2.3.4"}'
$ make run
$ dig -p 15353 example.local
$ BASE_DOMAIN=dns.com make run
$ dig -p 15353 example.local.dns.com
```

## Available env vars
- `RECORDS_JSON` (records json contents)
- `RECORDS_JSON_FILE` (records json file)
- `BASE_DOMAIN`
- `PORT` (Default: 15353)