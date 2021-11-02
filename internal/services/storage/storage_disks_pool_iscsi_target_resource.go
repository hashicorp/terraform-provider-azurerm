package storage

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/storagepool/mgmt/2021-08-01/storagepool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
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
	ACLMode           string                    `tfschema:"acl_mode"`
	DisksPoolId       string                    `tfschema:"disks_pool_id"`
	Endpoints         []string                  `tfschema:"endpoints"` // List of private IPv4 addresses to connect to the iSCSI Target.
	Luns              []DisksPoolIscsiLun       `tfschema:"luns"`
	ManagedBy         string                    `tfschema:"managed_by"`
	ManagedByExtended []string                  `tfschema:"managed_by_extended"`
	Name              string                    `tfschema:"name"`
	Port              int                       `tfschema:"port"` // The port used by iSCSI Target portal group.
	StaticAcls        []DisksPoolISCSITargetAcl `tfschema:"static_acls"`
	TargetIqn         string                    `tfschema:"target_iqn"` // iSCSI Target IQN (iSCSI Qualified Name); example: \"iqn.2005-03.org.iscsi:server\".

	// TODO: Should we add this?
	// Session []string //List of identifiers for active sessions on the iSCSI target
}

type DisksPoolISCSITargetAcl struct {
	InitiatorIqn string   `tfschema:"initiator_iqn"` // iSCSI initiator IQN (iSCSI Qualified Name); example: \"iqn.2005-03.org.iscsi:client\".
	MappedLuns   []string `tfschema:"mapped_luns"`   // List of LUN names mapped to the ACL.
}

