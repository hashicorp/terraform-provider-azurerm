package azurerm

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-06-01/netapp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	aznetapp "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetAppVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetAppVolumeCreateUpdate,
		Read:   resourceArmNetAppVolumeRead,
		Update: resourceArmNetAppVolumeCreateUpdate,
		Delete: resourceArmNetAppVolumeDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: aznetapp.ValidateNetAppVolumeName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: aznetapp.ValidateNetAppAccountName,
			},

			"pool_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: aznetapp.ValidateNetAppPoolName,
			},

			"volume_path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: aznetapp.ValidateNetAppVolumeVolumePath,
			},

			"service_level": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(netapp.Premium),
					string(netapp.Standard),
					string(netapp.Ultra),
				}, false),
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"storage_quota_in_gb": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 4096),
			},

			"export_policy_rule": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_index": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 5),
						},
						"allowed_clients": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.CIDR,
							},
						},
						"cifs_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"nfsv3_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"nfsv4_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"unix_read_only": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"unix_read_write": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmNetAppVolumeCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)
	poolName := d.Get("pool_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, accountName, poolName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_netapp_volume", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	volumePath := d.Get("volume_path").(string)
	serviceLevel := d.Get("service_level").(string)
	subnetId := d.Get("subnet_id").(string)
	storageQuotaInGB := int64(d.Get("storage_quota_in_gb").(int) * 1073741824)
	exportPolicyRule := d.Get("export_policy_rule").(*schema.Set).List()

	parameters := netapp.Volume{
		Location: utils.String(location),
		VolumeProperties: &netapp.VolumeProperties{
			CreationToken:  utils.String(volumePath),
			ServiceLevel:   netapp.ServiceLevel(serviceLevel),
			SubnetID:       utils.String(subnetId),
			UsageThreshold: utils.Int64(storageQuotaInGB),
			ExportPolicy:   expandArmNetAppVolumeExportPolicyRule(exportPolicyRule),
		},
	}

	future, err := client.CreateOrUpdate(ctx, parameters, resourceGroup, accountName, poolName, name)
	if err != nil {
		return fmt.Errorf("Error creating NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, accountName, poolName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read NetApp Volume %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmNetAppVolumeRead(d, meta)
}

func resourceArmNetAppVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["netAppAccounts"]
	poolName := id.Path["capacityPools"]
	name := id.Path["volumes"]

	resp, err := client.Get(ctx, resourceGroup, accountName, poolName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] NetApp Volumes %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading NetApp Volumes %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("account_name", accountName)
	d.Set("pool_name", poolName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.VolumeProperties; props != nil {
		d.Set("volume_path", props.CreationToken)
		d.Set("service_level", props.ServiceLevel)
		d.Set("subnet_id", props.SubnetID)

		if props.UsageThreshold != nil {
			d.Set("storage_quota_in_gb", *props.UsageThreshold/1073741824)
		}
		if err := d.Set("export_policy_rule", flattenArmNetAppVolumeExportPolicyRule(props.ExportPolicy)); err != nil {
			return fmt.Errorf("Error setting `export_policy_rule`: %+v", err)
		}
	}

	return nil
}

func resourceArmNetAppVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["netAppAccounts"]
	poolName := id.Path["capacityPools"]
	name := id.Path["volumes"]

	if _, err = client.Delete(ctx, resourceGroup, accountName, poolName, name); err != nil {
		return fmt.Errorf("Error deleting NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return waitForNetAppVolumeToBeDeleted(ctx, client, resourceGroup, accountName, poolName, name)
}

func waitForNetAppVolumeToBeDeleted(ctx context.Context, client *netapp.VolumesClient, resourceGroup, accountName, poolName, name string) error {
	// The resource NetApp Volume depends on the resource NetApp Pool.
	// Although the delete API returns 404 which means the NetApp Volume resource has been deleted.
	// Then it tries to immediately delete NetApp Pool but it still throws error `Can not delete resource before nested resources are deleted.`
	// In this case we're going to try triggering the Deletion again, in-case it didn't work prior to this attempt.
	// For more details, see related Bug: https://github.com/Azure/azure-sdk-for-go/issues/6485

	log.Printf("[DEBUG] Waiting for NetApp Volume Provisioning Service %q (Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200", "202"},
		Target:  []string{"404"},
		Refresh: netappVolumeDeleteStateRefreshFunc(ctx, client, resourceGroup, accountName, poolName, name),
		Timeout: 20 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NetApp Volume Provisioning Service %q (Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func netappVolumeDeleteStateRefreshFunc(ctx context.Context, client *netapp.VolumesClient, resourceGroupName string, accountName string, poolName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, accountName, poolName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("Error retrieving NetApp Volume %q (Resource Group %q): %s", name, resourceGroupName, err)
			}
		}

		if _, err := client.Delete(ctx, resourceGroupName, accountName, poolName, name); err != nil {
			log.Printf("Error reissuing NetApp Volume %q delete request (Resource Group %q): %+v", name, resourceGroupName, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func expandArmNetAppVolumeExportPolicyRule(input []interface{}) *netapp.VolumePropertiesExportPolicy {
	results := make([]netapp.ExportPolicyRule, 0)
	for _, item := range input {
		if item != nil {
			v := item.(map[string]interface{})
			ruleIndex := int32(v["rule_index"].(int))
			allowedClients := strings.Join(*utils.ExpandStringSlice(v["allowed_clients"].(*schema.Set).List()), ",")
			cifsEnabled := v["cifs_enabled"].(bool)
			nfsv3Enabled := v["nfsv3_enabled"].(bool)
			nfsv4Enabled := v["nfsv4_enabled"].(bool)
			unixReadOnly := v["unix_read_only"].(bool)
			unixReadWrite := v["unix_read_write"].(bool)

			result := netapp.ExportPolicyRule{
				AllowedClients: utils.String(allowedClients),
				Cifs:           utils.Bool(cifsEnabled),
				Nfsv3:          utils.Bool(nfsv3Enabled),
				Nfsv4:          utils.Bool(nfsv4Enabled),
				RuleIndex:      utils.Int32(ruleIndex),
				UnixReadOnly:   utils.Bool(unixReadOnly),
				UnixReadWrite:  utils.Bool(unixReadWrite),
			}

			results = append(results, result)
		}
	}

	return &netapp.VolumePropertiesExportPolicy{
		Rules: &results,
	}
}

func flattenArmNetAppVolumeExportPolicyRule(input *netapp.VolumePropertiesExportPolicy) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || input.Rules == nil {
		return results
	}

	for _, item := range *input.Rules {
		ruleIndex := int32(0)
		if v := item.RuleIndex; v != nil {
			ruleIndex = *v
		}
		allowedClients := []string{}
		if v := item.AllowedClients; v != nil {
			allowedClients = strings.Split(*v, ",")
		}
		cifsEnabled := false
		if v := item.Cifs; v != nil {
			cifsEnabled = *v
		}
		nfsv3Enabled := false
		if v := item.Nfsv3; v != nil {
			nfsv3Enabled = *v
		}
		nfsv4Enabled := false
		if v := item.Nfsv4; v != nil {
			nfsv4Enabled = *v
		}
		unixReadOnly := false
		if v := item.UnixReadOnly; v != nil {
			unixReadOnly = *v
		}
		unixReadWrite := false
		if v := item.UnixReadWrite; v != nil {
			unixReadWrite = *v
		}

		results = append(results, map[string]interface{}{
			"rule_index":      ruleIndex,
			"allowed_clients": utils.FlattenStringSlice(&allowedClients),
			"cifs_enabled":    cifsEnabled,
			"nfsv3_enabled":   nfsv3Enabled,
			"nfsv4_enabled":   nfsv4Enabled,
			"unix_read_only":  unixReadOnly,
			"unix_read_write": unixReadWrite,
		})
	}

	return results
}
