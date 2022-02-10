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
			_, err := parse.FrontdoorOriginID(id)
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

			"azure_origin": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},

			"deployment_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"enabled_state": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"enforce_certificate_name_check": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"host_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"http_port": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},

			"https_port": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},

			"origin_group_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"origin_host_header": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"priority": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},

			"provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"weight": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
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

	sdkId := afdorigins.NewOriginGroupOriginID(originGroupId.SubscriptionId, originGroupId.ResourceGroup, originGroupId.ProfileName, originGroupId.OriginGroupName, d.Get("name").(string))
	id := parse.NewFrontdoorOriginID(originGroupId.SubscriptionId, originGroupId.ResourceGroup, originGroupId.ProfileName, originGroupId.OriginGroupName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, sdkId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_frontdoor_origin", id.ID())
		}
	}

	enabledStateValue := afdorigins.EnabledState(d.Get("enabled_state").(string))
	props := afdorigins.AFDOrigin{
		Properties: &afdorigins.AFDOriginProperties{
			AzureOrigin:                 expandOriginGroupOriginResourceReference(d.Get("azure_origin").([]interface{})),
			EnabledState:                &enabledStateValue,
			EnforceCertificateNameCheck: utils.Bool(d.Get("enforce_certificate_name_check").(bool)),
			HostName:                    d.Get("host_name").(string),
			HttpPort:                    utils.Int64(int64(d.Get("http_port").(int))),
			HttpsPort:                   utils.Int64(int64(d.Get("https_port").(int))),
			OriginHostHeader:            utils.String(d.Get("origin_host_header").(string)),
			Priority:                    utils.Int64(int64(d.Get("priority").(int))),
			Weight:                      utils.Int64(int64(d.Get("weight").(int))),
		},
	}
	if err := client.CreateThenPoll(ctx, sdkId, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontdoorOriginRead(d, meta)
}

func resourceFrontdoorOriginRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sdkId, err := afdorigins.ParseOriginGroupOriginID(d.Id())
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorOriginID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *sdkId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.OriginName)

	d.Set("frontdoor_origin_group_id", afdorigingroups.NewOriginGroupID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.OriginGroupName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {

			if err := d.Set("azure_origin", flattenOriginGroupOriginResourceReference(props.AzureOrigin)); err != nil {
				return fmt.Errorf("setting `azure_origin`: %+v", err)
			}
			d.Set("deployment_status", props.DeploymentStatus)
			d.Set("enabled_state", props.EnabledState)
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

	sdkId, err := afdorigins.ParseOriginGroupOriginID(d.Id())
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorOriginID(d.Id())
	if err != nil {
		return err
	}

	enabledStateValue := afdorigins.EnabledState(d.Get("enabled_state").(string))
	props := afdorigins.AFDOriginUpdateParameters{
		Properties: &afdorigins.AFDOriginUpdatePropertiesParameters{
			AzureOrigin:                 expandOriginGroupOriginResourceReference(d.Get("azure_origin").([]interface{})),
			EnabledState:                &enabledStateValue,
			EnforceCertificateNameCheck: utils.Bool(d.Get("enforce_certificate_name_check").(bool)),
			HostName:                    utils.String(d.Get("host_name").(string)),
			HttpPort:                    utils.Int64(int64(d.Get("http_port").(int))),
			HttpsPort:                   utils.Int64(int64(d.Get("https_port").(int))),
			OriginHostHeader:            utils.String(d.Get("origin_host_header").(string)),
			Priority:                    utils.Int64(int64(d.Get("priority").(int))),
			Weight:                      utils.Int64(int64(d.Get("weight").(int))),
		},
	}
	if err := client.UpdateThenPoll(ctx, *sdkId, props); err != nil {

		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceFrontdoorOriginRead(d, meta)
}

func resourceFrontdoorOriginDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorOriginsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sdkId, err := afdorigins.ParseOriginGroupOriginID(d.Id())
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorOriginID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *sdkId); err != nil {

		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandOriginGroupOriginResourceReference(input []interface{}) *afdorigins.ResourceReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &afdorigins.ResourceReference{
		Id: utils.String(v["id"].(string)),
	}
}

func flattenOriginGroupOriginResourceReference(input *afdorigins.ResourceReference) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	if input.Id != nil {
		result["id"] = *input.Id
	}
	return append(results, result)
}
