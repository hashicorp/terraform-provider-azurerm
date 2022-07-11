package disks

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/sdk/2021-08-01/iscsitargets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DisksPoolIscsiTargetResource struct{}

var _ sdk.Resource = DisksPoolIscsiTargetResource{}

type DiskPoolIscsiTargetModel struct {
	ACLMode     string   `tfschema:"acl_mode"`
	DisksPoolId string   `tfschema:"disks_pool_id"`
	Endpoints   []string `tfschema:"endpoints"`
	Name        string   `tfschema:"name"`
	Port        int      `tfschema:"port"`
	TargetIqn   string   `tfschema:"target_iqn"`
}

func (d DisksPoolIscsiTargetResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(5, 223),
				validation.StringMatch(
					regexp.MustCompile(`[a-z\d.\-]*[a-z\d]$`),
					"The iSCSI target name can only contain lowercase letters, numbers, periods, or hyphens.",
				),
			),
		},

		"acl_mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice(
				[]string{string(iscsitargets.IscsiTargetAclModeDynamic)},
				false,
			),
		},

		"disks_pool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iscsitargets.ValidateDiskPoolID,
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
	return "azurerm_disk_pool_iscsi_target"
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

			poolId, err := iscsitargets.ParseDiskPoolID(m.DisksPoolId)
			if err != nil {
				return err
			}
			if poolId.SubscriptionId != metadata.Client.Account.SubscriptionId {
				return fmt.Errorf("Disk Pool subscription id %q is different from provider's subscription", poolId.SubscriptionId)
			}

			id := iscsitargets.NewIscsiTargetID(poolId.SubscriptionId, poolId.ResourceGroupName, poolId.DiskPoolName, m.Name)
			client := metadata.Client.Disks.DisksPoolIscsiTargetClient
			locks.ByID(poolId.ID())
			defer locks.UnlockByID(poolId.ID())

			existing, err := client.Get(ctx, id)
			notExistingResp := response.WasNotFound(existing.HttpResponse)
			if err != nil && !notExistingResp {
				return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
			}
			if !notExistingResp {
				return metadata.ResourceRequiresImport(d.ResourceType(), id)
			}

			future, err := client.CreateOrUpdate(ctx, id, iscsitargets.IscsiTargetCreate{
				Properties: iscsitargets.IscsiTargetCreateProperties{
					AclMode:   iscsitargets.IscsiTargetAclMode(m.ACLMode),
					TargetIqn: &m.TargetIqn,
				},
			})
			if err != nil {
				return fmt.Errorf("creating DisksPool iscsi target %q : %+v", id.ID(), err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id.ID())
			}
			//lintignore:R006
			return pluginsdk.Retry(time.Until(deadline), func() *resource.RetryError {
				//lintignore:R006
				if err := d.retryError("waiting for creation DisksPool iscsi target", id.ID(), future.Poller.PollUntilDone()); err != nil {
					return err
				}
				metadata.SetID(id)
				return nil
			})
		},
	}
}

func (d DisksPoolIscsiTargetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := iscsitargets.ParseIscsiTargetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			client := metadata.Client.Disks.DisksPoolIscsiTargetClient
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("malformed Disk Pool response %q : %+v", id.ID(), resp)
			}
			poolID := iscsitargets.NewDiskPoolID(id.SubscriptionId, id.ResourceGroupName, id.DiskPoolName)
			m := DiskPoolIscsiTargetModel{
				ACLMode:     string(resp.Model.Properties.AclMode),
				DisksPoolId: poolID.ID(),
				Name:        id.IscsiTargetName,
			}
			if endpoints := resp.Model.Properties.Endpoints; endpoints != nil {
				m.Endpoints = *endpoints
			}
			if port := resp.Model.Properties.Port; port != nil {
				m.Port = int(*port)
			}
			m.TargetIqn = resp.Model.Properties.TargetIqn
			return metadata.Encode(&m)
		},
	}
}

func (d DisksPoolIscsiTargetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := iscsitargets.ParseIscsiTargetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			client := metadata.Client.Disks.DisksPoolIscsiTargetClient
			future, err := client.Delete(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting DisksPool iscsi target %q: %+v", id.ID(), err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id)
			}
			//lintignore:R006
			return pluginsdk.Retry(time.Until(deadline), func() *resource.RetryError {
				return d.retryError("waiting for deletion of DisksPool iscsi target", id.ID(), future.Poller.PollUntilDone())
			})
		},
	}
}

func (d DisksPoolIscsiTargetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return iscsitargets.ValidateIscsiTargetID
}

func (DisksPoolIscsiTargetResource) retryError(action string, id string, err error) *resource.RetryError {
	if err == nil {
		return nil
	}

	// according to https://docs.microsoft.com/en-us/azure/virtual-machines/disks-pools-troubleshoot#common-failure-codes-when-enabling-iscsi-on-disk-pools the errors below are retryable.
	retryableErrors := []string{
		"GoalStateApplicationTimeoutError",
		"OngoingOperationInProgress",
	}
	for _, retryableError := range retryableErrors {
		if strings.Contains(err.Error(), retryableError) {
			return pluginsdk.RetryableError(fmt.Errorf("%s %s: %+v", action, id, err))
		}
	}
	return pluginsdk.NonRetryableError(fmt.Errorf("%s %s: %+v", action, id, err))
}
