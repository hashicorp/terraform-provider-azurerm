package storage

import (
	"github.com/Azure/azure-sdk-for-go/services/storagepool/mgmt/2021-08-01/storagepool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"regexp"
)

type DisksPoolISCSITargetResource struct{}

var _ sdk.ResourceWithUpdate = DisksPoolISCSITargetResource{}

type DiskPoolISCSITargetModel struct {
	AclMode           string                    `tfschema:"acl_mode"`
	ManagedByExtended []string                  `tfschema:"managed_by_extended"`
	Name              string                    `tfschema:"name"`
	TargetIQN         string                    `tfschema:"target_iqn"` //iSCSI Target IQN (iSCSI Qualified Name); example: \"iqn.2005-03.org.iscsi:server\".
	StaticAcls        []DisksPoolISCSITargetAcl `tfschema:"static_acls"`
	Luns              []DisksPoolISCSIILun      `tfschema:"luns"`
	Endpoints         []string                  `tfschema:"endpoints"` //List of private IPv4 addresses to connect to the iSCSI Target.
	Port              int                       `tfschema:"port"`      //The port used by iSCSI Target portal group.

	// TODO: Should we add this?
	// Session []string //List of identifiers for active sessions on the iSCSI target
}

type DisksPoolISCSITargetAcl struct {
	InitiatorIQN string   `tfschema:"initiator_iqn"` //iSCSI initiator IQN (iSCSI Qualified Name); example: \"iqn.2005-03.org.iscsi:client\".
	MappedLuns   []string `tfschema:"mapped_luns"`   //List of LUN names mapped to the ACL.
}

type DisksPoolISCSIILun struct {
	Name                       string `tfschema:"name"`                           //User defined name for iSCSI LUN; example: \"lun0\" min:1 max:90
	ManagedDiskAzureResourceId string `tfschema:"managed_disk_azure_resource_id"` //Azure Resource ID of the Managed Disk.
	Lun                        int    `tfschema:"lun"`                            //Specifies the Logical Unit Number of the iSCSI LUN. readonly
}

//{"properties":{"aclMode":"Dynamic","targetIqn":"iqn.2021-11.com.microsoft:teststesttestets","luns":[]}}
func (d DisksPoolISCSITargetResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"acl_mode": &schema.Schema{
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(possibleIscsiTargetACLModeValues(), false),
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
			// TODO: How to validate?
			ValidateFunc: validation.All(
				validation.StringLenBetween(5, 223),
				validation.StringMatch(
					regexp.MustCompile("[a-z\\d.\\-]*[a-z\\d]$"),
					"The iSCSI target name can only contain lowercase letters, numbers, periods, or hyphens.",
				),
			),
		},
		"target_iqn": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// TODO: How to validate?
			ValidateFunc: validate.ValidateIQN,
		},
		"static_acls": { // TODO: How to set static acl?
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"initiator_iqn": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validate.ValidateIQN,
					},
					"mapped_luns": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &schema.Schema{
							Type: pluginsdk.TypeString,
							// TODO: validate
						},
					},
				},
			},
		},
		"luns": { // Allow empty slice!
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						// TODO: validate
					},
					"managed_disk_azure_resource_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: computeValidate.ManagedDiskID,
					},
					"lun": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
}

func (d DisksPoolISCSITargetResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"managed_by": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"endpoints": {
			Type:     pluginsdk.TypeList,
			Computed: true,
		},
		"port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
	}
}

func (d DisksPoolISCSITargetResource) ModelObject() interface{} {
	return &DisksPoolISCSITargetResource{}
}

func (d DisksPoolISCSITargetResource) ResourceType() string {
	return "azurerm_storage_disks_pool_iscsi_target"
}

func (d DisksPoolISCSITargetResource) Create() sdk.ResourceFunc {
	panic("implement me")
}

func (d DisksPoolISCSITargetResource) Read() sdk.ResourceFunc {
	panic("implement me")
}

func (d DisksPoolISCSITargetResource) Delete() sdk.ResourceFunc {
	panic("implement me")
}

func (d DisksPoolISCSITargetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	panic("implement me")
}

func (d DisksPoolISCSITargetResource) Update() sdk.ResourceFunc {
	panic("implement me")
}

func possibleIscsiTargetACLModeValues() []string {
	var values []string
	for _, v := range storagepool.PossibleIscsiTargetACLModeValues() {
		values = append(values, string(v))
	}
	return values
}
