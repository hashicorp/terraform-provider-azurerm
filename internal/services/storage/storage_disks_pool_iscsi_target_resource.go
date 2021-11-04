package storage

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/storagepool/mgmt/2021-08-01/storagepool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"regexp"
	"time"
)

type DisksPoolIscsiTargetResource struct{}

var _ sdk.ResourceWithUpdate = DisksPoolIscsiTargetResource{}

type DiskPoolIscsiTargetModel struct {
	ACLMode     string              `tfschema:"acl_mode"`
	DisksPoolId string              `tfschema:"disks_pool_id"`
	Endpoints   []string            `tfschema:"endpoints"` // List of private IPv4 addresses to connect to the iSCSI Target.
	Luns        []DisksPoolIscsiLun `tfschema:"lun"`
	Name        string              `tfschema:"name"`
	Port        int                 `tfschema:"port"`       // The port used by iSCSI Target portal group.
	TargetIqn   string              `tfschema:"target_iqn"` // iSCSI Target IQN (iSCSI Qualified Name); example: \"iqn.2005-03.org.iscsi:server\".
}

type DisksPoolIscsiLun struct {
	Name          string `tfschema:"name"`            // User defined name for iSCSI LUN; example: \"lun0\" min:1 max:90
	ManagedDiskId string `tfschema:"managed_disk_id"` // Azure Resource ID of the Managed Disk.
	Number        int    `tfschema:"number"`          // Specifies the Logical Unit Number of the iSCSI LUN. readonly
}

func (d DisksPoolIscsiTargetResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"acl_mode": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(possibleIscsiTargetACLModeValues(), false),
		},
		"disks_pool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageDisksPoolID,
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(5, 223),
				validation.StringMatch(
					regexp.MustCompile("[a-z\\d.\\-]*[a-z\\d]$"),
					"The iSCSI target name can only contain lowercase letters, numbers, periods, or hyphens.",
				),
			),
		},

		"lun": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.All(
							validation.StringMatch(
								regexp.MustCompile("^[\\dA-Za-z-_.]+[\\dA-Za-z]$"),
								"supported characters include [0-9A-Za-z-_.]; name should end with an alphanumeric character"),
							validation.StringLenBetween(1, 90),
						),
					},
					"managed_disk_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: computeValidate.ManagedDiskID,
					},
					"number": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"target_iqn": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.IQN,
		},
	}
}

func (d DisksPoolIscsiTargetResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"endpoints": {
			Type: pluginsdk.TypeList,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Computed: true,
		},
		"port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
	}
}

func (d DisksPoolIscsiTargetResource) ModelObject() interface{} {
	return &DisksPoolIscsiTargetResource{}
}

func (d DisksPoolIscsiTargetResource) ResourceType() string {
	return "azurerm_storage_disks_pool_iscsi_target"
}

func (d DisksPoolIscsiTargetResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			m := DiskPoolIscsiTargetModel{}
			err := metadata.Decode(&m)
			if err != nil {
				return err
			}
			poolId, err := parse.StorageDisksPoolID(m.DisksPoolId)
			if err != nil {
				return err
			}
			if poolId.SubscriptionId != metadata.Client.Account.SubscriptionId {
				return fmt.Errorf("Disks Pool subscription id %q is different from provider's subscription", poolId.SubscriptionId)
			}
			id := parse.NewStorageDisksPoolISCSITargetID(poolId.SubscriptionId, poolId.ResourceGroup, poolId.DiskPoolName, m.Name)
			client := metadata.Client.Storage.DisksPoolIscsiTargetClient
			// TODO: test attach disk while creating iscsi target
			locks.ByID(poolId.ID())
			defer locks.UnlockByID(poolId.ID())

			existing, err := client.Get(ctx, id.ResourceGroup, id.DiskPoolName, id.IscsiTargetName)
			notExistingResp := utils.ResponseWasNotFound(existing.Response)
			if err != nil && !notExistingResp {
				return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
			}
			if !notExistingResp {
				return metadata.ResourceRequiresImport(d.ResourceType(), id)
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.DiskPoolName, id.IscsiTargetName, storagepool.IscsiTargetCreate{
				IscsiTargetCreateProperties: &storagepool.IscsiTargetCreateProperties{
					ACLMode:   storagepool.IscsiTargetACLMode(m.ACLMode),
					TargetIqn: &m.TargetIqn,
					Luns:      expandDisksPoolIscsiTargetLuns(m.Luns),
				},
				Name: utils.String(m.Name),
			})
			if err != nil {
				return err
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of DisksPool iscsi target %q : %+v", id.ID(), err)
			}
			metadata.SetID(id)
			return nil
		}}
}

func (d DisksPoolIscsiTargetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.StorageDisksPoolISCSITargetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			client := metadata.Client.Storage.DisksPoolIscsiTargetClient
			resp, err := client.Get(ctx, id.ResourceGroup, id.DiskPoolName, id.IscsiTargetName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			m := DiskPoolIscsiTargetModel{
				ACLMode:     string(resp.ACLMode),
				DisksPoolId: id.ID(),
				Name:        id.IscsiTargetName,
			}
			if resp.IscsiTargetProperties == nil {
				return metadata.Encode(&m)
			}
			if endpoints := resp.IscsiTargetProperties.Endpoints; endpoints != nil {
				m.Endpoints = *endpoints
			}
			if luns := resp.IscsiTargetProperties.Luns; luns != nil {
				m.Luns = flattenDisksPoolIscsiTargetLuns(resp.Luns)
			}
			if port := resp.IscsiTargetProperties.Port; port != nil {
				m.Port = int(*port)
			}
			m.TargetIqn = utils.NormalizeNilableString(resp.TargetIqn)
			return metadata.Encode(&m)
		},
	}
}

func (d DisksPoolIscsiTargetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.StorageDisksPoolISCSITargetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			client := metadata.Client.Storage.DisksPoolIscsiTargetClient
			future, err := client.Delete(ctx, id.ResourceGroup, id.DiskPoolName, id.IscsiTargetName)
			if err != nil {
				return err
			}
			err = future.WaitForCompletionRef(ctx, client.Client)
			if err != nil {
				return fmt.Errorf("waiting for deletion of DisksPool iscsi target %q: %+v", id.ID(), err)
			}
			return nil
		},
	}
}

func (d DisksPoolIscsiTargetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.StorageDisksPoolISCSITargetID
}

func (d DisksPoolIscsiTargetResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			r := metadata.ResourceData
			id, err := parse.StorageDisksPoolISCSITargetID(r.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			client := metadata.Client.Storage.DisksPoolIscsiTargetClient
			m := DiskPoolIscsiTargetModel{}
			err = metadata.Decode(&m)
			if err != nil {
				return err
			}

			patch := storagepool.IscsiTargetUpdate{
				IscsiTargetUpdateProperties: &storagepool.IscsiTargetUpdateProperties{},
			}
			if r.HasChange("lun") {
				patch.IscsiTargetUpdateProperties.Luns = expandDisksPoolIscsiTargetLuns(m.Luns)
			}

			future, err := client.Update(ctx, id.ResourceGroup, id.DiskPoolName, id.IscsiTargetName, patch)
			if err != nil {
				return err
			}
			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of DiskPool iscsi taraget %q : %+v", id.ID(), err)
			}
			return nil
		},
	}
}

func expandDisksPoolIscsiTargetLuns(modelLuns []DisksPoolIscsiLun) *[]storagepool.IscsiLun {
	luns := make([]storagepool.IscsiLun, 0)
	for i := 0; i < len(modelLuns); i++ {
		luns = append(luns, storagepool.IscsiLun{
			Name:                       &modelLuns[i].Name,
			ManagedDiskAzureResourceID: &modelLuns[i].ManagedDiskId,
			Lun:                        utils.Int32(int32(modelLuns[i].Number)),
		})
	}
	return &luns
}

func flattenDisksPoolIscsiTargetLuns(respLuns *[]storagepool.IscsiLun) []DisksPoolIscsiLun {
	var luns []DisksPoolIscsiLun
	if respLuns != nil {
		for i := 0; i < len(*respLuns); i++ {
			lun := (*respLuns)[i]
			l := DisksPoolIscsiLun{}
			if lun.Name != nil {
				l.Name = *lun.Name
			}
			if lun.ManagedDiskAzureResourceID != nil {
				l.ManagedDiskId = *lun.ManagedDiskAzureResourceID
			}
			if lun.Lun != nil {
				l.Number = int(*lun.Lun)
			}
			luns = append(luns, l)
		}
	}
	return luns
}

func possibleIscsiTargetACLModeValues() []string {
	var values []string
	for _, v := range storagepool.PossibleIscsiTargetACLModeValues() {
		values = append(values, string(v))
	}
	return values
}
