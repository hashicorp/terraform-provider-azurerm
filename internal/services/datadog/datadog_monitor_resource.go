// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datadog

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/monitorsresource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDatadogMonitor() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDatadogMonitorCreate,
		Read:   resourceDatadogMonitorRead,
		Update: resourceDatadogMonitorUpdate,
		Delete: resourceDatadogMonitorDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := monitorsresource.ParseMonitorID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatadogMonitorsName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"datadog_organization": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"api_key": {
							Type:      pluginsdk.TypeString,
							Required:  true,
							ForceNew:  true,
							Sensitive: true,
						},

						"application_key": {
							Type:      pluginsdk.TypeString,
							Required:  true,
							ForceNew:  true,
							Sensitive: true,
						},

						"enterprise_app_id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"linking_auth_code": {
							Type:      pluginsdk.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},

						"linking_client_id": {
							Type:      pluginsdk.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},

						"redirect_uri": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"sku_name": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				DiffSuppressFunc: SkuNameDiffSuppress,
			},

			"user": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DatadogUsersName,
						},

						"email": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.DatadogMonitorsEmailAddress,
						},

						"phone_number": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validate.DatadogMonitorsPhoneNumber,
						},
					},
				},
			},

			"monitoring_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"marketplace_subscription_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceDatadogMonitorCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Datadog.MonitorsResource
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := monitorsresource.NewMonitorID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.MonitorsGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_datadog_monitor", id.ID())
	}

	monitoringStatus := monitorsresource.MonitoringStatusDisabled
	if d.Get("monitoring_enabled").(bool) {
		monitoringStatus = monitorsresource.MonitoringStatusEnabled
	}

	payload := monitorsresource.DatadogMonitorResource{
		Location: location.Normalize(d.Get("location").(string)),
		Identity: expandMonitorIdentityProperties(d.Get("identity").([]interface{})),
		Sku: &monitorsresource.ResourceSku{
			Name: d.Get("sku_name").(string),
		},
		Properties: &monitorsresource.MonitorProperties{
			DatadogOrganizationProperties: expandMonitorOrganizationProperties(d.Get("datadog_organization").([]interface{})),
			UserInfo:                      expandMonitorUserInfo(d.Get("user").([]interface{})),
			MonitoringStatus:              pointer.To(monitoringStatus),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.MonitorsCreateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDatadogMonitorRead(d, meta)
}

func resourceDatadogMonitorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.MonitorsResource
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := monitorsresource.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.MonitorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Datadog monitor %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.MonitorName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if err := d.Set("identity", flattenMonitorIdentityProperties(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
		if props := model.Properties; props != nil {
			if err := d.Set("datadog_organization", flattenMonitorOrganizationProperties(props.DatadogOrganizationProperties, d)); err != nil {
				return fmt.Errorf("setting `datadog_organization`: %+v", err)
			}

			monitoringEnabled := false
			if props.MonitoringStatus != nil && *props.MonitoringStatus == monitorsresource.MonitoringStatusEnabled {
				monitoringEnabled = true
			}

			d.Set("monitoring_enabled", monitoringEnabled)
			d.Set("marketplace_subscription_status", string(pointer.From(props.MarketplaceSubscriptionStatus)))
		}

		skuName := ""
		if model.Sku != nil {
			skuName = model.Sku.Name
		}
		d.Set("sku_name", skuName)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceDatadogMonitorUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.MonitorsResource
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := monitorsresource.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	payload := monitorsresource.DatadogMonitorResourceUpdateParameters{
		Properties: &monitorsresource.MonitorUpdateProperties{},
	}
	if d.HasChange("sku_name") {
		payload.Sku = &monitorsresource.ResourceSku{
			Name: d.Get("sku_name").(string),
		}
	}
	if d.HasChange("monitoring_enabled") {
		monitoringStatus := monitorsresource.MonitoringStatusDisabled
		if d.Get("monitoring_enabled").(bool) {
			monitoringStatus = monitorsresource.MonitoringStatusEnabled
		}
		payload.Properties.MonitoringStatus = pointer.To(monitoringStatus)
	}
	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.MonitorsUpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}
	return resourceDatadogMonitorRead(d, meta)
}

func resourceDatadogMonitorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.MonitorsResource
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := monitorsresource.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	if err := client.MonitorsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting of %s: %+v", *id, err)
	}

	return nil
}

