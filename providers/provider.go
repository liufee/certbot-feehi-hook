package providers

type Provider interface {
	ResolveDomainName(dnsType string, pr string, value string) (bool, error)
	DeleteResolveDomainName(dnsType string, pr string) (bool, error)
}
