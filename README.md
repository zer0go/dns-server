# Simple DNS Server

## Install
```shell
$ make install
```

## Usage
```shell
$ export RECORDS_FILE=./records.example.json
$ make run
$ dig -p 5353 A example.com # example.com. 1800 IN A 0.0.0.0
$ dig -p 5353 A a.example.com # a.example.com. 3600 IN A 10.10.10.10
$ dig -p 5353 A b.example.com # b.example.com. 3600 IN A 20.20.20.20
$ dig -p 5353 SOA example.com # example.com. 7200 IN SOA ns1.example.com. email.address.com. 2023080700 3600 1800 2592000 3600
$ dig -p 5353 CNAME api.example.com # api.example.com. 1800 IN CNAME example.com.
$ dig -p 5353 NS example.com # ns1.example.com. 3600 IN NS ns1.example.com
```

## Available env vars
- `RECORDS_FILE` (records json file path)
- `RELOAD_RECORDS_INTERVAL` (reload records file at the specified intervals, default: 0 [disabled])
- `ADDR` (default: 0.0.0.0:5353)
- `LOG_FORMAT` (available formats: json)