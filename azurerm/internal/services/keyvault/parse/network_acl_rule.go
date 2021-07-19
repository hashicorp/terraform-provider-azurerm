package parse

import (
	"fmt"
	"strings"

	commonValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
)

type NetworkAclRuleType string

const (
	CIDR           NetworkAclRuleType = "cidr"
	Ip             NetworkAclRuleType = "ip"
	VirtualNetwork NetworkAclRuleType = "virtualNetwork"
)

type NetworkAclRuleId struct {
	VaultId VaultId
	Type    NetworkAclRuleType
	Source  string
}

func NewNetworkAclRuleId(vaultId, source string) (*NetworkAclRuleId, error) {
	myvaultId, err := VaultID(vaultId)
	if err != nil {
		return nil, err
	}

	ruleType, err := parseNetworkAclRuleType(source)
	if err != nil {
		return nil, err
	}

	return &NetworkAclRuleId{
		VaultId: *myvaultId,
		Source:  source,
		Type:    ruleType,
	}, nil
}

func (id NetworkAclRuleId) String() string {
	segments := []string{
		fmt.Sprintf("VaultID %q", id.VaultId),
		fmt.Sprintf("Source %q", id.Source),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Vault Network ACL Rule", segmentsStr)
}

func (id NetworkAclRuleId) ID() string {
	fmtString := "%s|%s"
	return fmt.Sprintf(fmtString, id.VaultId.ID(), id.Source)
}

func (id NetworkAclRuleId) Rule() string {
	if id.Type == Ip {
		return fmt.Sprintf("%s/32", id.Source)
	}
	return id.Source
}

// NetworkAclRuleId is a pseudo ID for storing Source parameter as this it not retrievable from API
// It is formed of the Azure Resource ID for the Vault and the Source it is created against
func NetworkAclRuleID(input string) (*NetworkAclRuleId, error) {
	parts := strings.Split(input, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("could not parse Network ACL Rule ID, invalid format %q", input)
	}

	vaultId, err := VaultID(parts[0])
	if err != nil {
		return nil, err
	}

	source := parts[1]

	ruleType, err := parseNetworkAclRuleType(source)
	if err != nil {
		return nil, err
	}

	return &NetworkAclRuleId{
		VaultId: *vaultId,
		Source:  parts[1],
		Type:    ruleType,
	}, nil
}

func parseNetworkAclRuleType(source string) (NetworkAclRuleType, error) {
	var errorList []error

	_, errorList = commonValidate.IPv4Address(source, "ip_rule")
	if len(errorList) == 0 {
		return Ip, nil
	}

	_, errorList = commonValidate.CIDR(source, "ip_rule")
	if len(errorList) == 0 {
		return CIDR, nil
	}

	_, errorList = networkValidate.SubnetID(source, "virtual_network_subnet_id")
	if len(errorList) == 0 {
		return VirtualNetwork, nil
	}

	return "", fmt.Errorf("could not determine rule type of %s", source)
}
