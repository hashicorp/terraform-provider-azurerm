// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/customermanagedkeys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	managedHsmValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsClusterCustomerManagedKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsClusterCustomerManagedKeyCreate,
		Read:   resourceLogAnalyticsClusterCustomerManagedKeyRead,
		Update: resourceLogAnalyticsClusterCustomerManagedKeyUpdate,
		Delete: resourceLogAnalyticsClusterCustomerManagedKeyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(6 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(6 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusters.ParseClusterID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ClusterCustomerManagedKeyV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"log_analytics_cluster_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: clusters.ValidateClusterID,
			},

			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
				ExactlyOneOf: []string{"key_vault_key_id", "managed_hsm_key_id"},
			},

			"managed_hsm_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.Any(managedHsmValidate.ManagedHSMDataPlaneVersionedKeyID, managedHsmValidate.ManagedHSMDataPlaneVersionlessKeyID),
				ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_key_id"},
			},
		},
	}
}

func resourceLogAnalyticsClusterCustomerManagedKeyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	accountClient := meta.(*clients.Client).Account
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Get("log_analytics_cluster_id").(string))
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `model` is nil", *id)
	}

	props := model.Properties
	if props == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `Properties` is nil", *id)
	}

	if props.KeyVaultProperties != nil {
		if keyProps := *props.KeyVaultProperties; keyProps.KeyName != nil && *keyProps.KeyName != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_cluster_customer_managed_key", id.ID())
		}
	}

	// Ensure `associatedWorkspaces` is not present in request, this is a read only property and cannot be sent to the API
	// Error: updating Customer Managed Key for Cluster
	//		performing CreateOrUpdate: unexpected status 400 (400 Bad Request) with error:
	//		InvalidParameter: 'properties.associatedWorkspaces' is a read only property and cannot be set.
	//		Please refer to https://docs.microsoft.com/en-us/azure/azure-monitor/log-query/logs-dedicated-clusters#link-a-workspace-to-the-cluster for more information on how to associate a workspace to the cluster.
	props.AssociatedWorkspaces = nil

	if cmkID, err := customermanagedkeys.ExpandKeyVaultOrManagedHSMKey(d, customermanagedkeys.VersionTypeAny, accountClient.Environment.KeyVault, accountClient.Environment.ManagedHSM); err != nil {
		return fmt.Errorf("expanding Customer Managed Key: %+v", err)
	} else if cmkID.IsSet() {
		model.Properties.KeyVaultProperties = &clusters.KeyVaultProperties{
			KeyVaultUri: pointer.To(cmkID.BaseUri()),
			KeyName:     pointer.To(cmkID.Name()),
			KeyVersion:  pointer.To(cmkID.Version()),
		}
	} else {
		return fmt.Errorf("expanding Customer Managed Key: no ID returned")
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
		return fmt.Errorf("creating Customer Managed Key for %s: %+v", *id, err)
	}

	updateWait, err := logAnalyticsClusterWaitForState(ctx, client, *id)
	if err != nil {
		return err
	}
	if _, err := updateWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish adding Customer Managed Key: %+v", *id, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsClusterCustomerManagedKeyRead(d, meta)
}

func resourceLogAnalyticsClusterCustomerManagedKeyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	accountClient := meta.(*clients.Client).Account
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `model` is nil", *id)
	}

	if props := model.Properties; props == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `Properties` is nil", *id)
	}

	// This is a read only property, please see comment in the create function.
	model.Properties.AssociatedWorkspaces = nil

	if cmkID, err := customermanagedkeys.ExpandKeyVaultOrManagedHSMKey(d, customermanagedkeys.VersionTypeAny, accountClient.Environment.KeyVault, accountClient.Environment.ManagedHSM); err != nil {
		return fmt.Errorf("expanding Customer Managed Key: %+v", err)
	} else if cmkID.IsSet() {
		model.Properties.KeyVaultProperties = &clusters.KeyVaultProperties{
			KeyVaultUri: pointer.To(cmkID.BaseUri()),
			KeyName:     pointer.To(cmkID.Name()),
			KeyVersion:  pointer.To(cmkID.Version()),
		}
	} else {
		return fmt.Errorf("expanding Customer Managed Key: no ID returned")
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
		return fmt.Errorf("updating Customer Managed Key for %s: %+v", *id, err)
	}

	return resourceLogAnalyticsClusterCustomerManagedKeyRead(d, meta)
}

func resourceLogAnalyticsClusterCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	accountClient := meta.(*clients.Client).Account
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if kvProps := props.KeyVaultProperties; kvProps != nil {
				cmkUri := pointer.From(kvProps.KeyVaultUri)
				cmkName := pointer.From(kvProps.KeyName)
				cmkVersion := pointer.From(kvProps.KeyVersion)

				if cmkUri != "" && cmkName != "" {
					keyId := fmt.Sprintf("%s/keys/%s", strings.TrimSuffix(cmkUri, "/"), cmkName)
					if cmkVersion != "" {
						keyId = fmt.Sprintf("%s/%s", keyId, cmkVersion)
					}

					if cmkID, err := customermanagedkeys.FlattenKeyVaultOrManagedHSMID(keyId, accountClient.Environment.ManagedHSM); err != nil {
						return fmt.Errorf("flattening %s: %+v", keyId, err)
					} else if cmkID.IsSet() {
						if cmkID.KeyVaultKeyId != nil {
							d.Set("key_vault_key_id", cmkID.KeyVaultKeyID())
						} else {
							d.Set("managed_hsm_key_id", cmkID.ManagedHSMKeyID())
						}
					} else {
						log.Printf("[DEBUG] %s has no Customer Managed Key - removing from state", *id)
						d.SetId("")
						return nil
					}
				}
			}
		}
	}

	d.Set("log_analytics_cluster_id", d.Id())

	return nil
}

func resourceLogAnalyticsClusterCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `model` is nil", *id)
	}

	props := model.Properties
	if props == nil {
		return fmt.Errorf("retrieving `azurerm_log_analytics_cluster` %s: `Properties` is nil", *id)
	}

	if props.KeyVaultProperties == nil {
		return fmt.Errorf("deleting `azurerm_log_analytics_cluster_customer_managed_key` %s: `customer managed key does not exist!`", *id)
	}

	if props.KeyVaultProperties.KeyName == nil || *props.KeyVaultProperties.KeyName == "" {
		return fmt.Errorf("deleting `azurerm_log_analytics_cluster_customer_managed_key` %s: `customer managed key does not exist!`", *id)
	}

	// This is a read only property, please see comment in the create function.
	props.AssociatedWorkspaces = nil

	// The API only removes the CMK when it is sent empty string values, sending nil for each property or an empty object does not work.
	model.Properties.KeyVaultProperties = &clusters.KeyVaultProperties{
		KeyVaultUri: pointer.To(""),
		KeyName:     pointer.To(""),
		KeyVersion:  pointer.To(""),
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
		return fmt.Errorf("deleting Customer Managed Key from %s: %+v", *id, err)
	}

	return nil
}
