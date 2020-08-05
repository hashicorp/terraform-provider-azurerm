package parse

import (
    "fmt"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)
type AvsHcxEnterpriseSiteId struct {
    ResourceGroup string
    PrivateCloudName string
    Name string
}

func AvsHcxEnterpriseSiteID(input string) (*AvsHcxEnterpriseSiteId, error) {
    id, err := azure.ParseAzureResourceID(input)
    if err != nil {
        return nil, fmt.Errorf("parsing avsHcxEnterpriseSite ID %q: %+v", input, err)
    }

    avsHcxEnterpriseSite := AvsHcxEnterpriseSiteId{
        ResourceGroup: id.ResourceGroup,
    }
    if avsHcxEnterpriseSite.PrivateCloudName, err = id.PopSegment("privateClouds"); err != nil {
        return nil, err
    }
    if avsHcxEnterpriseSite.Name, err = id.PopSegment("hcxEnterpriseSites"); err != nil {
        return nil, err
    }
    if err := id.ValidateNoEmptySegments(input); err != nil {
        return nil, err
    }

    return &avsHcxEnterpriseSite, nil
}
