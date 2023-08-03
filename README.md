# Simple DNS Server

## Install
```shell
$ make install
```

## Usage
```shell
$ export CONFIG_JSON='{"records":{"A":[{"name":"from-env.dns","ip":"1.2.3.4","ttl":1234}]}}'
$ export CONFIG_JSON_FILE=./config.example.json
$ make run
$ dig -p 5353 A from-env.com # from-env.dns. 1234 IN A 1.2.3.4
$ dig -p 5353 A example.com # example.com. 1800 IN A 0.0.0.0
$ dig -p 5353 A a.example.com # a.example.com. 3600 IN A 10.10.10.10
$ dig -p 5353 A b.example.com # b.example.com. 3600 IN A 20.20.20.20
$ dig -p 5353 SEO example.com # example.com. 7200 IN SOA ns1.example.com. email.address.com. 1691079983 60 60 60 60
```

## Available env vars
- `CONFIG_JSON` (conf json contents)
- `CONFIG_JSON_FILE` (config json file)
- `ADDR` (Default: 0.0.0.0:5353)
- `LOG_FORMAT` (available formats: json)