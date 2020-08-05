package parse

import (
    "fmt"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)
type AvsClusterId struct {
    ResourceGroup string
    PrivateCloudName string
    Name string
}

func AvsClusterID(input string) (*AvsClusterId, error) {
    id, err := azure.ParseAzureResourceID(input)
    if err != nil {
        return nil, fmt.Errorf("parsing avsCluster ID %q: %+v", input, err)
    }

    avsCluster := AvsClusterId{
        ResourceGroup: id.ResourceGroup,
    }
    if avsCluster.PrivateCloudName, err = id.PopSegment("privateClouds"); err != nil {
        return nil, err
    }
    if avsCluster.Name, err = id.PopSegment("clusters"); err != nil {
        return nil, err
    }
    if err := id.ValidateNoEmptySegments(input); err != nil {
        return nil, err
    }

    return &avsCluster, nil
}
