package template

import (
	"fmt"
	"regexp"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/arvancloud"
)

func isValidAPIKey(key string) bool {
	// Standard UUID regex: 8-4-4-4-12 hex characters
	re := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	return re.MatchString(key)
}

// Provider lets Caddy read and manipulate DNS records hosted by this DNS provider.
type Provider struct{ *arvancloud.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.arvancloud",
		New: func() caddy.Module { return &Provider{new(arvancloud.Provider)} },
	}
}

// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	p.Provider.AuthAPIKey = caddy.NewReplacer().ReplaceAll(p.Provider.AuthAPIKey, "")
	if !isValidAPIKey(p.Provider.AuthAPIKey) {
		return fmt.Errorf("ApiKey '%s' appears invalid; ensure it's correctly entered and not wrapped in braces nor quotes", p.Provider.AuthAPIKey)
	}
	return nil
}


// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	Single API Token
//
//	arvancloud <api_key>      // Domain read access and domain DNS write for all domains
//
//	Single API Token, alternative syntax
//
//	arvancloud {
//	  api_key <api_key>     // Domain read access and dmain DNS write for all domains
//	}
//
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.Provider.AuthAPIKey = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_key":
				if p.Provider.AuthAPIKey != "" {
					return d.Err("API token already set")
				}
				if d.NextArg() {
					p.Provider.AuthAPIKey = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.AuthAPIKey == "" {
		return d.Err("missing API key")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
