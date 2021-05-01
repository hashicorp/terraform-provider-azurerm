package parse

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"strings"
)

type InferenceClusterId struct {
	Name                 string
	ResourceGroup        string
	InferenceClusterName string
}

func InferenceClusterID(input string) (*InferenceClusterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	inference_cluster := InferenceClusterId{
		ResourceGroup: id.ResourceGroup,
	}

	if inference_cluster.InferenceClusterName, err = id.PopSegment("computes"); err != nil {
		return nil, err
	}

	if inference_cluster.Name, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	fmt.Printf("Debug: id = %q", id)

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &inference_cluster, nil
}

type KubernetesClusterId struct {
	SubscriptionId     string
	ResourceGroup      string
	ManagedClusterName string
}

func NewClusterID(subscriptionId, resourceGroup, managedClusterName string) KubernetesClusterId {
	return KubernetesClusterId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		ManagedClusterName: managedClusterName,
	}
}

func (id KubernetesClusterId) String() string {
	segments := []string{
		fmt.Sprintf("Managed Cluster Name %q", id.ManagedClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Cluster", segmentsStr)
}

func (id KubernetesClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedClusterName)
}

// KubernetesClusterId parses a Cluster ID into an KubernetesClusterId struct
func KubernetesClusterID(input string) (*KubernetesClusterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := KubernetesClusterId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ManagedClusterName, err = id.PopSegment("managedClusters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
