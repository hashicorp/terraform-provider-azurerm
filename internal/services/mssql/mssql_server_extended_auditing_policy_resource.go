// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
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
	client := meta.(*clients.Client).MSSQL.ServerExtendedBlobAuditingPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for MsSql Server Extended Auditing Policy creation.")

	serverId, err := commonids.ParseSqlServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, serverId.ResourceGroupName, serverId.ServerName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("retrieving MsSql Server Extended Auditing Policy %s: %+v", serverId, err)
			}
		}

		// if state is not disabled, we should import it.
		if existing.ID != nil && *existing.ID != "" && existing.ExtendedServerBlobAuditingPolicyProperties != nil && existing.ExtendedServerBlobAuditingPolicyProperties.State != sql.BlobAuditingPolicyStateDisabled {
			return tf.ImportAsExistsError("azurerm_mssql_server_extended_auditing_policy", *existing.ID)
		}
	}

	params := sql.ExtendedServerBlobAuditingPolicy{
		ExtendedServerBlobAuditingPolicyProperties: &sql.ExtendedServerBlobAuditingPolicyProperties{
			StorageEndpoint:             utils.String(d.Get("storage_endpoint").(string)),
			IsStorageSecondaryKeyInUse:  utils.Bool(d.Get("storage_account_access_key_is_secondary").(bool)),
			RetentionDays:               utils.Int32(int32(d.Get("retention_in_days").(int))),
			IsAzureMonitorTargetEnabled: utils.Bool(d.Get("log_monitoring_enabled").(bool)),
		},
	}

	if d.Get("enabled").(bool) {
		params.ExtendedServerBlobAuditingPolicyProperties.State = sql.BlobAuditingPolicyStateEnabled
	} else {
		params.ExtendedServerBlobAuditingPolicyProperties.State = sql.BlobAuditingPolicyStateDisabled
	}

	if v, ok := d.GetOk("storage_account_subscription_id"); ok {
		u, err := uuid.FromString(v.(string))
		if err != nil {
			return fmt.Errorf("while parsing storage_account_subscrption_id value %q as UUID: %+v", v.(string), err)
		}
		params.ExtendedServerBlobAuditingPolicyProperties.StorageAccountSubscriptionID = &u
	}

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		params.ExtendedServerBlobAuditingPolicyProperties.StorageAccountAccessKey = utils.String(v.(string))
	}

	if v, ok := d.GetOk("predicate_expression"); ok {
		params.ExtendedServerBlobAuditingPolicyProperties.PredicateExpression = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("audit_actions_and_groups"); ok && len(v.([]interface{})) > 0 {
		params.ExtendedServerBlobAuditingPolicyProperties.AuditActionsAndGroups = utils.ExpandStringSlice(v.([]interface{}))
	}

	future, err := client.CreateOrUpdate(ctx, serverId.ResourceGroupName, serverId.ServerName, params)
	if err != nil {
		return fmt.Errorf("creating MsSql Server Extended Auditing Policy %s: %+v", serverId, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of MsSql Server Extended Auditing Policy %s: %+v", serverId, err)
	}

	id := parse.NewServerExtendedAuditingPolicyID(subscriptionId, serverId.ResourceGroupName, serverId.ServerName, "default")

	d.SetId(id.ID())

	return resourceMsSqlServerExtendedAuditingPolicyRead(d, meta)
}

func resourceMsSqlServerExtendedAuditingPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerExtendedBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerExtendedAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading MsSql Server Extended Auditing Policy %s: %+v", id, err)
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)
	d.Set("server_id", serverId.ID())

	if props := resp.ExtendedServerBlobAuditingPolicyProperties; props != nil {
		d.Set("storage_endpoint", props.StorageEndpoint)
		d.Set("storage_account_access_key_is_secondary", props.IsStorageSecondaryKeyInUse)
		d.Set("retention_in_days", props.RetentionDays)
		d.Set("log_monitoring_enabled", props.IsAzureMonitorTargetEnabled)
		d.Set("enabled", props.State == sql.BlobAuditingPolicyStateEnabled)
		d.Set("predicate_expression", props.PredicateExpression)
		d.Set("audit_actions_and_groups", utils.FlattenStringSlice(props.AuditActionsAndGroups))

		if props.StorageAccountSubscriptionID.String() != "00000000-0000-0000-0000-000000000000" {
			d.Set("storage_account_subscription_id", props.StorageAccountSubscriptionID.String())
		}
	}

	return nil
}

func resourceMsSqlServerExtendedAuditingPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerExtendedBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerExtendedAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	params := sql.ExtendedServerBlobAuditingPolicy{
		ExtendedServerBlobAuditingPolicyProperties: &sql.ExtendedServerBlobAuditingPolicyProperties{
			State: sql.BlobAuditingPolicyStateDisabled,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, params)
	if err != nil {
		return fmt.Errorf("deleting MsSql Server %q Extended Auditing Policy(Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of MsSql Server %q Extended Auditing Policy (Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	return nil
}
