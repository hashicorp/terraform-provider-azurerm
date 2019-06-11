package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-08-01/analysisservices"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"log"
	"regexp"
	//"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
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

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

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

			//"backup_blob_container_uri": {
			//	Type:         schema.TypeString,
			//	Optional:     true,
			//	Sensitive:    true,
			//	ValidateFunc: validate.URLIsHTTPOrHTTPS,
			//},

			"gateway_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			//"enable_power_bi_service": {
			//	Type:     schema.TypeBool,
			//	Optional: true,
			//},

			//"ipv4_firewall_rules": {
			//	Type:     schema.TypeList,
			//	Optional: true,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"name": {
			//				Type:     schema.TypeString,
			//				Required: true,
			//			},
			//			"range_start": {
			//				Type:         schema.TypeString,
			//				Required:     true,
			//				ValidateFunc: validate.IPv4Address,
			//			},
			//			"range_end": {
			//				Type:         schema.TypeString,
			//				Required:     true,
			//				ValidateFunc: validate.IPv4Address,
			//			},
			//		},
			//	},
			//},

			"querypool_connection_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateQuerypoolConnectionMode(),
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAnalysisServicesServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).analysisServicesServerClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Analysis Services Server creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		server, err := client.GetDetails(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(server.Response) {
				return fmt.Errorf("Error checking for presence of existing Analysis Services Server %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if server.ID != nil && *server.ID != "" {
			return tf.ImportAsExistsError("azurerm_analysis_services_server", *server.ID)
		}
	}

	sku := d.Get("sku").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))

	serverProperties := expandAnalysisServicesServerProperties(d)

	tags := d.Get("tags").(map[string]interface{})

	analysisServicesServer := analysisservices.Server{
		Name:             &name,
		Location:         &location,
		Sku:              expandAnalysisServicesServerSku(sku),
		ServerProperties: serverProperties,
		Tags:             expandTags(tags),
	}

	future, err := client.Create(ctx, resGroup, name, analysisServicesServer)
	if err != nil {
		return fmt.Errorf("Error creating Analysis Services Server %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Analysis Services Server %q (Resource Group %q): %+v", name, resGroup, err)
	}

	server, err := client.GetDetails(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if server.ID == nil {
		return fmt.Errorf("Cannot read ID of Analysis Services Server %q (Resource Group %q)", name, resGroup)
	}

	d.SetId(*server.ID)

	return resourceArmAnalysisServicesServerRead(d, meta)
}

func resourceArmAnalysisServicesServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).analysisServicesServerClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["servers"]

	server, err := client.GetDetails(ctx, resGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(server.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Analysis Services Server %q: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)

	if location := server.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("sku", flattenAnalysisServicesServerSku(server.Sku))

	if serverProps := server.ServerProperties; serverProps != nil {
		d.Set("admin_users", *flattenAnalysisServicesServerAdminUsers(server.AsAdministrators))

		//if serverProps.BackupBlobContainerURI != nil {
		//	d.Set("backup_blob_container_uri", *serverProps.BackupBlobContainerURI)
		//}

		d.Set("gateway_resource_id", *flattenAnalysisServicesServerGatewayDetails(serverProps))

		enablePowerBi, fwRules := flattenAnalysisServicesServerFirewallSettings(serverProps)
		d.Set("enable_power_bi_service", &enablePowerBi)
		d.Set("ipv4_firewall_rules", fwRules)

		d.Set("querypool_connection_mode", *flattenAnalysisServicesServerQuerypoolConnectionMode(serverProps))
	}

	flattenAndSetTags(d, server.Tags)

	return nil
}

func resourceArmAnalysisServicesServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).analysisServicesServerClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Analysis Services Server creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		server, err := client.GetDetails(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(server.Response) {
				return fmt.Errorf("Error checking for presence of existing Analysis Services Server %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if server.ID != nil && *server.ID != "" {
			return tf.ImportAsExistsError("azurerm_analysis_services_server", *server.ID)
		}
	}

	sku := d.Get("sku").(string)

	serverProperties := expandAnalysisServicesServerMutableProperties(d)

	tags := d.Get("tags").(map[string]interface{})

	analysisServicesServer := analysisservices.ServerUpdateParameters{
		Sku:                     expandAnalysisServicesServerSku(sku),
		Tags:                    expandTags(tags),
		ServerMutableProperties: serverProperties,
	}

	future, err := client.Update(ctx, resGroup, name, analysisServicesServer)
	if err != nil {
		return fmt.Errorf("Error creating Analysis Services Server %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Analysis Services Server %q (Resource Group %q): %+v", name, resGroup, err)
	}

	server, err := client.GetDetails(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if server.ID == nil {
		return fmt.Errorf("Cannot read ID of Analysis Services Server %q (Resource Group %q)", name, resGroup)
	}

	d.SetId(*server.ID)

	return resourceArmAnalysisServicesServerRead(d, meta)
}

func resourceArmAnalysisServicesServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).analysisServicesServerClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

func expandAnalysisServicesServerSku(sku string) *analysisservices.ResourceSku {
	return &analysisservices.ResourceSku{
		Name: &sku,
	}
}

func flattenAnalysisServicesServerSku(sku *analysisservices.ResourceSku) string {
	return *sku.Name
}

