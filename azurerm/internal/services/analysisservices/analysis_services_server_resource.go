package analysisservices

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-08-01/analysisservices"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/analysisservices/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAnalysisServicesServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAnalysisServicesServerCreate,
		Read:   resourceArmAnalysisServicesServerRead,
		Update: resourceArmAnalysisServicesServerUpdate,
		Delete: resourceArmAnalysisServicesServerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ServerID(id)
			return err
		}),

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
					"S8v2",
					"S9v2",
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
				Type:     schema.TypeSet,
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
				Set: hashAnalysisServicesServerIpv4FirewallRule,
			},

			"querypool_connection_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateQuerypoolConnectionMode(),
			},

			"backup_blob_container_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"server_full_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmAnalysisServicesServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AnalysisServices.ServerClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Analysis Services Server creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
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
	client := meta.(*clients.Client).AnalysisServices.ServerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerID(d.Id())
	if err != nil {
		return err
	}

	server, err := client.GetDetails(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		if utils.ResponseWasNotFound(server.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Analytics Services Server %q (Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	d.Set("name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := server.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if server.Sku != nil {
		d.Set("sku", server.Sku.Name)
	}

	if serverProps := server.ServerProperties; serverProps != nil {
		if serverProps.AsAdministrators == nil {
			d.Set("admin_users", []string{})
		} else {
			d.Set("admin_users", serverProps.AsAdministrators.Members)
		}

		enablePowerBi, fwRules := flattenAnalysisServicesServerFirewallSettings(serverProps)
		d.Set("enable_power_bi_service", enablePowerBi)
		if err := d.Set("ipv4_firewall_rule", fwRules); err != nil {
			return fmt.Errorf("Error setting `ipv4_firewall_rule`: %s", err)
		}

		d.Set("querypool_connection_mode", string(serverProps.QuerypoolConnectionMode))
		d.Set("server_full_name", serverProps.ServerFullName)

		if containerUri, ok := d.GetOk("backup_blob_container_uri"); ok {
			d.Set("backup_blob_container_uri", containerUri)
		}
	}

	return tags.FlattenAndSet(d, server.Tags)
}

func resourceArmAnalysisServicesServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AnalysisServices.ServerClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Analysis Services Server update.")

	id, err := parse.ServerID(d.Id())
	if err != nil {
		return err
	}

	serverResp, err := client.GetDetails(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		return fmt.Errorf("Error retrieving Analysis Services Server %q (Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	if serverResp.State != analysisservices.StateSucceeded && serverResp.State != analysisservices.StatePaused {
		return fmt.Errorf("Error updating Analysis Services Server %q (Resource Group %q): State must be either Succeeded or Paused but got %q", id.ServerName, id.ResourceGroup, serverResp.State)
	}

	isPaused := serverResp.State == analysisservices.StatePaused

	if isPaused {
		resumeFuture, err := client.Resume(ctx, id.ResourceGroup, id.ServerName)
		if err != nil {
			return fmt.Errorf("Error starting Analysis Services Server %q (Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
		}

		if err = resumeFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for Analysis Services Server starting completion %q (Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
		}
	}

	serverProperties := expandAnalysisServicesServerMutableProperties(d)
	sku := d.Get("sku").(string)
	t := d.Get("tags").(map[string]interface{})

	analysisServicesServer := analysisservices.ServerUpdateParameters{
		Sku:                     &analysisservices.ResourceSku{Name: &sku},
		Tags:                    tags.Expand(t),
		ServerMutableProperties: serverProperties,
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ServerName, analysisServicesServer)
	if err != nil {
		return fmt.Errorf("Error creating Analysis Services Server %q (Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Analysis Services Server %q (Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	if isPaused {
		suspendFuture, err := client.Suspend(ctx, id.ResourceGroup, id.ServerName)
		if err != nil {
			return fmt.Errorf("Error re-pausing Analysis Services Server %q (Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
		}

		if err = suspendFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for Analysis Services Server pausing completion %q (Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
		}
	}

	return resourceArmAnalysisServicesServerRead(d, meta)
}

func resourceArmAnalysisServicesServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AnalysisServices.ServerClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		return fmt.Errorf("Error deleting Analysis Services Server %q (Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Analysis Services Server %q (Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
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

	if containerUri, ok := d.GetOk("backup_blob_container_uri"); ok {
		serverProperties.BackupBlobContainerURI = utils.String(containerUri.(string))
	}

	return &serverProperties
}

func expandAnalysisServicesServerMutableProperties(d *schema.ResourceData) *analysisservices.ServerMutableProperties {
	adminUsers := expandAnalysisServicesServerAdminUsers(d)

	serverProperties := analysisservices.ServerMutableProperties{AsAdministrators: adminUsers}

	serverProperties.IPV4FirewallSettings = expandAnalysisServicesServerFirewallSettings(d)

	serverProperties.QuerypoolConnectionMode = analysisservices.ConnectionMode(d.Get("querypool_connection_mode").(string))

	if containerUri, ok := d.GetOk("backup_blob_container_uri"); ok {
		serverProperties.BackupBlobContainerURI = utils.String(containerUri.(string))
	}

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

	firewallRules := d.Get("ipv4_firewall_rule").(*schema.Set).List()

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

func flattenAnalysisServicesServerFirewallSettings(serverProperties *analysisservices.ServerProperties) (*bool, *schema.Set) {
	if serverProperties == nil || serverProperties.IPV4FirewallSettings == nil {
		return utils.Bool(false), schema.NewSet(hashAnalysisServicesServerIpv4FirewallRule, make([]interface{}, 0))
	}

	firewallSettings := serverProperties.IPV4FirewallSettings

	enablePowerBi := utils.Bool(false)
	if firewallSettings.EnablePowerBIService != nil {
		enablePowerBi = firewallSettings.EnablePowerBIService
	}

	fwRules := make([]interface{}, 0)
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

	return enablePowerBi, schema.NewSet(hashAnalysisServicesServerIpv4FirewallRule, fwRules)
}

func hashAnalysisServicesServerIpv4FirewallRule(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["name"].(string))))
	buf.WriteString(fmt.Sprintf("%s-", m["range_start"].(string)))
	buf.WriteString(m["range_end"].(string))

	return hashcode.String(buf.String())
}
