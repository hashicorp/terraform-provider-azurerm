package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2020-03-01-preview/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogAnalyticsCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogAnalyticsClusterCreate,
		Read:   resourceArmLogAnalyticsClusterRead,
		Update: resourceArmLogAnalyticsClusterUpdate,
		Delete: resourceArmLogAnalyticsClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.LogAnalyticsClusterID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsClustersName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(operationalinsights.SystemAssigned),
								string(operationalinsights.None),
							}, false),
						},

						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"next_link": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"key_vault_property": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"key_vault_uri": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"key_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"sku": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(operationalinsights.CapacityReservation),
							}, false),
						},

						"capacity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceArmLogAnalyticsClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Log Analytics Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_operationalinsights_cluster", *existing.ID)
	}

	parameters := operationalinsights.Cluster{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Identity: expandArmLogAnalyticsClusterIdentity(d.Get("identity").([]interface{})),
		ClusterProperties: &operationalinsights.ClusterProperties{
			NextLink:           utils.String(d.Get("next_link").(string)),
			KeyVaultProperties: expandArmLogAnalyticsClusterKeyVaultProperties(d.Get("key_vault_property").([]interface{})),
		},
		Sku:  expandArmLogAnalyticsClusterClusterSku(d.Get("sku").([]interface{})),
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("creating Log Analytics Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating future for Log Analytics Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Log Analytics Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Log Analytics Cluster %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)
	return resourceArmLogAnalyticsClusterRead(d, meta)
}

func resourceArmLogAnalyticsClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Log Analytics %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Log Analytics Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if err := d.Set("identity", flattenArmLogAnalyticsIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	if props := resp.ClusterProperties; props != nil {
		if err := d.Set("key_vault_property", flattenArmLogAnalyticsKeyVaultProperties(props.KeyVaultProperties)); err != nil {
			return fmt.Errorf("setting `key_vault_property`: %+v", err)
		}
		d.Set("next_link", props.NextLink)
		d.Set("cluster_id", props.ClusterID)
	}
	if err := d.Set("sku", flattenArmLogAnalyticsClusterSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `sku`: %+v", err)
	}
	d.Set("type", resp.Type)
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmLogAnalyticsClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsClusterID(d.Id())
	if err != nil {
		return err
	}

	parameters := operationalinsights.ClusterPatch{
		ClusterPatchProperties: &operationalinsights.ClusterPatchProperties{},
	}
	if d.HasChange("key_vault_property") {
		parameters.ClusterPatchProperties.KeyVaultProperties = expandArmLogAnalyticsClusterKeyVaultProperties(d.Get("key_vault_property").([]interface{}))
	}
	if d.HasChange("sku") {
		parameters.Sku = expandArmLogAnalyticsClusterClusterSku(d.Get("sku").([]interface{}))
	}
	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, parameters); err != nil {
		return fmt.Errorf("updating Log Analytics Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return resourceArmLogAnalyticsClusterRead(d, meta)
}

func resourceArmLogAnalyticsClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Log Analytics Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Log Analytics Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}

func expandArmLogAnalyticsClusterIdentity(input []interface{}) *operationalinsights.Identity {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &operationalinsights.Identity{
		Type: operationalinsights.IdentityType(v["type"].(string)),
	}
}

func expandArmLogAnalyticsClusterKeyVaultProperties(input []interface{}) *operationalinsights.KeyVaultProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &operationalinsights.KeyVaultProperties{
		KeyVaultURI: utils.String(v["key_vault_uri"].(string)),
		KeyName:     utils.String(v["key_name"].(string)),
		KeyVersion:  utils.String(v["key_version"].(string)),
	}
}

func expandArmLogAnalyticsClusterClusterSku(input []interface{}) *operationalinsights.ClusterSku {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &operationalinsights.ClusterSku{
		Capacity: utils.Int64(int64(v["capacity"].(int))),
		Name:     operationalinsights.ClusterSkuNameEnum(v["name"].(string)),
	}
}

func flattenArmLogAnalyticsIdentity(input *operationalinsights.Identity) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var t operationalinsights.IdentityType
	if input.Type != "" {
		t = input.Type
	}
	var principalId string
	if input.PrincipalID != nil {
		principalId = *input.PrincipalID
	}
	var tenantId string
	if input.TenantID != nil {
		tenantId = *input.TenantID
	}
	return []interface{}{
		map[string]interface{}{
			"type":         t,
			"principal_id": principalId,
			"tenant_id":    tenantId,
		},
	}
}

func flattenArmLogAnalyticsKeyVaultProperties(input *operationalinsights.KeyVaultProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var keyName string
	if input.KeyName != nil {
		keyName = *input.KeyName
	}
	var keyVaultUri string
	if input.KeyVaultURI != nil {
		keyVaultUri = *input.KeyVaultURI
	}
	var keyVersion string
	if input.KeyVersion != nil {
		keyVersion = *input.KeyVersion
	}
	return []interface{}{
		map[string]interface{}{
			"key_name":      keyName,
			"key_vault_uri": keyVaultUri,
			"key_version":   keyVersion,
		},
	}
}

func flattenArmLogAnalyticsClusterSku(input *operationalinsights.ClusterSku) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var name operationalinsights.ClusterSkuNameEnum
	if input.Name != "" {
		name = input.Name
	}
	var capacity int64
	if input.Capacity != nil {
		capacity = *input.Capacity
	}
	return []interface{}{
		map[string]interface{}{
			"name":     name,
			"capacity": capacity,
		},
	}
}
