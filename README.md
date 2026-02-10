Arvancloud module for Caddy
===========================

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records with Arvancloud accounts.

## Caddy module name

```
dns.providers.arvancloud
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
	"module": "acme",
	"challenges": {
		"dns": {
			"provider": {
				"name": "arvancloud",
				"api_key": "your_Arvancloud_api_key"
			}
		}
	}
}
```

or with the Caddyfile:

```
# globally
{
	acme_dns arvancloud your_Arvancloud_api_key
}
```

```
# one site
tls {
	dns arvancloud your_Arvancloud_api_key
}
```

## Caveats
Arvancloud doesn't currently allow more than 1 IP address for a single A record on the free plan. this might change in the future but for now that means the API will return an error in this case.