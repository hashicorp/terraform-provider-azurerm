package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ManagedHSMKeyId struct {
	SubscriptionId string
	ResourceGroup  string
	ManagedHSMName string
	KeyName        string
	VersionName    string
}

func NewManagedHSMKeyID(subscriptionId, resourceGroup, managedHSMName, keyName, versionName string) ManagedHSMKeyId {
	return ManagedHSMKeyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ManagedHSMName: managedHSMName,
		KeyName:        keyName,
		VersionName:    versionName,
	}
}

func (id ManagedHSMKeyId) String() string {
	segments := []string{
		fmt.Sprintf("Version Name %q", id.VersionName),
		fmt.Sprintf("Key Name %q", id.KeyName),
		fmt.Sprintf("Managed H S M Name %q", id.ManagedHSMName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed H S M Key", segmentsStr)
}

func (id ManagedHSMKeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/managedHSMs/%s/keys/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedHSMName, id.KeyName, id.VersionName)
}

// ManagedHSMKeyID parses a ManagedHSMKey ID into an ManagedHSMKeyId struct
func ManagedHSMKeyID(input string) (*ManagedHSMKeyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ManagedHSMKeyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ManagedHSMName, err = id.PopSegment("managedHSMs"); err != nil {
		return nil, err
	}
	if resourceId.KeyName, err = id.PopSegment("keys"); err != nil {
		return nil, err
	}
	if resourceId.VersionName, err = id.PopSegment("versions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
