package accounts

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SubDomainType string

const (
	BlobSubDomainType          SubDomainType = "blob"
	DataLakeStoreSubDomainType SubDomainType = "dfs"
	FileSubDomainType          SubDomainType = "file"
	QueueSubDomainType         SubDomainType = "queue"
	TableSubDomainType         SubDomainType = "table"
)

func PossibleValuesForSubDomainType() []SubDomainType {
	return []SubDomainType{
		BlobSubDomainType,
		DataLakeStoreSubDomainType,
		FileSubDomainType,
		QueueSubDomainType,
		TableSubDomainType,
	}
}

// TODO: update this to implement `resourceids.ResourceId` once
// https://github.com/hashicorp/go-azure-helpers/issues/187 is fixed
var _ resourceids.Id = AccountId{}

type AccountId struct {
	AccountName   string
	ZoneName      *string
	SubDomainType SubDomainType
	DomainSuffix  string
	IsEdgeZone    bool
}

func (a AccountId) ID() string {
	components := []string{
		a.AccountName,
	}
	if a.IsEdgeZone {
		// Storage Accounts hosted in an Edge Zone
		//   `{accountname}.{component}.{edgezone}.edgestorage.azure.net`
		components = append(components, string(a.SubDomainType))
		if a.ZoneName != nil {
			components = append(components, *a.ZoneName)
		}
	} else {
		// Storage Accounts using a DNS Zone
		//   `{accountname}.{dnszone}.{component}.storage.azure.net`
		// or a Regular Storage Account
		//   `{accountname}.{component}.core.windows.net`
		if a.ZoneName != nil {
			components = append(components, *a.ZoneName)
		}
		components = append(components, string(a.SubDomainType))
	}
	components = append(components, a.DomainSuffix)
	return fmt.Sprintf("https://%s", strings.Join(components, "."))
}

func (a AccountId) String() string {
	components := []string{
		fmt.Sprintf("IsEdgeZone %t", a.IsEdgeZone),
		fmt.Sprintf("ZoneName %q", pointer.From(a.ZoneName)),
		fmt.Sprintf("Subdomain Type %q", string(a.SubDomainType)),
		fmt.Sprintf("DomainSuffix %q", a.DomainSuffix),
	}
	return fmt.Sprintf("Account %q (%s)", a.AccountName, strings.Join(components, " / "))
}

func ParseAccountID(input, regularDomainSuffix string) (*AccountId, error) {
	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a URL: %s", input, err)
	}

	// There's 3 different types of Storage Account ID:
	// 1. Regular ol' Storage Accounts
	//   `{name}.{component}.core.windows.net` (e.g. `example1.blob.core.windows.net`)
	// 2. Storage Accounts using a DNS Zone
	//   `{accountname}.{dnszone}.{component}.storage.azure.net`
	// 3. Storage Accounts hosted in an Edge Zone
	//   `{accountname}.{component}.{edgezone}.edgestorage.azure.net`
	// since both `dnszone` and `edgezone` are the only two identifiers, we need to check if `domainSuffix` includes `edge`
	// to know how to treat these
	type domainType int
	const (
		unknownDomain domainType = iota
		regularDomain
		dnsZoneDomain
		edgeZoneDomain
	)

	validDomainSuffixes := map[domainType]string{
		regularDomain:  regularDomainSuffix,
		dnsZoneDomain:  "storage.azure.net",
		edgeZoneDomain: "edgestorage.azure.net",
	}

	dt := unknownDomain
	for k, suffix := range validDomainSuffixes {
		if strings.HasSuffix(uri.Host, suffix) {
			dt = k
			break
		}
	}
	if dt == unknownDomain {
		return nil, fmt.Errorf("invalid domain suffix of account %q", uri.Host)
	}

	domainSuffix := validDomainSuffixes[dt]
	hostName := strings.TrimSuffix(uri.Host, fmt.Sprintf(".%s", domainSuffix))
	components := strings.Split(hostName, ".")
	accountId := AccountId{
		DomainSuffix: domainSuffix,
		IsEdgeZone:   strings.Contains(strings.ToLower(domainSuffix), "edge"),
	}

	switch dt {
	case regularDomain:
		if l := len(components); l != 2 {
			return nil, fmt.Errorf("expect 2 uri components before the domain suffix %q, got=%d", domainSuffix, l)
		}
		accountId.AccountName = components[0]
		subDomainType, err := parseSubDomainType(components[1])
		if err != nil {
			return nil, err
		}
		accountId.SubDomainType = *subDomainType
		return &accountId, nil
	case dnsZoneDomain:
		if l := len(components); l != 3 {
			return nil, fmt.Errorf("expect 3 uri components before the domain suffix %q, got=%d", domainSuffix, l)
		}
		accountId.AccountName = components[0]
		subDomainType, err := parseSubDomainType(components[2])
		if err != nil {
			return nil, err
		}
		accountId.SubDomainType = *subDomainType
		accountId.ZoneName = pointer.To(components[1])
		return &accountId, nil
	case edgeZoneDomain:
		if l := len(components); l != 3 {
			return nil, fmt.Errorf("expect 3 uri components before the domain suffix %q, got=%d", domainSuffix, l)
		}
		accountId.AccountName = components[0]
		subDomainType, err := parseSubDomainType(components[1])
		if err != nil {
			return nil, err
		}
		accountId.SubDomainType = *subDomainType
		accountId.ZoneName = pointer.To(components[2])
		return &accountId, nil
	}

	return nil, fmt.Errorf("unknown storage account domain type %q", input)
}

func parseSubDomainType(input string) (*SubDomainType, error) {
	for _, k := range PossibleValuesForSubDomainType() {
		if strings.EqualFold(input, string(k)) {
			return pointer.To(k), nil
		}
	}

	return nil, fmt.Errorf("expected the subdomain type to be one of [%+v] but got %q", PossibleValuesForSubDomainType(), input)
}
