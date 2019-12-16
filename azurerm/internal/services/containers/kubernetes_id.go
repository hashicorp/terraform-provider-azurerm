package containers

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KubernetesClusterID struct {
	Name          string
	ResourceGroup string
}

func KubernetesClusterIDSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: ValidateKubernetesClusterID,
	}
}

func ParseKubernetesClusterID(input string) (*KubernetesClusterID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	cluster := KubernetesClusterID{
		ResourceGroup: id.ResourceGroup,
	}

	if cluster.Name, err = id.PopSegment("managedClusters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &cluster, nil
}

func ValidateKubernetesClusterID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := ParseKubernetesClusterID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a Resource Id: %v", v, err))
	}

	return warnings, errors
}
