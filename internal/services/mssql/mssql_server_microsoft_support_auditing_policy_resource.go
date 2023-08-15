// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/gofrs/uuid"
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

	serverId, err := parse.ServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name, "default")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Failed to check for presence of existing Server %q Sql Microsoft Support Auditing (Resource Group %q): %s", serverId.Name, serverId.ResourceGroup, err)
			}
		}

		// if state is not disabled, we should import it.
		if existing.ID != nil && *existing.ID != "" && existing.ServerDevOpsAuditSettingsProperties != nil && existing.ServerDevOpsAuditSettingsProperties.State != sql.BlobAuditingPolicyStateDisabled {
			return tf.ImportAsExistsError("azurerm_mssql_server_microsoft_support_auditing_policy", *existing.ID)
		}
	}

	params := sql.ServerDevOpsAuditingSettings{
		ServerDevOpsAuditSettingsProperties: &sql.ServerDevOpsAuditSettingsProperties{
			StorageEndpoint:             utils.String(d.Get("blob_storage_endpoint").(string)),
			IsAzureMonitorTargetEnabled: utils.Bool(d.Get("log_monitoring_enabled").(bool)),
		},
	}

	if d.Get("enabled").(bool) {
		params.ServerDevOpsAuditSettingsProperties.State = sql.BlobAuditingPolicyStateEnabled
	} else {
		params.ServerDevOpsAuditSettingsProperties.State = sql.BlobAuditingPolicyStateDisabled
	}

	if v, ok := d.GetOk("storage_account_subscription_id"); ok {
		u, err := uuid.FromString(v.(string))
		if err != nil {
			return fmt.Errorf("while parsing storage_account_subscrption_id value %q as UUID: %+v", v.(string), err)
		}
		params.ServerDevOpsAuditSettingsProperties.StorageAccountSubscriptionID = &u
	}

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		params.ServerDevOpsAuditSettingsProperties.StorageAccountAccessKey = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, serverId.ResourceGroup, serverId.Name, "default", params)
	if err != nil {
		return fmt.Errorf("creating MsSql Server %q Microsoft Support Auditing Policy (Resource Group %q): %+v", serverId.Name, serverId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of MsSql Server %q Microsoft Support Auditing Policy (Resource Group %q): %+v", serverId.Name, serverId.ResourceGroup, err)
	}

	read, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name, "default")
	if err != nil {
		return fmt.Errorf("retrieving MsSql Server %q Microsoft Support Auditing Policy (Resource Group %q): %+v", serverId.Name, serverId.ResourceGroup, err)
	}

	if read.Name == nil || *read.Name == "" {
		return fmt.Errorf("reading MsSql Server %q Microsoft Support Auditing Policy (Resource Group %q) Name is empty or nil", serverId.Name, serverId.ResourceGroup)
	}
	id := parse.NewServerMicrosoftSupportAuditingPolicyID(subscriptionId, serverId.ResourceGroup, serverId.Name, *read.Name)

	d.SetId(id.ID())

	return resourceMsSqlServerMicrosoftSupportAuditingPolicyRead(d, meta)
}

func resourceMsSqlServerMicrosoftSupportAuditingPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServerDevOpsAuditSettingsClient
	serverClient := meta.(*clients.Client).MSSQL.ServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerMicrosoftSupportAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, "default")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading MsSql Server %s Microsoft Support Auditing Policy (Resource Group %q): %s", id.ServerName, id.ResourceGroup, err)
	}

	serverResp, err := serverClient.Get(ctx, id.ResourceGroup, id.ServerName, "")
	if err != nil || serverResp.ID == nil || *serverResp.ID == "" {
		return fmt.Errorf("reading MsSql Server %q ID is empty or nil(Resource Group %q): %s", id.ServerName, id.ResourceGroup, err)
	}

	d.Set("server_id", serverResp.ID)

	if props := resp.ServerDevOpsAuditSettingsProperties; props != nil {
		d.Set("blob_storage_endpoint", props.StorageEndpoint)
		d.Set("log_monitoring_enabled", props.IsAzureMonitorTargetEnabled)
		d.Set("enabled", props.State == sql.BlobAuditingPolicyStateEnabled)

		if props.StorageAccountSubscriptionID.String() != "00000000-0000-0000-0000-000000000000" {
			d.Set("storage_account_subscription_id", props.StorageAccountSubscriptionID.String())
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

	params := sql.ServerDevOpsAuditingSettings{
		ServerDevOpsAuditSettingsProperties: &sql.ServerDevOpsAuditSettingsProperties{
			State: sql.BlobAuditingPolicyStateDisabled,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, "default", params)
	if err != nil {
		return fmt.Errorf("deleting MsSql Server %q Microsoft Support Auditing Policy(Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of MsSql Server %q Microsoft Support Auditing Policy (Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	return nil
}
