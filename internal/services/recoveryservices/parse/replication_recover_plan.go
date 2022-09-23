package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ReplicationRecoverPlanId struct {
	SubscriptionId              string
	ResourceGroup               string
	VaultName                   string
	ReplicationRecoveryPlanName string
}

func NewReplicationRecoverPlanID(subscriptionId, resourceGroup, vaultName, replicationRecoveryPlanName string) ReplicationRecoverPlanId {
	return ReplicationRecoverPlanId{
		SubscriptionId:              subscriptionId,
		ResourceGroup:               resourceGroup,
		VaultName:                   vaultName,
		ReplicationRecoveryPlanName: replicationRecoveryPlanName,
	}
}

func (id ReplicationRecoverPlanId) String() string {
	segments := []string{
		fmt.Sprintf("Replication Recovery Plan Name %q", id.ReplicationRecoveryPlanName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Replication Recover Plan", segmentsStr)
}

func (id ReplicationRecoverPlanId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationRecoveryPlans/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.ReplicationRecoveryPlanName)
}

// ReplicationRecoverPlanID parses a ReplicationRecoverPlan ID into an ReplicationRecoverPlanId struct
func ReplicationRecoverPlanID(input string) (*ReplicationRecoverPlanId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ReplicationRecoverPlanId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VaultName, err = id.PopSegment("vaults"); err != nil {
		return nil, err
	}
	if resourceId.ReplicationRecoveryPlanName, err = id.PopSegment("replicationRecoveryPlans"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
