package datadog

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datadog/mgmt/2021-03-01/datadog"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
			_, err := parse.DatadogMonitorID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatadogMonitorsName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"datadog_organization_properties": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"api_key": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"application_key": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"enterprise_app_id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"linking_auth_code": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"linking_client_id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"redirect_uri": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"identity": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
								"UserAssigned",
							}, false),
						},

						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"sku": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"user_info": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.DatadogUsersName,
						},

						"email_address": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.DatadogMonitorsEmailAddress,
						},

						"phone_number": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.DatadogMonitorsPhoneNumber,
						},
					},
				},
			},

			"monitoring_status": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"liftr_resource_category": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"liftr_resource_preference": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"marketplace_subscription_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceDatadogMonitorCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Datadog.MonitorsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewDatadogMonitorID(subscriptionId, resourceGroup, name).ID()

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing Datadog Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_datadog_monitor", id)
	}

	monitoringStatus := datadog.MonitoringStatusDisabled
	if d.Get("monitoring_status").(bool) {
		monitoringStatus = datadog.MonitoringStatusEnabled
	}

	body := datadog.MonitorResource{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Identity: expandMonitorIdentityProperties(d.Get("identity").([]interface{})),
		Sku:      expandMonitorResourceSku(d.Get("sku").([]interface{})),
		Properties: &datadog.MonitorProperties{
			DatadogOrganizationProperties: expandMonitorOrganizationProperties(d.Get("datadog_organization_properties").([]interface{})),
			UserInfo:                      expandMonitorUserInfo(d.Get("user_info").([]interface{})),
			MonitoringStatus:              monitoringStatus,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	future, err := client.Create(ctx, resourceGroup, name, &body)
	if err != nil {
		return fmt.Errorf("creating Datadog Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of the Datadog Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id)
	return resourceDatadogMonitorRead(d, meta)
}

func resourceDatadogMonitorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.MonitorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatadogMonitorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MonitorName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] datadog %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Datadog Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}
	d.Set("name", id.MonitorName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if err := d.Set("identity", flattenMonitorIdentityProperties(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	if props := resp.Properties; props != nil {
		if err := d.Set("datadog_organization_properties", flattenMonitorOrganizationProperties(props.DatadogOrganizationProperties)); err != nil {
			return fmt.Errorf("setting `datadog_organization_properties`: %+v", err)
		}
		d.Set("monitoring_status", props.MonitoringStatus == datadog.MonitoringStatusEnabled)
		if err := d.Set("user_info", flattenMonitorUserInfo(props.UserInfo)); err != nil {
			return fmt.Errorf("setting `user_info`: %+v", err)
		}
		d.Set("liftr_resource_category", props.LiftrResourceCategory)
		d.Set("liftr_resource_preference", props.LiftrResourcePreference)
		d.Set("marketplace_subscription_status", props.MarketplaceSubscriptionStatus)
	}
	if err := d.Set("sku", flattenMonitorResourceSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `sku`: %+v", err)
	}
	d.Set("type", resp.Type)
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDatadogMonitorUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.MonitorsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatadogMonitorID(d.Id())
	if err != nil {
		return err
	}

	body := datadog.MonitorResourceUpdateParameters{
		Properties: &datadog.MonitorUpdateProperties{},
	}
	if d.HasChange("monitoring_status") {
		monitoringStatus := datadog.MonitoringStatusDisabled
		if d.Get("monitoring_status").(bool) {
			monitoringStatus = datadog.MonitoringStatusEnabled
		}
		body.Properties.MonitoringStatus = monitoringStatus
	}
	if d.HasChange("tags") {
		body.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.MonitorName, &body); err != nil {
		return fmt.Errorf("updating Datadog Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}
	return resourceDatadogMonitorRead(d, meta)
}

func resourceDatadogMonitorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datadog.MonitorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatadogMonitorID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.MonitorName)
	if err != nil {
		return fmt.Errorf("deleting Datadog Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the Datadog Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}
	return nil
}

func expandMonitorIdentityProperties(input []interface{}) *datadog.IdentityProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &datadog.IdentityProperties{
		Type: datadog.ManagedIdentityTypes(v["type"].(string)),
	}
}

func expandMonitorResourceSku(input []interface{}) *datadog.ResourceSku {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &datadog.ResourceSku{
		Name: utils.String(v["name"].(string)),
	}
}

func expandMonitorOrganizationProperties(input []interface{}) *datadog.OrganizationProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &datadog.OrganizationProperties{
		LinkingAuthCode: utils.String(v["linking_auth_code"].(string)),
		LinkingClientID: utils.String(v["linking_client_id"].(string)),
		RedirectURI:     utils.String(v["redirect_uri"].(string)),
		APIKey:          utils.String(v["api_key"].(string)),
		ApplicationKey:  utils.String(v["application_key"].(string)),
		EnterpriseAppID: utils.String(v["enterprise_app_id"].(string)),
	}
}

func expandMonitorUserInfo(input []interface{}) *datadog.UserInfo {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &datadog.UserInfo{
		Name:         utils.String(v["name"].(string)),
		EmailAddress: utils.String(v["email_address"].(string)),
		PhoneNumber:  utils.String(v["phone_number"].(string)),
	}
}

func flattenMonitorIdentityProperties(input *datadog.IdentityProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var t datadog.ManagedIdentityTypes
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

func flattenMonitorOrganizationProperties(input *datadog.OrganizationProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var name string
	if input.Name != nil {
		name = *input.Name
	}
	var apiKey string
	if input.APIKey != nil {
		apiKey = *input.APIKey
	}
	var applicationKey string
	if input.ApplicationKey != nil {
		applicationKey = *input.ApplicationKey
	}
	var enterpriseAppId string
	if input.EnterpriseAppID != nil {
		enterpriseAppId = *input.EnterpriseAppID
	}
	var linkingAuthCode string
	if input.LinkingAuthCode != nil {
		linkingAuthCode = *input.LinkingAuthCode
	}
	var linkingClientId string
	if input.LinkingClientID != nil {
		linkingClientId = *input.LinkingClientID
	}
	var redirectUri string
	if input.RedirectURI != nil {
		redirectUri = *input.RedirectURI
	}
	var id string
	if input.ID != nil {
		id = *input.ID
	}
	return []interface{}{
		map[string]interface{}{
			"name":              name,
			"api_key":           apiKey,
			"application_key":   applicationKey,
			"enterprise_app_id": enterpriseAppId,
			"linking_auth_code": linkingAuthCode,
			"linking_client_id": linkingClientId,
			"redirect_uri":      redirectUri,
			"id":                id,
		},
	}
}

func flattenMonitorUserInfo(input *datadog.UserInfo) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var name string
	if input.Name != nil {
		name = *input.Name
	}
	var emailAddress string
	if input.EmailAddress != nil {
		emailAddress = *input.EmailAddress
	}
	var phoneNumber string
	if input.PhoneNumber != nil {
		phoneNumber = *input.PhoneNumber
	}
	return []interface{}{
		map[string]interface{}{
			"name":          name,
			"email_address": emailAddress,
			"phone_number":  phoneNumber,
		},
	}
}

func flattenMonitorResourceSku(input *datadog.ResourceSku) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var name string
	if input.Name != nil {
		name = *input.Name
	}
	return []interface{}{
		map[string]interface{}{
			"name": name,
		},
	}
}
