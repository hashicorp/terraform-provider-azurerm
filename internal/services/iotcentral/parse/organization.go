package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = OrganizationId{}

type OrganizationId struct {
	DomainSuffix   string
	SubDomain      string
	OrganizationId string
}

func NewOrganizationID(subDomain string, domainSuffix string, id string) (*OrganizationId, error) {
	return &OrganizationId{
		DomainSuffix:   domainSuffix,
		SubDomain:      subDomain,
		OrganizationId: id,
	}, nil
}

func (id OrganizationId) ID() string {
	return fmt.Sprintf("https://%s.%s/api/organizations/%s", id.SubDomain, id.DomainSuffix, id.OrganizationId)
}

func (id OrganizationId) String() string {
	components := []string{
		fmt.Sprintf("DomainSuffix %q", id.DomainSuffix),
		fmt.Sprintf("Sub Domain %q", id.SubDomain),
		fmt.Sprintf("Id %q", id.OrganizationId),
	}
	return fmt.Sprintf("Iot Central OrganizationId %s", strings.Join(components, " / "))
}

func ParseOrganizationID(id string, subDomain string, domainSuffix string) (*OrganizationId, error) {
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("cannot parse azure iot central organization ID: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	if len(components) != 3 {
		return nil, fmt.Errorf("iot central organization should have 3 segments, found %d segment(s) in %q", len(components), id)
	}

	apiString := components[0]
	if apiString != "api" {
		return nil, fmt.Errorf("iot central organization should have api as first segment, found %q", apiString)
	}

	organizationsString := components[1]
	if organizationsString != "organizations" {
		return nil, fmt.Errorf("iot central organization should have organizations as second segment, found %q", organizationsString)
	}

	parsedOrganizationId := components[2]

	parsedSubDomain := strings.Split(idURL.Host, ".")[0]
	parsedDomainSuffix := strings.Split(idURL.Host, ".")[1]

	if subDomain != "" { // subDomain is empty when importing
		if parsedSubDomain != subDomain {
			return nil, fmt.Errorf("iot central organization subdomain should be %q, got %q", subDomain, parsedSubDomain)
		}
	}

	if parsedDomainSuffix != domainSuffix {
		return nil, fmt.Errorf("iot central organization domain suffix should be %q, got %q", domainSuffix, parsedDomainSuffix)
	}

	organizationId := OrganizationId{
		DomainSuffix:   parsedDomainSuffix,
		SubDomain:      parsedSubDomain,
		OrganizationId: parsedOrganizationId,
	}

	return &organizationId, nil
}
