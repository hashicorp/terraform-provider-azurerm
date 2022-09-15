package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskencryptionsets"
)

func retrieveDiskEncryptionSetEncryptionType(ctx context.Context, client *diskencryptionsets.DiskEncryptionSetsClient, diskEncryptionSetId string) (*diskencryptionsets.DiskEncryptionSetType, error) {
	id, err := diskencryptionsets.ParseDiskEncryptionSetID(diskEncryptionSetId)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	var encryptionType *diskencryptionsets.DiskEncryptionSetType

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.EncryptionType != nil {
			encryptionType = props.EncryptionType
		}
	}

	if encryptionType == nil {
		return nil, fmt.Errorf("retrieving %s: EncryptionType was nil", *id)
	}

	return encryptionType, nil
}
