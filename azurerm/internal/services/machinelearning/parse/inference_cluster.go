package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type InferenceClusterId struct {
	SubscriptionId string
	ResourceGroup  string
	WorkspaceName  string
	ComputeName    string
}

func NewInferenceClusterID(subscriptionId, resourceGroup, workspaceName, computeName string) InferenceClusterId {
	return InferenceClusterId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WorkspaceName:  workspaceName,
		ComputeName:    computeName,
	}
}

func (id InferenceClusterId) String() string {
	segments := []string{
		fmt.Sprintf("Compute Name %q", id.ComputeName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Inference Cluster", segmentsStr)
}

func (id InferenceClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MachineLearningServices/workspaces/%s/computes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.ComputeName)
}

// InferenceClusterID parses a InferenceCluster ID into an InferenceClusterId struct
func InferenceClusterID(input string) (*InferenceClusterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := InferenceClusterId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if resourceId.ComputeName, err = id.PopSegment("computes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
