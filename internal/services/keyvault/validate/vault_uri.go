package validate

import (
	"fmt"
	"net/url"
	"strings"
)

func VaultURI(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	url, err := url.ParseRequestURI(value)
	if err != nil || url == nil {
		errors = append(errors, fmt.Errorf("%q may only contain a valid vault Uri", k))
		return
	}

	// https://learn.microsoft.com/en-us/azure/key-vault/general/about-keys-secrets-certificates#dns-suffixes-for-base-url
	dnsSuffixes := []string{
		".vault.azure.net",
		".vault.azure.cn",
		".vault.usgovcloudapi.net",
		".vault.microsoftazure.de",
	}

	for _, dnsSuffix := range dnsSuffixes {
		if strings.HasSuffix(url.Host, dnsSuffix) {
			return
		}
	}

	errors = append(errors, fmt.Errorf("%q may only contain a valid vault Uri", k))

	return
}
