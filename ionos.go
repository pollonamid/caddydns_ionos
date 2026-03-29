package ionos

import (
	"context"

	caddy "github.com/caddyserver/caddy/v2"
	caddyfile "github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/ionos"
	"github.com/libdns/libdns"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct {
	*ionos.Provider
	// OverrideDomain specifies the actual zone for DNS-01 challenge records.
	// Use this when delegating ACME DNS validation via CNAME records.
	OverrideDomain string `json:"override_domain,omitempty"`
}

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.ionos",
		New: func() caddy.Module { return &Provider{Provider: new(ionos.Provider)} },
	}
}

// Before using the provider config, resolve placeholders in the API token.
// Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.AuthAPIToken = repl.ReplaceAll(p.AuthAPIToken, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	ionos [<api_token>] {
//	    api_token <api_token>
//	    override_domain <ionos_domain>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.AuthAPIToken = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_token":
				if p.AuthAPIToken != "" {
					return d.Err("API token already set")
				}
				if d.NextArg() {
					p.AuthAPIToken = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "override_domain":
				if p.OverrideDomain != "" {
					return d.Err("Override domain already set")
				}
				if d.NextArg() {
					p.OverrideDomain = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.AuthAPIToken == "" {
		return d.Err("missing API token")
	}
	return nil
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	if p.OverrideDomain != "" {
		zone = p.OverrideDomain
	}
	return p.Provider.GetRecords(ctx, zone)
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	if p.OverrideDomain != "" {
		zone = p.OverrideDomain
	}
	return p.Provider.AppendRecords(ctx, zone, records)
}

// DeleteRecords deletes the records from the zone.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	if p.OverrideDomain != "" {
		zone = p.OverrideDomain
	}
	return p.Provider.DeleteRecords(ctx, zone, records)
}

// SetRecords sets the records in the zone, either by updating existing records
// or creating new ones. It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	if p.OverrideDomain != "" {
		zone = p.OverrideDomain
	}
	return p.Provider.SetRecords(ctx, zone, records)
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
