package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdorigingroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdorigins"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontdoorOrigin() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontdoorOriginCreate,
		Read:   resourceFrontdoorOriginRead,
		Update: resourceFrontdoorOriginUpdate,
		Delete: resourceFrontdoorOriginDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := afdorigins.ParseOriginGroupOriginID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"frontdoor_origin_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: afdorigingroups.ValidateOriginGroupID,
			},

			// HostName cannot be null or empty.;
			"host_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"azure_origin_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"deployment_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"enable_health_probes": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enforce_certificate_name_check": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			// Property 'AfdOrigin.Priority' cannot be set to '10000000'. Acceptable values are within range [1, 5];
			// Property 'AfdOrigin.Weight' cannot be set to '10000000'. Acceptable values are within range [1, 1000]"
			"http_port": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      80,
				ValidateFunc: validation.IntBetween(1, 65535),
			},

			"https_port": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      443,
				ValidateFunc: validation.IntBetween(1, 65535),
			},

			// Must be a valid domain name, IP version 4, or IP version 6
			"origin_host_header": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: IsValidDomain,
			},

			"priority": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 5),
			},

			"weight": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      500,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			"origin_group_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFrontdoorOriginCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	originGroupId, err := parse.FrontdoorOriginGroupID(d.Get("frontdoor_origin_group_id").(string))
	if err != nil {
		return err
	}

	id := afdorigins.NewOriginGroupOriginID(originGroupId.SubscriptionId, originGroupId.ResourceGroup, originGroupId.ProfileName, originGroupId.OriginGroupName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_frontdoor_origin", id.ID())
		}
	}

	props := afdorigins.AFDOrigin{
		Properties: &afdorigins.AFDOriginProperties{
			AzureOrigin:                 expandOriginGroupOriginResourceReference(d.Get("azure_origin_id").(string)),
			EnabledState:                ConvertBoolToOriginsEnabledState(d.Get("enable_health_probes").(bool)),
			EnforceCertificateNameCheck: utils.Bool(d.Get("enforce_certificate_name_check").(bool)),
			HostName:                    d.Get("host_name").(string),
			HttpPort:                    utils.Int64(int64(d.Get("http_port").(int))),
			HttpsPort:                   utils.Int64(int64(d.Get("https_port").(int))),
			OriginHostHeader:            utils.String(d.Get("origin_host_header").(string)),
			Priority:                    utils.Int64(int64(d.Get("priority").(int))),
			Weight:                      utils.Int64(int64(d.Get("weight").(int))),
		},
	}

	if err := client.CreateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontdoorOriginRead(d, meta)
}

func resourceFrontdoorOriginRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdorigins.ParseOriginGroupOriginID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.OriginName)

	d.Set("frontdoor_origin_group_id", afdorigingroups.NewOriginGroupID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.OriginGroupName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {

			if err := d.Set("azure_origin_id", flattenOriginGroupOriginResourceReference(props.AzureOrigin)); err != nil {
				return fmt.Errorf("setting `azure_origin_id`: %+v", err)
			}

			d.Set("deployment_status", props.DeploymentStatus)
			d.Set("enable_health_probes", ConvertOriginsEnabledStateToBool(props.EnabledState))
			d.Set("enforce_certificate_name_check", props.EnforceCertificateNameCheck)
			d.Set("host_name", props.HostName)
			d.Set("http_port", props.HttpPort)
			d.Set("https_port", props.HttpsPort)
			d.Set("origin_group_name", props.OriginGroupName)
			d.Set("origin_host_header", props.OriginHostHeader)
			d.Set("priority", props.Priority)
			d.Set("provisioning_state", props.ProvisioningState)
			d.Set("weight", props.Weight)
		}
	}
	return nil
}

func resourceFrontdoorOriginUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdorigins.ParseOriginGroupOriginID(d.Id())
	if err != nil {
		return err
	}

	props := afdorigins.AFDOriginUpdateParameters{
		Properties: &afdorigins.AFDOriginUpdatePropertiesParameters{
			AzureOrigin:                 expandOriginGroupOriginResourceReference(d.Get("azure_origin_id").(string)),
			EnabledState:                ConvertBoolToOriginsEnabledState(d.Get("enable_health_probes").(bool)),
			EnforceCertificateNameCheck: utils.Bool(d.Get("enforce_certificate_name_check").(bool)),
			HostName:                    utils.String(d.Get("host_name").(string)),
			HttpPort:                    utils.Int64(int64(d.Get("http_port").(int))),
			HttpsPort:                   utils.Int64(int64(d.Get("https_port").(int))),
			OriginHostHeader:            utils.String(d.Get("origin_host_header").(string)),
			Priority:                    utils.Int64(int64(d.Get("priority").(int))),
			Weight:                      utils.Int64(int64(d.Get("weight").(int))),
		},
	}
	if err := client.UpdateThenPoll(ctx, *id, props); err != nil {

		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceFrontdoorOriginRead(d, meta)
}

func resourceFrontdoorOriginDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdorigins.ParseOriginGroupOriginID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {

		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandOriginGroupOriginResourceReference(input string) *afdorigins.ResourceReference {
	if len(input) == 0 {
		return nil
	}

	return &afdorigins.ResourceReference{
		Id: utils.String(input),
	}
}

func flattenOriginGroupOriginResourceReference(input *afdorigins.ResourceReference) string {
	result := ""
	if input == nil {
		return result
	}

	if input.Id != nil {
		result = *input.Id
	}
	return result
}
