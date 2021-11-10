package analysisservices

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/analysisservices/sdk/2017-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/analysisservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAnalysisServicesServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAnalysisServicesServerCreate,
		Read:   resourceAnalysisServicesServerRead,
		Update: resourceAnalysisServicesServerUpdate,
		Delete: resourceAnalysisServicesServerDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := servers.ParseServerID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku": {
				Type:     pluginsdk.TypeString,
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
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"enable_power_bi_service": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"ipv4_firewall_rule": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"range_start": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azValidate.IPv4Address,
						},
						"range_end": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azValidate.IPv4Address,
						},
					},
				},
				Set: hashAnalysisServicesServerIpv4FirewallRule,
			},

			"querypool_connection_mode": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.QueryPoolConnectionMode(),
			},

			"backup_blob_container_uri": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"server_full_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceAnalysisServicesServerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AnalysisServices.ServerClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Analysis Services Server creation.")

	id := servers.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		server, err := client.GetDetails(ctx, id)
		if err != nil {
			if !response.WasNotFound(server.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(server.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_analysis_services_server", id.ID())
		}
	}

	serverProperties := expandAnalysisServicesServerProperties(d)
	analysisServicesServer := servers.AnalysisServicesServer{
		Location: azure.NormalizeLocation(d.Get("location").(string)),
		Sku: servers.ResourceSku{
			Name: d.Get("sku").(string),
		},
		Properties: serverProperties,
		Tags:       expandTags(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateThenPoll(ctx, id, analysisServicesServer); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceAnalysisServicesServerRead(d, meta)
}

func resourceAnalysisServicesServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AnalysisServices.ServerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseServerID(d.Id())
	if err != nil {
		return err
	}

	server, err := client.GetDetails(ctx, *id)
	if err != nil {
		if response.WasNotFound(server.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := server.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("sku", model.Sku.Name)

		if props := model.Properties; props != nil {
			adminUsers := make([]string, 0)
			if props.AsAdministrators != nil && props.AsAdministrators.Members != nil {
				adminUsers = *props.AsAdministrators.Members
			}
			d.Set("admin_users", adminUsers)

			enablePowerBi, fwRules := flattenAnalysisServicesServerFirewallSettings(props)
			d.Set("enable_power_bi_service", enablePowerBi)
			if err := d.Set("ipv4_firewall_rule", fwRules); err != nil {
				return fmt.Errorf("setting `ipv4_firewall_rule`: %s", err)
			}

			connectionMode := ""
			if props.QuerypoolConnectionMode != nil {
				connectionMode = string(*props.QuerypoolConnectionMode)
			}
			d.Set("querypool_connection_mode", connectionMode)
			d.Set("server_full_name", props.ServerFullName)

			if containerUri, ok := d.GetOk("backup_blob_container_uri"); ok {
				d.Set("backup_blob_container_uri", containerUri)
			}
		}

		if err := tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}

	return nil
}

func resourceAnalysisServicesServerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AnalysisServices.ServerClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Analysis Services Server update.")

	id, err := servers.ParseServerID(d.Id())
	if err != nil {
		return err
	}

	serverResp, err := client.GetDetails(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if serverResp.Model == nil {
		return fmt.Errorf("retrieving %s: response model was nil", *id)
	}
	if serverResp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", *id)
	}
	if serverResp.Model.Properties.State == nil {
		return fmt.Errorf("retrieving %s: state was nil", *id)
	}
	state := *serverResp.Model.Properties.State
	if state != servers.StateSucceeded && state != servers.StatePaused {
		return fmt.Errorf("updating %s: State must be either Succeeded or Paused but got %q", *id, string(state))
	}
	isPaused := state == servers.StatePaused

	if isPaused {
		if err := client.ResumeThenPoll(ctx, *id); err != nil {
			return fmt.Errorf("starting %s: %+v", *id, err)
		}
	}

	serverProperties := expandAnalysisServicesServerMutableProperties(d)
	sku := d.Get("sku").(string)
	t := d.Get("tags").(map[string]interface{})

	analysisServicesServer := servers.AnalysisServicesServerUpdateParameters{
		Sku: &servers.ResourceSku{
			Name: sku,
		},
		Tags:       expandTags(t),
		Properties: serverProperties,
	}

	if err := client.UpdateThenPoll(ctx, *id, analysisServicesServer); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if isPaused {
		if err := client.SuspendThenPoll(ctx, *id); err != nil {
			return fmt.Errorf("re-pausing %s: %+v", *id, err)
		}
	}

	return resourceAnalysisServicesServerRead(d, meta)
}

func resourceAnalysisServicesServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AnalysisServices.ServerClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseServerID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandAnalysisServicesServerProperties(d *pluginsdk.ResourceData) *servers.AnalysisServicesServerProperties {
	adminUsers := expandAnalysisServicesServerAdminUsers(d)

	serverProperties := servers.AnalysisServicesServerProperties{
		AsAdministrators: adminUsers,
	}

	serverProperties.IpV4FirewallSettings = expandAnalysisServicesServerFirewallSettings(d)

	if querypoolConnectionMode, ok := d.GetOk("querypool_connection_mode"); ok {
		connectionMode := servers.ConnectionMode(querypoolConnectionMode.(string))
		serverProperties.QuerypoolConnectionMode = &connectionMode
	}

	if containerUri, ok := d.GetOk("backup_blob_container_uri"); ok {
		serverProperties.BackupBlobContainerUri = utils.String(containerUri.(string))
	}

	return &serverProperties
}

func expandAnalysisServicesServerMutableProperties(d *pluginsdk.ResourceData) *servers.AnalysisServicesServerMutableProperties {
	adminUsers := expandAnalysisServicesServerAdminUsers(d)

	serverProperties := servers.AnalysisServicesServerMutableProperties{
		AsAdministrators: adminUsers,
	}

	serverProperties.IpV4FirewallSettings = expandAnalysisServicesServerFirewallSettings(d)

	connectionMode := servers.ConnectionMode(d.Get("querypool_connection_mode").(string))
	serverProperties.QuerypoolConnectionMode = &connectionMode

	if containerUri, ok := d.GetOk("backup_blob_container_uri"); ok {
		serverProperties.BackupBlobContainerUri = utils.String(containerUri.(string))
	}

	return &serverProperties
}

func expandAnalysisServicesServerAdminUsers(d *pluginsdk.ResourceData) *servers.ServerAdministrators {
	adminUsers := d.Get("admin_users").(*pluginsdk.Set)
	members := make([]string, 0)

	for _, admin := range adminUsers.List() {
		if adm, ok := admin.(string); ok {
			members = append(members, adm)
		}
	}

	return &servers.ServerAdministrators{
		Members: &members,
	}
}

func expandAnalysisServicesServerFirewallSettings(d *pluginsdk.ResourceData) *servers.IPv4FirewallSettings {
	firewallSettings := servers.IPv4FirewallSettings{
		EnablePowerBIService: utils.Bool(d.Get("enable_power_bi_service").(bool)),
	}

	firewallRules := d.Get("ipv4_firewall_rule").(*pluginsdk.Set).List()

	fwRules := make([]servers.IPv4FirewallRule, len(firewallRules))
	for i, v := range firewallRules {
		fwRule := v.(map[string]interface{})
		fwRules[i] = servers.IPv4FirewallRule{
			FirewallRuleName: utils.String(fwRule["name"].(string)),
			RangeStart:       utils.String(fwRule["range_start"].(string)),
			RangeEnd:         utils.String(fwRule["range_end"].(string)),
		}
	}
	firewallSettings.FirewallRules = &fwRules

	return &firewallSettings
}

func flattenAnalysisServicesServerFirewallSettings(serverProperties *servers.AnalysisServicesServerProperties) (*bool, *pluginsdk.Set) {
	if serverProperties == nil || serverProperties.IpV4FirewallSettings == nil {
		return utils.Bool(false), pluginsdk.NewSet(hashAnalysisServicesServerIpv4FirewallRule, make([]interface{}, 0))
	}

	firewallSettings := serverProperties.IpV4FirewallSettings

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

	return enablePowerBi, pluginsdk.NewSet(hashAnalysisServicesServerIpv4FirewallRule, fwRules)
}

func hashAnalysisServicesServerIpv4FirewallRule(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["name"].(string))))
	buf.WriteString(fmt.Sprintf("%s-", m["range_start"].(string)))
	buf.WriteString(m["range_end"].(string))

	return pluginsdk.HashString(buf.String())
}
