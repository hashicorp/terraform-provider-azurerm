package containers

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KubernetesClusterID struct {
	Name          string
	ResourceGroup string

	ID azure.ResourceID
}

func KubernetesClusterIDSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: ValidateKubernetesClusterID,
	}
}

func ParseKubernetesClusterID(id string) (*KubernetesClusterID, error) {
	clusterId, err := azure.ParseAzureResourceID(id)
	if err != nil {
		return nil, err
	}

	resourceGroup := clusterId.ResourceGroup
	if resourceGroup == "" {
		return nil, fmt.Errorf("%q is missing a Resource Group", id)
	}

	clusterName := clusterId.Path["managedClusters"]
	if clusterName == "" {
		return nil, fmt.Errorf("%q is missing the `managedClusters` segment", id)
	}

	output := KubernetesClusterID{
		Name:          clusterName,
		ResourceGroup: resourceGroup,
		ID:            *clusterId,
	}
	return &output, nil
}

func ValidateKubernetesClusterID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	id, err := azure.ParseAzureResourceID(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a Resource Id: %v", v, err))
	}

	if id != nil {
		if id.Path["managedClusters"] == "" {
			errors = append(errors, fmt.Errorf("The 'managedClusters' segment is missing from Resource ID %q", v))
		}
	}

	return warnings, errors
}
