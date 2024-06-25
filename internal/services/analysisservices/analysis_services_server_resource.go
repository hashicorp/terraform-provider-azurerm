// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package analysisservices

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/analysisservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAnalysisServicesServer() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
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

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

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

			"power_bi_service_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: !features.FourPointOhBeta(),
				ConflictsWith: func() []string {
					if !features.FourPointOhBeta() {
						return []string{"enable_power_bi_service"}
					}
					return []string{}
				}(),
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
				Set: hashAnalysisServicesServerIPv4FirewallRule,
			},

			"querypool_connection_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(servers.ConnectionModeAll),
				ValidateFunc: validation.StringInSlice([]string{
					string(servers.ConnectionModeAll),
					string(servers.ConnectionModeReadOnly),
				}, false),
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

			"tags": commonschema.Tags(),
		},
	}
	if !features.FourPointOhBeta() {
		resource.Schema["querypool_connection_mode"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(servers.ConnectionModeAll),
				string(servers.ConnectionModeReadOnly),
			}, false),
		}
		resource.Schema["enable_power_bi_service"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			Deprecated:    "The property `enable_power_bi_service` has been superseded by `power_bi_service_enabled` and will be removed in v4.0 of the AzureRM Provider.",
			ConflictsWith: []string{"power_bi_service_enabled"},
		}
	}
	return resource
}

func resourceAnalysisServicesServerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AnalysisServices.Servers
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Analysis Services Server creation.")

	id := servers.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	server, err := client.GetDetails(ctx, id)
	if err != nil {
		if !response.WasNotFound(server.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(server.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_analysis_services_server", id.ID())
	}

	analysisServicesServer := servers.AnalysisServicesServer{
		Location: location.Normalize(d.Get("location").(string)),
		Sku: servers.ResourceSku{
			Name: d.Get("sku").(string),
		},
		Properties: &servers.AnalysisServicesServerProperties{
			AsAdministrators:     expandAnalysisServicesServerAdminUsers(d),
			IPV4FirewallSettings: expandAnalysisServicesServerFirewallSettings(d),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("enable_power_bi_service"); ok && !features.FourPointOhBeta() {
		if analysisServicesServer.Properties.IPV4FirewallSettings == nil {
			analysisServicesServer.Properties.IPV4FirewallSettings = &servers.IPv4FirewallSettings{
				FirewallRules: pointer.To(make([]servers.IPv4FirewallRule, 0)),
			}
		}
		analysisServicesServer.Properties.IPV4FirewallSettings.EnablePowerBIService = pointer.To(v.(bool))
	}
	if v, ok := d.GetOk("power_bi_service_enabled"); ok {
		if analysisServicesServer.Properties.IPV4FirewallSettings == nil {
			analysisServicesServer.Properties.IPV4FirewallSettings = &servers.IPv4FirewallSettings{
				FirewallRules: pointer.To(make([]servers.IPv4FirewallRule, 0)),
			}
		}
		analysisServicesServer.Properties.IPV4FirewallSettings.EnablePowerBIService = pointer.To(v.(bool))
	}

	if querypoolConnectionMode, ok := d.GetOk("querypool_connection_mode"); ok {
		analysisServicesServer.Properties.QuerypoolConnectionMode = pointer.To(servers.ConnectionMode(querypoolConnectionMode.(string)))
	}

	if containerUri, ok := d.GetOk("backup_blob_container_uri"); ok {
		analysisServicesServer.Properties.BackupBlobContainerUri = pointer.To(containerUri.(string))
	}

	if err := client.CreateThenPoll(ctx, id, analysisServicesServer); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceAnalysisServicesServerRead(d, meta)
}

func resourceAnalysisServicesServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AnalysisServices.Servers
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

	d.Set("name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

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
			if !features.FourPointOhBeta() {
				d.Set("enable_power_bi_service", enablePowerBi)
			}
			d.Set("power_bi_service_enabled", enablePowerBi)
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

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceAnalysisServicesServerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AnalysisServices.Servers
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

	analysisServicesServer := servers.AnalysisServicesServerUpdateParameters{
		Sku: &servers.ResourceSku{
			Name: d.Get("sku").(string),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &servers.AnalysisServicesServerMutableProperties{
			AsAdministrators:        expandAnalysisServicesServerAdminUsers(d),
			IPV4FirewallSettings:    expandAnalysisServicesServerFirewallSettings(d),
			QuerypoolConnectionMode: pointer.To(servers.ConnectionMode(d.Get("querypool_connection_mode").(string))),
		},
	}

	if d.HasChange("power_bi_service_enabled") {
		if analysisServicesServer.Properties.IPV4FirewallSettings == nil {
			analysisServicesServer.Properties.IPV4FirewallSettings = &servers.IPv4FirewallSettings{
				FirewallRules: pointer.To(make([]servers.IPv4FirewallRule, 0)),
			}
		}
		analysisServicesServer.Properties.IPV4FirewallSettings.EnablePowerBIService = pointer.To(d.Get("power_bi_service_enabled").(bool))
	}

	if !features.FourPointOhBeta() {
		if d.HasChange("enable_power_bi_service") {
			if analysisServicesServer.Properties.IPV4FirewallSettings == nil {
				analysisServicesServer.Properties.IPV4FirewallSettings = &servers.IPv4FirewallSettings{
					FirewallRules: pointer.To(make([]servers.IPv4FirewallRule, 0)),
				}
			}
			analysisServicesServer.Properties.IPV4FirewallSettings.EnablePowerBIService = pointer.To(d.Get("enable_power_bi_service").(bool))
		}
	}

	if containerUri, ok := d.GetOk("backup_blob_container_uri"); ok {
		analysisServicesServer.Properties.BackupBlobContainerUri = pointer.To(containerUri.(string))
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
	client := meta.(*clients.Client).AnalysisServices.Servers
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
	firewallRules := d.Get("ipv4_firewall_rule").(*pluginsdk.Set).List()

	if len(firewallRules) == 0 {
		return nil
	}

	firewallSettings := servers.IPv4FirewallSettings{}
	fwRules := make([]servers.IPv4FirewallRule, len(firewallRules))
	for i, v := range firewallRules {
		fwRule := v.(map[string]interface{})
		fwRules[i] = servers.IPv4FirewallRule{
			FirewallRuleName: pointer.To(fwRule["name"].(string)),
			RangeStart:       pointer.To(fwRule["range_start"].(string)),
			RangeEnd:         pointer.To(fwRule["range_end"].(string)),
		}
	}
	firewallSettings.FirewallRules = &fwRules

	return &firewallSettings
}

func flattenAnalysisServicesServerFirewallSettings(serverProperties *servers.AnalysisServicesServerProperties) (*bool, *pluginsdk.Set) {
	if serverProperties == nil || serverProperties.IPV4FirewallSettings == nil {
		return pointer.To(false), pluginsdk.NewSet(hashAnalysisServicesServerIPv4FirewallRule, make([]interface{}, 0))
	}

	firewallSettings := serverProperties.IPV4FirewallSettings

	enablePowerBi := pointer.To(false)
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

	return enablePowerBi, pluginsdk.NewSet(hashAnalysisServicesServerIPv4FirewallRule, fwRules)
}

func hashAnalysisServicesServerIPv4FirewallRule(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["name"].(string))))
	buf.WriteString(fmt.Sprintf("%s-", m["range_start"].(string)))
	buf.WriteString(m["range_end"].(string))

	return pluginsdk.HashString(buf.String())
}
