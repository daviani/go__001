package scanner

type Scanner interface {
	Scan(domain string) string
	Name() string
}
type DNSScanner struct{}
type SSLScanner struct{}

func (d DNSScanner) Name() string { return "dns" }
func (s SSLScanner) Name() string { return "ssl" }

func (d DNSScanner) Scan(domain string) string {
	return "DNS result for: " + domain
}

func (s SSLScanner) Scan(domain string) string {
	return "SSL result for: " + domain
}
