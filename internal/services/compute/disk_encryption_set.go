package compute

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
)

// retrieveDiskEncryptionSetEncryptionType returns encryption type of the disk encryption set
func retrieveDiskEncryptionSetEncryptionType(ctx context.Context, client *compute.DiskEncryptionSetsClient, diskEncryptionSetId string) (*compute.EncryptionType, error) {
	diskEncryptionSet, err := parse.DiskEncryptionSetID(diskEncryptionSetId)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, diskEncryptionSet.ResourceGroup, diskEncryptionSet.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *diskEncryptionSet, err)
	}

	var encryptionType *compute.EncryptionType
	if props := resp.EncryptionSetProperties; props != nil && string(props.EncryptionType) != "" {
		v := compute.EncryptionType(props.EncryptionType)
		encryptionType = &v
	}

	if encryptionType == nil {
		return nil, fmt.Errorf("retrieving %s: EncryptionType was nil", *diskEncryptionSet)
	}

	return encryptionType, nil
}
