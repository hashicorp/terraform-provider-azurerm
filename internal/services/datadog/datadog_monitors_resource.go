package datadog

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datadog/mgmt/2021-03-01/datadog" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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

	id := parse.NewDatadogMonitorID(subscriptionId, resourceGroup, name)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_datadog_monitor", id.ID())
	}

	monitoringStatus := datadog.MonitoringStatusDisabled
	if d.Get("monitoring_enabled").(bool) {
		monitoringStatus = datadog.MonitoringStatusEnabled
	}

	body := datadog.MonitorResource{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Identity: expandMonitorIdentityProperties(d.Get("identity").([]interface{})),
		Sku: &datadog.ResourceSku{
			Name: utils.String(d.Get("sku_name").(string)),
		},
		Properties: &datadog.MonitorProperties{
			DatadogOrganizationProperties: expandMonitorOrganizationProperties(d.Get("datadog_organization").([]interface{})),
			UserInfo:                      expandMonitorUserInfo(d.Get("user").([]interface{})),
			MonitoringStatus:              monitoringStatus,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	future, err := client.Create(ctx, resourceGroup, name, &body)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
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
			log.Printf("[INFO] Datadog monitor %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.MonitorName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if err := d.Set("identity", flattenMonitorIdentityProperties(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	if props := resp.Properties; props != nil {
		if err := d.Set("datadog_organization", flattenMonitorOrganizationProperties(props.DatadogOrganizationProperties, d)); err != nil {
			return fmt.Errorf("setting `datadog_organization`: %+v", err)
		}
		d.Set("monitoring_enabled", props.MonitoringStatus == datadog.MonitoringStatusEnabled)
		d.Set("marketplace_subscription_status", props.MarketplaceSubscriptionStatus)
	}
	if resp.Sku != nil {
		d.Set("sku_name", *resp.Sku.Name)
	}

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
	if d.HasChange("sku_name") {
		body.Sku = &datadog.ResourceSku{Name: utils.String(d.Get("sku_name").(string))}
	}
	if d.HasChange("monitoring_enabled") {
		monitoringStatus := datadog.MonitoringStatusDisabled
		if d.Get("monitoring_enabled").(bool) {
			monitoringStatus = datadog.MonitoringStatusEnabled
		}
		body.Properties.MonitoringStatus = monitoringStatus
	}
	if d.HasChange("tags") {
		body.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.MonitorName, &body)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for updating of %s: %+v", id, err)
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
		return fmt.Errorf("deleting of %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
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

func expandMonitorIdentityProperties(input []interface{}) *datadog.IdentityProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &datadog.IdentityProperties{
		Type: datadog.ManagedIdentityTypes(v["type"].(string)),
	}
}

func expandMonitorOrganizationProperties(input []interface{}) *datadog.OrganizationProperties {
	if len(input) == 0 || input[0] == nil {
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
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &datadog.UserInfo{
		Name:         utils.String(v["name"].(string)),
		EmailAddress: utils.String(v["email"].(string)),
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

func flattenMonitorOrganizationProperties(input *datadog.OrganizationProperties, d *pluginsdk.ResourceData) []interface{} {
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
	if input.ID != nil {
		id = *input.ID
	}
	var redirectUri string
	if input.RedirectURI != nil {
		redirectUri = *input.RedirectURI
	}
	var enterpriseAppID string
	if input.EnterpriseAppID != nil {
		enterpriseAppID = *input.EnterpriseAppID
	}
	return []interface{}{
		map[string]interface{}{
			"name":              name,
			"api_key":           utils.String(v["api_key"].(string)),
			"application_key":   utils.String(v["application_key"].(string)),
			"enterprise_app_id": enterpriseAppID,
			"linking_auth_code": utils.String(v["linking_auth_code"].(string)),
			"linking_client_id": utils.String(v["linking_client_id"].(string)),
			"redirect_uri":      redirectUri,
			"id":                id,
		},
	}
}
