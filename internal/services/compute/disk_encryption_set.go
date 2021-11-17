package compute

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
)

// retrieveDiskEncryptionSetEncryptionType returns encryption type of the disk encryption set
func retrieveDiskEncryptionSetEncryptionType(ctx context.Context, client *compute.DiskEncryptionSetsClient, diskEncryptionSetId string) (*string, error) {
	diskEncryptionSet, err := parse.DiskEncryptionSetID(diskEncryptionSetId)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, diskEncryptionSet.ResourceGroup, diskEncryptionSet.Name)
	if err != nil {
		return nil, err
	}

	if properties := resp.EncryptionSetProperties; properties != nil {
		encryptionType := string(properties.EncryptionType)
		return &encryptionType, nil
	} else {
		return nil, fmt.Errorf("could not get EncryptionSetProperties")
	}
}
