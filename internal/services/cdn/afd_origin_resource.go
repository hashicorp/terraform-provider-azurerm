package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAfdOrigin() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAfdOriginsCreate,
		Read:   resourceAfdOriginsRead,
		Update: resourceAfdOriginsUpdate,
		Delete: resourceAfdOriginsDelete,

		SchemaVersion: 1,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AfdOriginGroupsID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"origin_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AfdOriginGroupsID,
			},

			"host_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"azure_origin": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"origin_host_header": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
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
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

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

			"private_link": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"approval_message": { //   A custom message to be included in the approval request to connect to the Private Link.
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"location": { // The location of the Private Link resource. Required only if 'privateLinkResourceId' is populated
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"resource_id": { // The Resource Id of the Private Link resource. Populating this optional field indicates that this backend is 'Private'
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},
		},
	}
}

func resourceAfdOriginsCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// parse origin_group_id
	originGroupId := d.Get("origin_group_id").(string)
	originGroup, err := parse.AfdOriginGroupsID(originGroupId)
	if err != nil {
		return err
	}

	originname := d.Get("name").(string)
	hostname := d.Get("host_name").(string)
	priority := int32(d.Get("priority").(int))
	weight := int32(d.Get("weight").(int))
	httpPort := int32(d.Get("http_port").(int))
	httpsPort := int32(d.Get("https_port").(int))
	originHostHeader := d.Get("origin_host_header").(string)
	azureOrigin := d.Get("azure_origin").(string)
	privateLinkSettings := d.Get("private_link").([]interface{})
	enabledState := d.Get("enabled").(bool)

	id := parse.NewAfdOriginsID(originGroup.SubscriptionId, originGroup.ResourceGroup, originGroup.ProfileName, originGroup.OriginGroupName, originname)

	afdOrigin := cdn.AFDOrigin{
		AFDOriginProperties: &cdn.AFDOriginProperties{
			HostName:  &hostname,
			Priority:  &priority,
			Weight:    &weight,
			HTTPPort:  &httpPort,
			HTTPSPort: &httpsPort,
		},
	}

	if enabledState {
		afdOrigin.EnabledState = cdn.EnabledStateEnabled
	} else {
		afdOrigin.EnabledState = cdn.EnabledStateDisabled
	}

	if azureOrigin != "" {
		afdOrigin.AzureOrigin = &cdn.ResourceReference{
			ID: &azureOrigin,
		}
	}

	if originHostHeader != "" {
		afdOrigin.OriginHostHeader = &originHostHeader
	}

	if privateLinkSettings != nil {
		afdOrigin.SharedPrivateLinkResource = expandPrivateLinkSettings(privateLinkSettings)
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, originGroup.OriginGroupName, originname, afdOrigin)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceAfdOriginsRead(d, meta)
}

func expandPrivateLinkSettings(input []interface{}) interface{} {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})

	resourceId := config["resource_id"].(string)
	location := config["location"].(string)
	approvalMessage := config["approval_message"].(string)

	privateLinkResource := cdn.SharedPrivateLinkResourceProperties{
		PrivateLink: &cdn.ResourceReference{
			ID: &resourceId,
		},
		PrivateLinkLocation: &location,
		RequestMessage:      &approvalMessage,
	}

	return privateLinkResource
}

func resourceAfdOriginsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdOriginsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.OriginName)
	d.Set("http_port", resp.HTTPPort)
	d.Set("https_port", resp.HTTPSPort)
	d.Set("priority", resp.Priority)
	d.Set("weight", resp.Weight)
	d.Set("host_name", resp.HostName)
	d.Set("origin_host_header", resp.OriginHostHeader)

	d.Set("origin_group_id", parse.NewAfdOriginGroupsID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.OriginGroupName).ID())

	if resp.EnabledState == cdn.EnabledStateEnabled {
		d.Set("enabled", true)
	} else {
		d.Set("enabled", false)
	}

	return nil
}

func resourceAfdOriginsUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdOriginsID(d.Id())
	if err != nil {
		return err
	}

	var originUpdate cdn.AFDOriginUpdateParameters
	var originUpdateProperties cdn.AFDOriginUpdatePropertiesParameters

	if d.HasChange("origin_host_header") {
		originHostHeader := d.Get("origin_host_header").(string)
		originUpdateProperties.OriginHostHeader = &originHostHeader
	}

	if d.HasChange("enabled") {
		enabledState := d.Get("enabled").(bool)
		if enabledState {
			originUpdateProperties.EnabledState = cdn.EnabledStateEnabled
		} else {
			originUpdateProperties.EnabledState = cdn.EnabledStateDisabled
		}
	}

	if d.HasChange("http_port") {
		httpPort := int32(d.Get("http_port").(int))
		originUpdateProperties.HTTPPort = &httpPort
	}

	if d.HasChange("https_port") {
		httpsPort := int32(d.Get("https_port").(int))
		originUpdateProperties.HTTPSPort = &httpsPort
	}

	if d.HasChange("priority") {
		priority := int32(d.Get("priority").(int))
		originUpdateProperties.Priority = &priority
	}

	if d.HasChange("private_link") {
		privateLink := d.Get("private_link").([]interface{})
		originUpdateProperties.SharedPrivateLinkResource = expandPrivateLinkSettings(privateLink)
	}

	if d.HasChange("weight") {
		weight := int32(d.Get("weight").(int))
		originUpdateProperties.Weight = &weight
	}

	originUpdate.AFDOriginUpdatePropertiesParameters = &originUpdateProperties

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName, originUpdate)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceAfdOriginsRead(d, meta)
}

func resourceAfdOriginsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDOriginsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdOriginsID(d.Id())
	if err != nil {
		return err
	}

	// TODO - Cannot delete the last origin when the origin group is still associated with a route under an endpoint. Please disassociate the origin group and try again.

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName)

	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}
