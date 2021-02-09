package loganalytics

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceLogAnalyticsClusterCustomerManagedKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogAnalyticsClusterCustomerManagedKeyCreate,
		Read:   resourceLogAnalyticsClusterCustomerManagedKeyRead,
		Update: resourceLogAnalyticsClusterCustomerManagedKeyUpdate,
		Delete: resourceLogAnalyticsClusterCustomerManagedKeyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Hour),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Hour),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"log_analytics_cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsClusterID,
			},

			"key_vault_key_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
			},
		},
	}
}

func resourceLogAnalyticsClusterCustomerManagedKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterIdRaw := d.Get("log_analytics_cluster_id").(string)
	clusterId, err := parse.LogAnalyticsClusterID(clusterIdRaw)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, clusterId.ResourceGroup, clusterId.ClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Log Analytics Cluster %q (resource group %q) was not found", clusterId.ClusterName, clusterId.ResourceGroup)
		}
		return fmt.Errorf("failed to get details of Log Analytics Cluster %q (resource group %q): %+v", clusterId.ClusterName, clusterId.ResourceGroup, err)
	}
	if resp.ClusterProperties != nil && resp.ClusterProperties.KeyVaultProperties != nil {
		keyProps := *resp.ClusterProperties.KeyVaultProperties
		if keyProps.KeyName != nil && *keyProps.KeyName != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_cluster_customer_managed_key", fmt.Sprintf("%s/CMK", clusterIdRaw))
		}
	}

	d.SetId(fmt.Sprintf("%s/CMK", clusterIdRaw))
	return resourceLogAnalyticsClusterCustomerManagedKeyUpdate(d, meta)
}

func resourceLogAnalyticsClusterCustomerManagedKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(d.Get("key_vault_key_id").(string))
	if err != nil {
		return fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
	}

	clusterId, err := parse.LogAnalyticsClusterID(d.Get("log_analytics_cluster_id").(string))
	if err != nil {
		return err
	}

	clusterPatch := operationalinsights.ClusterPatch{
		ClusterPatchProperties: &operationalinsights.ClusterPatchProperties{
			KeyVaultProperties: &operationalinsights.KeyVaultProperties{
				KeyVaultURI: utils.String(keyId.KeyVaultBaseUrl),
				KeyName:     utils.String(keyId.Name),
				KeyVersion:  utils.String(keyId.Version),
			},
		},
	}

	if _, err := client.Update(ctx, clusterId.ResourceGroup, clusterId.ClusterName, clusterPatch); err != nil {
		return fmt.Errorf("updating Log Analytics Cluster %q (Resource Group %q): %+v", clusterId.ClusterName, clusterId.ResourceGroup, err)
	}

	updateWait := logAnalyticsClusterWaitForState(ctx, meta, d.Timeout(schema.TimeoutUpdate), clusterId.ResourceGroup, clusterId.ClusterName)

	if _, err := updateWait.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Log Analytics Cluster to finish updating %q (Resource Group %q): %v", clusterId.ClusterName, clusterId.ResourceGroup, err)
	}

	return resourceLogAnalyticsClusterCustomerManagedKeyRead(d, meta)
}

func resourceLogAnalyticsClusterCustomerManagedKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	idRaw := strings.TrimRight(d.Id(), "/CMK")

	id, err := parse.LogAnalyticsClusterID(idRaw)
	if err != nil {
		return err
	}

	d.Set("log_analytics_cluster_id", idRaw)

	resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Log Analytics %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Log Analytics Cluster %q (Resource Group %q): %+v", id.ClusterName, id.ResourceGroup, err)
	}

	if props := resp.ClusterProperties; props != nil {
		if kvProps := props.KeyVaultProperties; kvProps != nil {
			var keyVaultUri, keyName, keyVersion string
			if kvProps.KeyVaultURI != nil && *kvProps.KeyVaultURI != "" {
				keyVaultUri = *kvProps.KeyVaultURI
			} else {
				return fmt.Errorf("empty value returned for Key Vault URI")
			}
			if kvProps.KeyName != nil && *kvProps.KeyName != "" {
				keyName = *kvProps.KeyName
			} else {
				return fmt.Errorf("empty value returned for Key Vault Key Name")
			}
			if kvProps.KeyVersion != nil {
				keyVersion = *kvProps.KeyVersion
			}
			keyVaultKeyId, err := keyVaultParse.NewNestedItemID(keyVaultUri, "keys", keyName, keyVersion)
			if err != nil {
				return err
			}
			d.Set("key_vault_key_id", keyVaultKeyId.ID())
		}
	}

	return nil
}

func resourceLogAnalyticsClusterCustomerManagedKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterId, err := parse.LogAnalyticsClusterID(d.Get("log_analytics_cluster_id").(string))
	if err != nil {
		return err
	}

	clusterPatch := operationalinsights.ClusterPatch{
		ClusterPatchProperties: &operationalinsights.ClusterPatchProperties{
			KeyVaultProperties: &operationalinsights.KeyVaultProperties{
				KeyVaultURI: nil,
				KeyName:     nil,
				KeyVersion:  nil,
			},
		},
	}

	if _, err = client.Update(ctx, clusterId.ResourceGroup, clusterId.ClusterName, clusterPatch); err != nil {
		return fmt.Errorf("removing Log Analytics Cluster Customer Managed Key from cluster %q (resource group %q)", clusterId.ClusterName, clusterId.ResourceGroup)
	}

	deleteWait := logAnalyticsClusterWaitForState(ctx, meta, d.Timeout(schema.TimeoutDelete), clusterId.ResourceGroup, clusterId.ClusterName)

	if _, err := deleteWait.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Log Analytics Cluster to finish updating %q (Resource Group %q): %v", clusterId.ClusterName, clusterId.ResourceGroup, err)
	}

	return nil
}
