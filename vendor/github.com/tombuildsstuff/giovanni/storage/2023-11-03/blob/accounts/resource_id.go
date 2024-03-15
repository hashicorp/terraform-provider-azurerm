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

func ParseAccountID(input, domainSuffix string) (*AccountId, error) {
	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a URL: %s", input, err)
	}

	if !strings.HasSuffix(uri.Host, domainSuffix) {
		return nil, fmt.Errorf("expected the account %q to use a domain suffix of %q", uri.Host, domainSuffix)
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

	hostName := strings.TrimSuffix(uri.Host, fmt.Sprintf(".%s", domainSuffix))
	components := strings.Split(hostName, ".")
	accountId := AccountId{
		DomainSuffix: domainSuffix,
		IsEdgeZone:   strings.Contains(strings.ToLower(domainSuffix), "edge"),
	}

	if len(components) == 2 {
		// this will be a regular Storage Account (e.g. `example1.blob.core.windows.net`)
		accountId.AccountName = components[0]
		subDomainType, err := parseSubDomainType(components[1])
		if err != nil {
			return nil, err
		}
		accountId.SubDomainType = *subDomainType
		return &accountId, nil
	}

	if len(components) == 3 {
		// This can either be a Zone'd Storage Account or a Storage Account within an Edge Zone
		accountName := ""
		subDomainTypeRaw := ""
		zone := ""
		if accountId.IsEdgeZone {
			// `{accountname}.{component}.{edgezone}.edgestorage.azure.net`
			accountName = components[0]
			subDomainTypeRaw = components[1]
			zone = components[2]
		} else {
			// `{accountname}.{dnszone}.{component}.storage.azure.net`
			accountName = components[0]
			zone = components[1]
			subDomainTypeRaw = components[2]
		}

		accountId.AccountName = accountName
		subDomainType, err := parseSubDomainType(subDomainTypeRaw)
		if err != nil {
			return nil, err
		}
		accountId.SubDomainType = *subDomainType
		accountId.ZoneName = pointer.To(zone)
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
