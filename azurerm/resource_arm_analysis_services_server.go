package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-08-01/analysisservices"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAnalysisServicesServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAnalysisServicesServerCreate,
		Read:   resourceArmAnalysisServicesServerRead,
		Update: resourceArmAnalysisServicesServerUpdate,
		Delete: resourceArmAnalysisServicesServerDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAnalysisServicesServerName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"D1",
					"B1",
					"B2",
					"S0",
					"S1",
					"S2",
					"S4",
					"S8",
					"S9",
				}, false),
			},

			"admin_users": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"enable_power_bi_service": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"ipv4_firewall_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"range_start": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.IPv4Address,
						},
						"range_end": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.IPv4Address,
						},
					},
				},
			},

			"querypool_connection_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateQuerypoolConnectionMode(),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmAnalysisServicesServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).analysisservices.ServerClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Analysis Services Server creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		server, err := client.GetDetails(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(server.Response) {
				return fmt.Errorf("Error checking for presence of existing Analysis Services Server %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if server.ID != nil && *server.ID != "" {
			return tf.ImportAsExistsError("azurerm_analysis_services_server", *server.ID)
		}
	}

	sku := d.Get("sku").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))

	serverProperties := expandAnalysisServicesServerProperties(d)

	t := d.Get("tags").(map[string]interface{})

	analysisServicesServer := analysisservices.Server{
		Name:             &name,
		Location:         &location,
		Sku:              &analysisservices.ResourceSku{Name: &sku},
		ServerProperties: serverProperties,
		Tags:             tags.Expand(t),
	}

	future, err := client.Create(ctx, resourceGroup, name, analysisServicesServer)
	if err != nil {
		return fmt.Errorf("Error creating Analysis Services Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Analysis Services Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, getDetailsErr := client.GetDetails(ctx, resourceGroup, name)
	if getDetailsErr != nil {
		return fmt.Errorf("Error retrieving Analytics Services Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Analytics Services Server %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmAnalysisServicesServerRead(d, meta)
}

func resourceArmAnalysisServicesServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).analysisservices.ServerClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["servers"]

	server, err := client.GetDetails(ctx, resourceGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(server.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Analytics Services Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if location := server.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if server.Sku != nil {
		d.Set("sku", server.Sku.Name)
	}

	if serverProps := server.ServerProperties; serverProps != nil {
		if serverProps.AsAdministrators == nil || serverProps.AsAdministrators.Members == nil {
			d.Set("admin_users", []string{})
		} else {
			d.Set("admin_users", *serverProps.AsAdministrators.Members)
		}

		enablePowerBi, fwRules := flattenAnalysisServicesServerFirewallSettings(serverProps)
		d.Set("enable_power_bi_service", enablePowerBi)
		if err := d.Set("ipv4_firewall_rule", fwRules); err != nil {
			return fmt.Errorf("Error setting `ipv4_firewall_rule`: %s", err)
		}

		d.Set("querypool_connection_mode", string(serverProps.QuerypoolConnectionMode))
	}

	return tags.FlattenAndSet(d, server.Tags)
}

func resourceArmAnalysisServicesServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).analysisservices.ServerClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Analysis Services Server creation.")

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["servers"]

	serverProperties := expandAnalysisServicesServerMutableProperties(d)
	sku := d.Get("sku").(string)
	t := d.Get("tags").(map[string]interface{})

	analysisServicesServer := analysisservices.ServerUpdateParameters{
		Sku:                     &analysisservices.ResourceSku{Name: &sku},
		Tags:                    tags.Expand(t),
		ServerMutableProperties: serverProperties,
	}

	future, err := client.Update(ctx, resourceGroup, name, analysisServicesServer)
	if err != nil {
		return fmt.Errorf("Error creating Analysis Services Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Analysis Services Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return resourceArmAnalysisServicesServerRead(d, meta)
}

func resourceArmAnalysisServicesServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).analysisservices.ServerClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["servers"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Analysis Services Server %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Analysis Services Server %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func validateAnalysisServicesServerName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-z][0-9a-z]{2,62}$`).Match([]byte(value)) {
		errors = append(errors, fmt.Errorf("%q must begin with a letter, be lowercase alphanumeric, and be between 3 and 63 characters in length", k))
	}

	return warnings, errors
}

func validateQuerypoolConnectionMode() schema.SchemaValidateFunc {
	connectionModes := make([]string, len(analysisservices.PossibleConnectionModeValues()))
	for i, v := range analysisservices.PossibleConnectionModeValues() {
		connectionModes[i] = string(v)
	}

	return validation.StringInSlice(connectionModes, true)
}

func expandAnalysisServicesServerProperties(d *schema.ResourceData) *analysisservices.ServerProperties {
	adminUsers := expandAnalysisServicesServerAdminUsers(d)

	serverProperties := analysisservices.ServerProperties{AsAdministrators: adminUsers}

	serverProperties.IPV4FirewallSettings = expandAnalysisServicesServerFirewallSettings(d)

	if querypoolConnectionMode, ok := d.GetOk("querypool_connection_mode"); ok {
		serverProperties.QuerypoolConnectionMode = analysisservices.ConnectionMode(querypoolConnectionMode.(string))
	}

	return &serverProperties
}

func expandAnalysisServicesServerMutableProperties(d *schema.ResourceData) *analysisservices.ServerMutableProperties {
	adminUsers := expandAnalysisServicesServerAdminUsers(d)

	serverProperties := analysisservices.ServerMutableProperties{AsAdministrators: adminUsers}

	serverProperties.IPV4FirewallSettings = expandAnalysisServicesServerFirewallSettings(d)

	serverProperties.QuerypoolConnectionMode = analysisservices.ConnectionMode(d.Get("querypool_connection_mode").(string))

	return &serverProperties
}

func expandAnalysisServicesServerAdminUsers(d *schema.ResourceData) *analysisservices.ServerAdministrators {
	adminUsers := d.Get("admin_users").(*schema.Set)
	members := make([]string, 0)

	for _, admin := range adminUsers.List() {
		if adm, ok := admin.(string); ok {
			members = append(members, adm)
		}
	}

	return &analysisservices.ServerAdministrators{Members: &members}
}

func expandAnalysisServicesServerFirewallSettings(d *schema.ResourceData) *analysisservices.IPv4FirewallSettings {
	firewallSettings := analysisservices.IPv4FirewallSettings{
		EnablePowerBIService: utils.Bool(d.Get("enable_power_bi_service").(bool)),
	}

	firewallRules := d.Get("ipv4_firewall_rule").([]interface{})

	fwRules := make([]analysisservices.IPv4FirewallRule, len(firewallRules))
	for i, v := range firewallRules {
		fwRule := v.(map[string]interface{})
		fwRules[i] = analysisservices.IPv4FirewallRule{
			FirewallRuleName: utils.String(fwRule["name"].(string)),
			RangeStart:       utils.String(fwRule["range_start"].(string)),
			RangeEnd:         utils.String(fwRule["range_end"].(string)),
		}
	}
	firewallSettings.FirewallRules = &fwRules

	return &firewallSettings
}

func flattenAnalysisServicesServerFirewallSettings(serverProperties *analysisservices.ServerProperties) (enablePowerBi *bool, fwRules []interface{}) {
	if serverProperties == nil || serverProperties.IPV4FirewallSettings == nil {
		return utils.Bool(false), make([]interface{}, 0)
	}

	firewallSettings := serverProperties.IPV4FirewallSettings

	enablePowerBi = utils.Bool(false)
	if firewallSettings.EnablePowerBIService != nil {
		enablePowerBi = firewallSettings.EnablePowerBIService
	}

	fwRules = make([]interface{}, 0)
	if firewallSettings.FirewallRules != nil {
		for _, fwRule := range *firewallSettings.FirewallRules {
			output := make(map[string]interface{})
			if fwRule.FirewallRuleName != nil {
				output["name"] = *fwRule.FirewallRuleName
			}

			if fwRule.RangeStart != nil {
				output["range_start"] = *fwRule.RangeStart
			}

			if fwRule.RangeEnd != nil {
				output["range_end"] = *fwRule.RangeEnd
			}

			fwRules = append(fwRules, output)
		}
	}

	return enablePowerBi, fwRules
}
