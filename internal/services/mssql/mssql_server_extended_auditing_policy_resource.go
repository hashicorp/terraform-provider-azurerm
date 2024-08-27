// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/blobauditing"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlServerExtendedAuditingPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlServerExtendedAuditingPolicyCreateUpdate,
		Read:   resourceMsSqlServerExtendedAuditingPolicyRead,
		Update: resourceMsSqlServerExtendedAuditingPolicyCreateUpdate,
		Delete: resourceMsSqlServerExtendedAuditingPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ServerExtendedAuditingPolicyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"storage_endpoint": {
				// TODO 4.0: rename to `blob_storage_endpoint`
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPS,
			},

			"storage_account_access_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_account_access_key_is_secondary": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"retention_in_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 3285),
			},

			"log_monitoring_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"storage_account_subscription_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.IsUUID,
			},

			"predicate_expression": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"audit_actions_and_groups": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				// audit_actions_and_groups seems to be pre-populated with values ["SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP", "FAILED_DATABASE_AUTHENTICATION_GROUP", "BATCH_COMPLETED_GROUP"],
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceMsSqlServerExtendedAuditingPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.BlobAuditingPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for MsSql Server Extended Auditing Policy creation.")

	serverId, err := commonids.ParseSqlServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.ExtendedServerBlobAuditingPoliciesGet(ctx, *serverId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("retrieving MsSql Server Extended Auditing Policy %s: %+v", serverId, err)
			}
		}

		// if state is not disabled, we should import it.
		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" && existing.Model.Properties != nil && existing.Model.Properties.State != blobauditing.BlobAuditingPolicyStateDisabled {
			return tf.ImportAsExistsError("azurerm_mssql_server_extended_auditing_policy", *existing.Model.Id)
		}
	}

	params := blobauditing.ExtendedServerBlobAuditingPolicy{
		Properties: &blobauditing.ExtendedServerBlobAuditingPolicyProperties{
			StorageEndpoint:             utils.String(d.Get("storage_endpoint").(string)),
			IsStorageSecondaryKeyInUse:  utils.Bool(d.Get("storage_account_access_key_is_secondary").(bool)),
			RetentionDays:               utils.Int64(int64(d.Get("retention_in_days").(int))),
			IsAzureMonitorTargetEnabled: utils.Bool(d.Get("log_monitoring_enabled").(bool)),
		},
	}

	if d.Get("enabled").(bool) {
		params.Properties.State = blobauditing.BlobAuditingPolicyStateEnabled
	} else {
		params.Properties.State = blobauditing.BlobAuditingPolicyStateDisabled
	}

	if v, ok := d.GetOk("storage_account_subscription_id"); ok {
		params.Properties.StorageAccountSubscriptionId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		params.Properties.StorageAccountAccessKey = utils.String(v.(string))
	}

	if v, ok := d.GetOk("predicate_expression"); ok {
		params.Properties.PredicateExpression = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("audit_actions_and_groups"); ok && len(v.([]interface{})) > 0 {
		params.Properties.AuditActionsAndGroups = utils.ExpandStringSlice(v.([]interface{}))
	}

	err = client.ExtendedServerBlobAuditingPoliciesCreateOrUpdateThenPoll(ctx, *serverId, params)
	if err != nil {
		return fmt.Errorf("creating MsSql Server Extended Auditing Policy %s: %+v", serverId, err)
	}

	id := parse.NewServerExtendedAuditingPolicyID(subscriptionId, serverId.ResourceGroupName, serverId.ServerName, "default")

	d.SetId(id.ID())

	return resourceMsSqlServerExtendedAuditingPolicyRead(d, meta)
}

func resourceMsSqlServerExtendedAuditingPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.BlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerExtendedAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	resp, err := client.ExtendedServerBlobAuditingPoliciesGet(ctx, serverId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading MsSql Server Extended Auditing Policy %s: %+v", id, err)
	}

	d.Set("server_id", serverId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("storage_endpoint", props.StorageEndpoint)
			d.Set("storage_account_access_key_is_secondary", props.IsStorageSecondaryKeyInUse)
			d.Set("retention_in_days", props.RetentionDays)
			d.Set("log_monitoring_enabled", props.IsAzureMonitorTargetEnabled)
			d.Set("enabled", props.State == blobauditing.BlobAuditingPolicyStateEnabled)
			d.Set("predicate_expression", props.PredicateExpression)
			d.Set("audit_actions_and_groups", utils.FlattenStringSlice(props.AuditActionsAndGroups))

			if pointer.From(props.StorageAccountSubscriptionId) != "00000000-0000-0000-0000-000000000000" {
				d.Set("storage_account_subscription_id", props.StorageAccountSubscriptionId)
			}
		}
	}

	return nil
}

func resourceMsSqlServerExtendedAuditingPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.BlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerExtendedAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	params := blobauditing.ExtendedServerBlobAuditingPolicy{
		Properties: &blobauditing.ExtendedServerBlobAuditingPolicyProperties{
			State: blobauditing.BlobAuditingPolicyStateDisabled,
		},
	}

	err = client.ExtendedServerBlobAuditingPoliciesCreateOrUpdateThenPoll(ctx, serverId, params)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", serverId, err)
	}
	return nil
}
