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
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverdevopsaudit"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlServerMicrosoftSupportAuditingPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlServerMicrosoftSupportAuditingPolicyCreateUpdate,
		Read:   resourceMsSqlServerMicrosoftSupportAuditingPolicyRead,
		Update: resourceMsSqlServerMicrosoftSupportAuditingPolicyCreateUpdate,
		Delete: resourceMsSqlServerMicrosoftSupportAuditingPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ServerMicrosoftSupportAuditingPolicyID(id)
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

			"blob_storage_endpoint": {
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
		},
	}
}

func resourceMsSqlServerMicrosoftSupportAuditingPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerDevOpsAuditSettingsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for MsSql Server Microsoft Support Auditing Policy creation.")

	serverId, err := commonids.ParseSqlServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.SettingsGet(ctx, *serverId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("retrieving MsSql Server Microsoft Support Auditing Policy %s: %+v", serverId, err)
			}
		}

		// if state is not disabled, we should import it.
		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" && existing.Model.Properties != nil && existing.Model.Properties.State != serverdevopsaudit.BlobAuditingPolicyStateDisabled {
			return tf.ImportAsExistsError("azurerm_mssql_server_microsoft_support_auditing_policy", *existing.Model.Id)
		}
	}

	params := serverdevopsaudit.ServerDevOpsAuditingSettings{
		Properties: &serverdevopsaudit.ServerDevOpsAuditSettingsProperties{
			IsAzureMonitorTargetEnabled: utils.Bool(d.Get("log_monitoring_enabled").(bool)),
		},
	}

	if v := d.Get("blob_storage_endpoint").(string); v != "" {
		params.Properties.StorageEndpoint = utils.String(v)
	}

	if d.Get("enabled").(bool) {
		params.Properties.State = serverdevopsaudit.BlobAuditingPolicyStateEnabled
	} else {
		params.Properties.State = serverdevopsaudit.BlobAuditingPolicyStateDisabled
	}

	if v, ok := d.GetOk("storage_account_subscription_id"); ok {
		params.Properties.StorageAccountSubscriptionId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		params.Properties.StorageAccountAccessKey = utils.String(v.(string))
	}

	err = client.SettingsCreateOrUpdateThenPoll(ctx, *serverId, params)
	if err != nil {
		return fmt.Errorf("creating MsSql Server Microsoft Support Auditing Policy %s: %+v", serverId, err)
	}

	id := parse.NewServerMicrosoftSupportAuditingPolicyID(subscriptionId, serverId.ResourceGroupName, serverId.ServerName, "default")

	d.SetId(id.ID())

	return resourceMsSqlServerMicrosoftSupportAuditingPolicyRead(d, meta)
}

func resourceMsSqlServerMicrosoftSupportAuditingPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerDevOpsAuditSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerMicrosoftSupportAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	resp, err := client.SettingsGet(ctx, serverId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %s", *id, err)
	}

	d.Set("server_id", serverId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("blob_storage_endpoint", props.StorageEndpoint)
			d.Set("log_monitoring_enabled", props.IsAzureMonitorTargetEnabled)
			d.Set("enabled", props.State == serverdevopsaudit.BlobAuditingPolicyStateEnabled)

			if pointer.From(props.StorageAccountSubscriptionId) != "00000000-0000-0000-0000-000000000000" {
				d.Set("storage_account_subscription_id", pointer.From(props.StorageAccountSubscriptionId))
			}
		}
	}

	return nil
}

func resourceMsSqlServerMicrosoftSupportAuditingPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerDevOpsAuditSettingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerMicrosoftSupportAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	params := serverdevopsaudit.ServerDevOpsAuditingSettings{
		Properties: &serverdevopsaudit.ServerDevOpsAuditSettingsProperties{
			State: serverdevopsaudit.BlobAuditingPolicyStateDisabled,
		},
	}

	err = client.SettingsCreateOrUpdateThenPoll(ctx, serverId, params)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