func SkuNameDiffSuppress(_, old, new string, _ *pluginsdk.ResourceData) bool {
	// During creating, accepts any sku name.
	if old == "" {
		return false
	}
	// Sku name of the datadog monitor has two kinds:
	// - Concrete sku: E.g. "payg_v2_Monthly". These will be returned unchanged by API.
	// - Linked sku:  The value is named "Linked". This will be accepted and changed by the API, which will then return you the linked concrete sku.
	if new == "Linked" {
		return true
	}
	return old == new
}

func expandMonitorIdentityProperties(input []interface{}) *monitorsresource.IdentityProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &monitorsresource.IdentityProperties{
		// @tombuildsstuff: this should be normalized in Pandora to a common identity type? is this a Swagger bag with SA & UA omitted?
		Type: pointer.To(monitorsresource.ManagedIdentityTypes(v["type"].(string))),
	}
}

func expandMonitorOrganizationProperties(input []interface{}) *monitorsresource.DatadogOrganizationProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &monitorsresource.DatadogOrganizationProperties{
		LinkingAuthCode: utils.String(v["linking_auth_code"].(string)),
		LinkingClientId: utils.String(v["linking_client_id"].(string)),
		RedirectUri:     utils.String(v["redirect_uri"].(string)),
		ApiKey:          utils.String(v["api_key"].(string)),
		ApplicationKey:  utils.String(v["application_key"].(string)),
		EnterpriseAppId: utils.String(v["enterprise_app_id"].(string)),
	}
}

func expandMonitorUserInfo(input []interface{}) *monitorsresource.UserInfo {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &monitorsresource.UserInfo{
		Name:         utils.String(v["name"].(string)),
		EmailAddress: utils.String(v["email"].(string)),
		PhoneNumber:  utils.String(v["phone_number"].(string)),
	}
}

func flattenMonitorIdentityProperties(input *monitorsresource.IdentityProperties) []interface{} {
	if input == nil || input.Type == nil {
		return make([]interface{}, 0)
	}

	var t string
	if *input.Type != "" {
		t = string(*input.Type)
	}
	var principalId string
	if input.PrincipalId != nil {
		principalId = *input.PrincipalId
	}
	var tenantId string
	if input.TenantId != nil {
		tenantId = *input.TenantId
	}
	return []interface{}{
		map[string]interface{}{
			"type":         t,
			"principal_id": principalId,
			"tenant_id":    tenantId,
		},
	}
}

func flattenMonitorOrganizationProperties(input *monitorsresource.DatadogOrganizationProperties, d *pluginsdk.ResourceData) []interface{} {
	organisationProperties := d.Get("datadog_organization").([]interface{})
	if len(organisationProperties) == 0 {
		return make([]interface{}, 0)
	}
	v := organisationProperties[0].(map[string]interface{})

	var name string
	if input.Name != nil {
		name = *input.Name
	}
	var id string
	if input.Id != nil {
		id = *input.Id
	}
	var redirectUri string
	if input.RedirectUri != nil {
		redirectUri = *input.RedirectUri
	}
	var enterpriseAppId string
	if input.EnterpriseAppId != nil {
		enterpriseAppId = *input.EnterpriseAppId
	}
	return []interface{}{
		map[string]interface{}{
			"name":              name,
			"api_key":           utils.String(v["api_key"].(string)),
			"application_key":   utils.String(v["application_key"].(string)),
			"enterprise_app_id": enterpriseAppId,
			"linking_auth_code": utils.String(v["linking_auth_code"].(string)),
			"linking_client_id": utils.String(v["linking_client_id"].(string)),
			"redirect_uri":      redirectUri,
			"id":                id,
		},
	}
}
