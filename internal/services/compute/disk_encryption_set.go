package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/tombuildsstuff/kermit/sdk/compute/2022-08-01/compute"
)

// retrieveDiskEncryptionSetEncryptionType returns encryption type of the disk encryption set
func retrieveDiskEncryptionSetEncryptionType(ctx context.Context, client *compute.DiskEncryptionSetsClient, diskEncryptionSetId string) (*disks.EncryptionType, error) {
	diskEncryptionSet, err := parse.DiskEncryptionSetID(diskEncryptionSetId)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, diskEncryptionSet.ResourceGroup, diskEncryptionSet.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *diskEncryptionSet, err)
	}

	var encryptionType *disks.EncryptionType
	if props := resp.EncryptionSetProperties; props != nil && string(props.EncryptionType) != "" {
		v := disks.EncryptionType(props.EncryptionType)
		encryptionType = &v
	}

	if encryptionType == nil {
		return nil, fmt.Errorf("retrieving %s: EncryptionType was nil", *diskEncryptionSet)
	}

	return encryptionType, nil
}
