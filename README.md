# PowerDNS Go client

See PowerDNS API spec at https://doc.powerdns.com/md/httpapi/api_spec

This is a partial implementation of the API.

Currently only the following is supported
```
CreateZone via a POST to /api/v1/servers/localhost/zones

CreateRecord via POST to /api/v1/servers/localhost/zones/:zone_id
```


## Building

```shell
go build
```