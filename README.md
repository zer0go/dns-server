# Simple DNS Server

## Install
```shell
$ make install
```

## Usage
```shell
$ export RECORDS_JSON='{"example.local":"1.2.3.4"}'
$ make run
$ dig -p 5353 example.local # example.local.          3600 IN A 1.2.3.4
$ BASE_DOMAIN=dns.com make run
$ dig -p 5353 example.local.dns.com # example.local.dns.com. 3600 IN A 1.2.3.4
```

## Available env vars
- `RECORDS_JSON` (records json contents)
- `RECORDS_JSON_FILE` (records json file)
- `RELOAD_CONFIG` (reload records before every request, default: false)
- `BASE_DOMAIN`
- `ADDR` (Default: 0.0.0.0:15353)