type DisksPoolIscsiLun struct {
	Name                       string `tfschema:"name"`                           // User defined name for iSCSI LUN; example: \"lun0\" min:1 max:90
	ManagedDiskAzureResourceId string `tfschema:"managed_disk_azure_resource_id"` // Azure Resource ID of the Managed Disk.
	Lun                        int    `tfschema:"lun"`                            // Specifies the Logical Unit Number of the iSCSI LUN. readonly
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
		"luns": { // Allow empty slice!
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.All(
							validation.StringIsNotEmpty,
							validation.StringLenBetween(1, 90),
						),
					},
					"managed_disk_azure_resource_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: computeValidate.ManagedDiskID,
					},
					"lun": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"managed_by_extended": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
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
		"static_acls": { // TODO: How to set static acl?
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"initiator_iqn": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.IQN,
					},
					"mapped_luns": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &schema.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
		"target_iqn": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
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
		"managed_by": {
			Type:     pluginsdk.TypeString,
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
			targetModel := DiskPoolIscsiTargetModel{}
			err := metadata.Decode(&targetModel)
			if err != nil {
				return err
			}
			poolId, err := parse.StorageDisksPoolID(targetModel.DisksPoolId)
			if err != nil {
				return err
			}
			id := parse.NewStorageDisksPoolISCSITargetID(poolId.SubscriptionId, poolId.ResourceGroup, poolId.DiskPoolName, targetModel.Name)

			locks.ByID(poolId.ID())
			defer locks.UnlockByID(poolId.ID())

			client := metadata.Client.Storage.DisksPoolIscsiTargetClient

			if metadata.ResourceData.IsNewResource() {
				existing, err := client.Get(ctx, id.ResourceGroup, id.DiskPoolName, id.IscsiTargetName)
				notExistingResp := utils.ResponseWasNotFound(existing.Response)
				if err != nil && !notExistingResp {
					return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
				}
				if !notExistingResp {
					return metadata.ResourceRequiresImport(d.ResourceType(), id)
				}
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.DiskPoolName, id.IscsiTargetName, storagepool.IscsiTargetCreate{
				IscsiTargetCreateProperties: &storagepool.IscsiTargetCreateProperties{
					ACLMode:    storagepool.IscsiTargetACLMode(targetModel.ACLMode),
					TargetIqn:  &targetModel.TargetIqn,
					StaticAcls: expandDisksPoolIscsiTargetStaticAcls(targetModel),
					Luns:       expandDisksPoolIscsiTargetLuns(targetModel),
				},
				ManagedByExtended: &targetModel.ManagedByExtended,
				Name:              utils.String(targetModel.Name),
			})

			if err != nil {
				return err
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of DisksPool iscsi target %q : %+v", id.String(), err)
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
			poolID := parse.NewStorageDisksPoolID(id.SubscriptionId, id.ResourceGroup, id.DiskPoolName)
			client := metadata.Client.Storage.DisksPoolIscsiTargetClient
			resp, err := client.Get(ctx, id.ResourceGroup, id.DiskPoolName, id.IscsiTargetName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			model := DiskPoolIscsiTargetModel{
				ACLMode:     string(resp.ACLMode),
				DisksPoolId: poolID.ID(),
				ManagedBy:   utils.NormalizeNilableString(resp.ManagedBy),
				Name:        *resp.Name,
				TargetIqn:   utils.NormalizeNilableString(resp.TargetIqn),
			}
			if resp.ManagedByExtended != nil {
				model.ManagedByExtended = *resp.ManagedByExtended
			}
			model.StaticAcls = flattenDisksPoolIscsiTargetStaticAcls(resp)
			model.Luns = flattenDisksPoolIscsiTargetLuns(resp)
			if resp.Endpoints != nil {
				model.Endpoints = *resp.Endpoints
			}
			if resp.Port != nil {
				model.Port = int(*resp.Port)
			}
			return metadata.Encode(&model)
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
			poolID := parse.NewStorageDisksPoolID(id.SubscriptionId, id.ResourceGroup, id.DiskPoolName)
			locks.ByID(poolID.ID())
			defer locks.UnlockByID(poolID.ID())

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
			poolId := parse.NewStorageDisksPoolID(id.SubscriptionId, id.ResourceGroup, id.DiskPoolName)
			locks.ByID(poolId.ID())
			defer locks.UnlockByID(poolId.ID())
			client := metadata.Client.Storage.DisksPoolIscsiTargetClient
			patch := storagepool.IscsiTargetUpdate{
				IscsiTargetUpdateProperties: &storagepool.IscsiTargetUpdateProperties{},
			}
			m := DiskPoolIscsiTargetModel{}

			err = metadata.Decode(&m)
			if err != nil {
				return err
			}
			if r.HasChange("managed_by_extended") {
				patch.ManagedByExtended = &m.ManagedByExtended
			}
			if r.HasChange("static_acls") {
				patch.StaticAcls = expandDisksPoolIscsiTargetStaticAcls(m)
			}

			if r.HasChange("luns") {
				patch.Luns = expandDisksPoolIscsiTargetLuns(m)
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

func expandDisksPoolIscsiTargetLuns(model DiskPoolIscsiTargetModel) *[]storagepool.IscsiLun {
	var luns []storagepool.IscsiLun
	for _, lun := range model.Luns {
		luns = append(luns, storagepool.IscsiLun{
			Name:                       &lun.Name,
			ManagedDiskAzureResourceID: &lun.ManagedDiskAzureResourceId,
			Lun:                        utils.Int32(int32(lun.Lun)),
		})
	}
	if len(luns) == 0 {
		return nil
	}
	return &luns
}

func flattenDisksPoolIscsiTargetLuns(resp storagepool.IscsiTarget) []DisksPoolIscsiLun {
	var luns []DisksPoolIscsiLun
	if resp.Luns != nil {
		for _, lun := range *resp.Luns {
			l := DisksPoolIscsiLun{}
			if lun.Name != nil {
				l.Name = *lun.Name
			}
			if lun.ManagedDiskAzureResourceID != nil {
				l.ManagedDiskAzureResourceId = *lun.ManagedDiskAzureResourceID
			}
			if lun.Lun != nil {
				l.Lun = int(*lun.Lun)
			}
			luns = append(luns, l)
		}
	}
	return luns
}

func expandDisksPoolIscsiTargetStaticAcls(model DiskPoolIscsiTargetModel) *[]storagepool.ACL {
	var acls []storagepool.ACL
	for _, acl := range model.StaticAcls {
		acls = append(acls, storagepool.ACL{
			InitiatorIqn: &acl.InitiatorIqn,
			MappedLuns:   &acl.MappedLuns,
		})
	}
	return &acls
}

func flattenDisksPoolIscsiTargetStaticAcls(resp storagepool.IscsiTarget) []DisksPoolISCSITargetAcl {
	var acls []DisksPoolISCSITargetAcl
	if resp.StaticAcls != nil {
		for _, acl := range *resp.StaticAcls {
			a := DisksPoolISCSITargetAcl{}
			if acl.InitiatorIqn != nil {
				a.InitiatorIqn = *acl.InitiatorIqn
			}
			if acl.MappedLuns != nil {
				a.MappedLuns = *acl.MappedLuns
			}
			acls = append(acls, a)
		}
	}
	return acls
}

func possibleIscsiTargetACLModeValues() []string {
	var values []string
	for _, v := range storagepool.PossibleIscsiTargetACLModeValues() {
		values = append(values, string(v))
	}
	return values
}
