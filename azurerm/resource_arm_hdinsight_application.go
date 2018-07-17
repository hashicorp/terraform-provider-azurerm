package azurerm

import (
	"context"
	"fmt"
	"log"

	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2015-03-01-preview/hdinsight"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmHDInsightApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHDInsightApplicationCreate,
		Read:   resourceArmHDInsightApplicationRead,
		Delete: resourceArmHDInsightApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_name": resourceGroupNameSchema(),

			"marketplace_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"edge_node": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hardware_profile": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vm_size": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},

			"install_script_action": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"uri": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"roles": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"edgenode",
									"headnode",
									"workernode",
									"zookeepernode",
								}, false),
							},
						},
					},
				},
			},

			"uninstall_script_action": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"uri": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"roles": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"edgenode",
									"headnode",
									"workernode",
									"zookeepernode",
								}, false),
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmHDInsightApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdinsightApplicationsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for HDInsight Application creation.")

	name := d.Get("name").(string)
	clusterName := d.Get("cluster_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	marketplaceIdentifier := d.Get("marketplace_identifier").(string)

	computeProfile := expandHDInsightApplicationComputeProfile(d)
	installScriptActions := expandHDInsightApplicationScriptActions(d.Get("install_script_action").([]interface{}))
	uninstallScriptActions := expandHDInsightApplicationScriptActions(d.Get("uninstall_script_action").([]interface{}))

	properties := hdinsight.ApplicationProperties{
		ApplicationType:        utils.String("CustomApplication"),
		MarketplaceIdentifier:  utils.String(marketplaceIdentifier),
		ComputeProfile:         computeProfile,
		InstallScriptActions:   installScriptActions,
		UninstallScriptActions: uninstallScriptActions,
	}
	future, err := client.Create(ctx, resourceGroup, clusterName, name, properties)
	if err != nil {
		if response.WasOK(future.Response()) {
			err = resourceArmHDInsightApplicationReadError(client, ctx, resourceGroup, clusterName, name)
		}

		return fmt.Errorf("Error creating HDInsight Application %q (Cluster %q / Resource Group %q): %+v", name, clusterName, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasOK(future.Response()) {
			err = resourceArmHDInsightApplicationReadError(client, ctx, resourceGroup, clusterName, name)
		}

		return fmt.Errorf("Error waiting for creation of HDInsight Application %q (Cluster %q / Resource Group %q): %+v", name, clusterName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, clusterName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving HDInsights Application %q (Cluster %q / Resource Group %q): %+v", name, clusterName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("[ERROR] Cannot read HDInsight Application %q (Cluster %q / Resource Group %q) ID", name, clusterName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmHDInsightApplicationRead(d, meta)
}

func resourceArmHDInsightApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdinsightApplicationsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	clusterName := id.Path["clusters"]
	name := id.Path["applications"]

	resp, err := client.Get(ctx, resourceGroup, clusterName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] HDInsight Application %q (Cluster %q / Resource Group %q) was not found - removing from state!", name, clusterName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HDInsight Application %q (Cluster %q / Resource Group %q): %+v", name, clusterName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("cluster_name", clusterName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.Properties; props != nil {
		d.Set("marketplace_identifier", props.MarketplaceIdentifier)

		computeProfile := flattenHDInsightApplicationComputeProfile(props.ComputeProfile)
		if err := d.Set("edge_node", computeProfile); err != nil {
			return fmt.Errorf("Error setting `edge_node`: %+v", err)
		}

		installActions := flattenHDInsightApplicationScriptActions(props.InstallScriptActions)
		if err := d.Set("install_script_action", installActions); err != nil {
			return fmt.Errorf("Error setting `install_script_action`: %+v", err)
		}

		uninstallActions := flattenHDInsightApplicationScriptActions(props.UninstallScriptActions)
		if err := d.Set("uninstall_script_action", uninstallActions); err != nil {
			return fmt.Errorf("Error setting `uninstall_script_action`: %+v", err)
		}
	}

	return nil
}

func resourceArmHDInsightApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdinsightApplicationsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	clusterName := id.Path["clusters"]
	name := id.Path["applications"]

	future, err := client.Delete(ctx, resourceGroup, clusterName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting HDInsight Application %q (Cluster %q / Resource Group %q): %+v", name, clusterName, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error waiting for deletion of HDInsight Application %q (Cluster %q / Resource Group %q): %+v", name, clusterName, resourceGroup, err)
	}

	return nil
}

func expandHDInsightApplicationComputeProfile(d *schema.ResourceData) *hdinsight.ComputeProfile {
	vs := d.Get("edge_node").([]interface{})
	v := vs[0].(map[string]interface{})

	hardwareProfiles := v["hardware_profile"].([]interface{})
	hardwareProfile := hardwareProfiles[0].(map[string]interface{})
	vmSize := hardwareProfile["vm_size"].(string)

	return &hdinsight.ComputeProfile{
		Roles: &[]hdinsight.Role{
			{
				// these two have to be hard-coded
				Name:                utils.String("edgenode"),
				TargetInstanceCount: utils.Int32(int32(1)),
				HardwareProfile: &hdinsight.HardwareProfile{
					VMSize: utils.String(vmSize),
				},
			},
		},
	}
}

func flattenHDInsightApplicationComputeProfile(input *hdinsight.ComputeProfile) []interface{} {
	roles := make([]interface{}, 0)

	for _, v := range *input.Roles {
		if v.Name == nil || strings.EqualFold(*v.Name, "edgenode") {
			continue
		}

		role := make(map[string]interface{}, 0)

		hardwareProfiles := make([]interface{}, 0)
		if profile := v.HardwareProfile; profile != nil {
			if size := profile.VMSize; size != nil {
				hardwareProfile := map[string]interface{}{
					"vm_size": *size,
				}
				hardwareProfiles = append(hardwareProfiles, hardwareProfile)
			}
		}
		role["hardware_profile"] = hardwareProfiles

		roles = append(roles, role)
	}

	return roles
}

func expandHDInsightApplicationScriptActions(input []interface{}) *[]hdinsight.RuntimeScriptAction {
	actions := make([]hdinsight.RuntimeScriptAction, 0)

	for _, v := range input {
		val := v.(map[string]interface{})

		name := val["name"].(string)
		uri := val["uri"].(string)

		rolesRaw := val["roles"].([]interface{})
		roles := make([]string, 0)
		for _, v := range rolesRaw {
			role := v.(string)
			roles = append(roles, role)
		}

		action := hdinsight.RuntimeScriptAction{
			Name:  utils.String(name),
			URI:   utils.String(uri),
			Roles: &roles,
		}

		actions = append(actions, action)
	}

	return &actions
}

func flattenHDInsightApplicationScriptActions(input *[]hdinsight.RuntimeScriptAction) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	for _, action := range *input {
		output := make(map[string]interface{}, 0)

		if name := action.Name; name != nil {
			output["name"] = *name
		}

		if uri := action.URI; uri != nil {
			output["uri"] = *uri
		}

		roles := make([]string, 0)
		if action.Roles != nil {
			for _, r := range *action.Roles {
				roles = append(roles, r)
			}
		}
		output["roles"] = roles
		outputs = append(outputs, output)
	}

	return outputs
}

func resourceArmHDInsightApplicationReadError(client hdinsight.ApplicationsClient, ctx context.Context, resourceGroup, clusterName, name string) error {
	// the HDInsights errors returns errors as a 200 to the SDK
	// meaning we need to retrieve the cluster to get the error
	resp, err := client.Get(ctx, resourceGroup, clusterName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving HDInsights Application %q (Cluster %q / Resource Group %q): %+v", name, clusterName, resourceGroup, err)
	}

	if props := resp.Properties; props != nil {
		if errors := props.Errors; errors != nil {
			var err error

			for _, e := range *errors {
				if message := e.Message; message != nil {
					err = multierror.Append(err, fmt.Errorf(*message))
				}
			}

			return fmt.Errorf("Error updating HDInsight Application %q (Cluster %q / Resource Group %q): %+v", name, clusterName, resourceGroup, err)
		}
	}

	return nil
}