func expandAnalysisServicesServerProperties(d *schema.ResourceData) *analysisservices.ServerProperties {
	adminUsers := expandAnalysisServicesServerAdminUsers(d)

	serverProperties := analysisservices.ServerProperties{AsAdministrators: adminUsers}

	//if backupBlobContainerUri, ok := d.GetOk("backup_blob_container_uri"); ok {
	//	serverProperties.BackupBlobContainerURI = backupBlobContainerUri.(*string)
	//}

	if gatewayDetails := expandAnalysisServicesServerGatewayDetails(d); gatewayDetails != nil {
		serverProperties.GatewayDetails = gatewayDetails
	}

	serverProperties.IPV4FirewallSettings = expandAnalysisServicesServerFirewallSettings(d)

	if connectionMode := expandAnalysisServicesServerQuerypoolConnectionMode(d); connectionMode != nil {
		serverProperties.QuerypoolConnectionMode = *connectionMode
	}

	return &serverProperties
}

func expandAnalysisServicesServerMutableProperties(d *schema.ResourceData) *analysisservices.ServerMutableProperties {
	adminUsers := expandAnalysisServicesServerAdminUsers(d)

	serverProperties := analysisservices.ServerMutableProperties{AsAdministrators: adminUsers}

	//if backupBlobContainerUri, ok := d.GetOk("backup_blob_container_uri"); ok {
	//	serverProperties.BackupBlobContainerURI = backupBlobContainerUri.(*string)
	//}

	if gatewayDetails := expandAnalysisServicesServerGatewayDetails(d); gatewayDetails != nil {
		serverProperties.GatewayDetails = gatewayDetails
	}

	serverProperties.IPV4FirewallSettings = expandAnalysisServicesServerFirewallSettings(d)

	if connectionMode := expandAnalysisServicesServerQuerypoolConnectionMode(d); connectionMode != nil {
		serverProperties.QuerypoolConnectionMode = *connectionMode
	}

	return &serverProperties
}

func expandAnalysisServicesServerAdminUsers(d *schema.ResourceData) *analysisservices.ServerAdministrators {
	adminUsers := d.Get("admin_users").([]interface{})
	members := make([]string, 0)

	for _, admin := range adminUsers {
		if adm, ok := admin.(string); ok {
			members = append(members, adm)
		}
	}

	return &analysisservices.ServerAdministrators{Members: &members}
}

func expandAnalysisServicesServerGatewayDetails(d *schema.ResourceData) *analysisservices.GatewayDetails {
	if gatewayResourceId, ok := d.GetOk("gateway_resource_id"); ok {
		rId := gatewayResourceId.(string)
		return &analysisservices.GatewayDetails{GatewayResourceID: &rId}
	}

	return nil
}

func expandAnalysisServicesServerQuerypoolConnectionMode(d *schema.ResourceData) *analysisservices.ConnectionMode {
	if querypoolConnectionMode, ok := d.GetOk("querypool_connection_mode"); ok {
		connMode := analysisservices.ConnectionMode(querypoolConnectionMode.(string))
		return &connMode
	}

	return nil
}

func expandAnalysisServicesServerFirewallSettings(d *schema.ResourceData) *analysisservices.IPv4FirewallSettings {
	firewallSettings := analysisservices.IPv4FirewallSettings{}

	if enablePowerBi, exists := d.GetOkExists("enable_power_bi_service"); exists {
		firewallSettings.EnablePowerBIService = enablePowerBi.(*bool)
	}

	firewallRules := d.Get("ipv4_firewall_rules").([]interface{})
	if len(firewallRules) > 0 {
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
	}

	return &firewallSettings
}

func flattenAnalysisServicesServerAdminUsers(serverAdministrators *analysisservices.ServerAdministrators) *[]string {
	if serverAdministrators == nil || serverAdministrators.Members == nil {
		return &[]string{}
	}

	return serverAdministrators.Members
}

func flattenAnalysisServicesServerGatewayDetails(serverProperties *analysisservices.ServerProperties) *string {
	if serverProperties.GatewayDetails == nil {
		return nil
	}

	return serverProperties.GatewayDetails.GatewayResourceID
}

func flattenAnalysisServicesServerFirewallSettings(serverProperties *analysisservices.ServerProperties) (enablePowerBi *bool, fwRules []interface{}) {
	if serverProperties.IPV4FirewallSettings == nil {
		return nil, nil
	}

	firewallSettings := serverProperties.IPV4FirewallSettings

	if firewallSettings.FirewallRules != nil && len(*firewallSettings.FirewallRules) > 0 {
		fwRules := make([]map[string]interface{}, len(*firewallSettings.FirewallRules))
		for i, fwRule := range *firewallSettings.FirewallRules {
			fwRules[i]["name"] = *fwRule.FirewallRuleName
			fwRules[i]["range_start"] = *fwRule.RangeStart
			fwRules[i]["range_end"] = *fwRule.RangeEnd
		}
	}

	return firewallSettings.EnablePowerBIService, fwRules
}

func flattenAnalysisServicesServerQuerypoolConnectionMode(serverProperties *analysisservices.ServerProperties) *string {
	connMode := string(serverProperties.QuerypoolConnectionMode)
	return &connMode
}
