package azurerm

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-06-01/netapp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	aznetapp "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp"
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
				ValidateFunc: aznetapp.ValidateNetAppAccountName,
			},

			"pool_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: aznetapp.ValidateNetAppPoolName,
			},

			"creation_token": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: aznetapp.ValidateNetAppVolumeCreationToken,
			},

			"service_level": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(netapp.Premium),
					string(netapp.Standard),
					string(netapp.Ultra),
				}, true),
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"usage_threshold": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 4096),
			},

			"export_policy_rule": {
				Type:     schema.TypeList,
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.CIDR,
						},
						"cifs": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"nfsv3": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"nfsv4": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"unix_read_only": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"unix_read_write": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmNetAppVolumeCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Netapp.VolumeClient
	ctx := meta.(*ArmClient).StopContext

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
	creationToken := d.Get("creation_token").(string)
	serviceLevel := d.Get("service_level").(string)
	subnetId := d.Get("subnet_id").(string)
	usageThreshold := int64(d.Get("usage_threshold").(int) * 1073741824)
	exportPolicyRule := d.Get("export_policy_rule").([]interface{})

	parameters := netapp.Volume{
		Location: utils.String(location),
		VolumeProperties: &netapp.VolumeProperties{
			CreationToken:  utils.String(creationToken),
			ServiceLevel:   netapp.ServiceLevel(serviceLevel),
			SubnetID:       utils.String(subnetId),
			UsageThreshold: utils.Int64(usageThreshold),
			ExportPolicy:   expandArmExportPolicyRule(exportPolicyRule),
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
	if resp.ID == nil {
		return fmt.Errorf("Cannot read NetApp Volume %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmNetAppVolumeRead(d, meta)
}

func resourceArmNetAppVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Netapp.VolumeClient
	ctx := meta.(*ArmClient).StopContext

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
		d.Set("creation_token", props.CreationToken)
		d.Set("service_level", props.ServiceLevel)
		d.Set("subnet_id", props.SubnetID)

		if props.UsageThreshold != nil {
			d.Set("usage_threshold", *props.UsageThreshold/1073741824)
		}
		if props.ExportPolicy != nil {
			if err := d.Set("export_policy_rule", flattenArmExportPolicyRule(props.ExportPolicy.Rules)); err != nil {
				return fmt.Errorf("Error setting `export_policy_rule`: %+v", err)
			}
		}
	}

	return nil
}

func resourceArmNetAppVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Netapp.VolumeClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["netAppAccounts"]
	poolName := id.Path["capacityPools"]
	name := id.Path["volumes"]

	_, err = client.Delete(ctx, resourceGroup, accountName, poolName, name)
	if err != nil {
		return fmt.Errorf("Error deleting NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return waitForNetAppVolumeToBeDeleted(ctx, client, resourceGroup, accountName, poolName, name)
}

func waitForNetAppVolumeToBeDeleted(ctx context.Context, client *netapp.VolumesClient, resourceGroup, accountName, poolName, name string) error {
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

func expandArmExportPolicyRule(input []interface{}) *netapp.VolumePropertiesExportPolicy {
	results := make([]netapp.ExportPolicyRule, 0)
	for _, item := range input {
		if item != nil {
			v := item.(map[string]interface{})
			ruleIndex := int32(v["rule_index"].(int))
			allowedClients := v["allowed_clients"].(string)
			cifs := v["cifs"].(bool)
			nfsv3 := v["nfsv3"].(bool)
			nfsv4 := v["nfsv4"].(bool)
			unixReadOnly := v["unix_read_only"].(bool)
			unixReadWrite := v["unix_read_write"].(bool)

			result := netapp.ExportPolicyRule{
				AllowedClients: utils.String(allowedClients),
				Cifs:           utils.Bool(cifs),
				Nfsv3:          utils.Bool(nfsv3),
				Nfsv4:          utils.Bool(nfsv4),
				RuleIndex:      utils.Int32(ruleIndex),
				UnixReadOnly:   utils.Bool(unixReadOnly),
				UnixReadWrite:  utils.Bool(unixReadWrite),
			}

			results = append(results, result)
		}
	}

	result := netapp.VolumePropertiesExportPolicy{
		Rules: &results,
	}

	return &result
}

func flattenArmExportPolicyRule(input *[]netapp.ExportPolicyRule) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		c := make(map[string]interface{})

		if v := item.RuleIndex; v != nil {
			c["rule_index"] = *v
		}
		if v := item.AllowedClients; v != nil {
			c["allowed_clients"] = *v
		}
		if v := item.Cifs; v != nil {
			c["cifs"] = *v
		}
		if v := item.Nfsv3; v != nil {
			c["nfsv3"] = *v
		}
		if v := item.Nfsv4; v != nil {
			c["nfsv4"] = *v
		}
		if v := item.UnixReadOnly; v != nil {
			c["unix_read_only"] = *v
		}
		if v := item.UnixReadWrite; v != nil {
			c["unix_read_write"] = *v
		}

		results = append(results, c)
	}

	return results
}
