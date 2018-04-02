package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNotificationHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNotificationHubCreateUpdate,
		Read:   resourceArmNotificationHubRead,
		Update: resourceArmNotificationHubCreateUpdate,
		Delete: resourceArmNotificationHubDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"apns_credential": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// NOTE: APNS supports two modes, certificate auth (v1) and token auth (v2)
						// certificate authentication/v1 is marked for deprecation; as such we're not
						// supporting it at this time.
						"application_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"application_mode": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Production",
								"Sandbox",
							}, false),
						},
						"application_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"token": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmNotificationHubCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).notificationHubsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	apnsCredential, err := expandNotificationHubsAPNSCredentials(d)
	if err != nil {
		return err
	}

	parameters := notificationhubs.CreateOrUpdateParameters{
		Location: utils.String(location),
		Tags:     expandTags(tags),
		Properties: &notificationhubs.Properties{
			ApnsCredential: apnsCredential,
		},
	}

	_, err = client.CreateOrUpdate(ctx, resourceGroup, namespaceName, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Notification Hub %q (Namespace %q / Resource Group %q) ID", name, namespaceName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmNotificationHubRead(d, meta)
}

func resourceArmNotificationHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).notificationHubsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	name := id.Path["notificationHubs"]

	resp, err := client.Get(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Notification Hub %q (Namespace %q / Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("namespace_name", namespaceName)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.Properties; props != nil {
		apns := flattenNotificationHubsAPNSCredentials(props.ApnsCredential)
		if d.Set("apns_settings", apns); err != nil {
			return fmt.Errorf("Error flattening `apns_settings`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmNotificationHubDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).notificationHubsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	name := id.Path["notificationHubs"]

	resp, err := client.Delete(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Notification Hub %q (Namespace %q / Resource Group %q): %+v", name, namespaceName, resourceGroup, err)
		}
	}

	return nil
}

func expandNotificationHubsAPNSCredentials(d *schema.ResourceData) (*notificationhubs.ApnsCredential, error) {
	inputs := d.Get("apns_credential").([]interface{})
	if len(inputs) == 0 {
		return nil, nil
	}

	input := inputs[0].(map[string]interface{})
	applicationMode := input["application_mode"].(string)
	applicationId := input["application_id"].(string)
	applicationName := input["application_name"].(string)
	keyId := input["key_id"].(string)
	token := input["token"].(string)

	applicationEndpoints := map[string]string{
		"Production": "https://api.push.apple.com:443/3/device",
		"Sandbox":    "https://api.development.push.apple.com:443/3/device",
	}
	endpoint := applicationEndpoints[applicationMode]

	credentials := notificationhubs.ApnsCredential{
		ApnsCredentialProperties: &notificationhubs.ApnsCredentialProperties{
			AppID:    utils.String(applicationId),
			AppName:  utils.String(applicationName),
			Endpoint: utils.String(endpoint),
			KeyID:    utils.String(keyId),
			Token:    utils.String(token),
		},
	}
	return &credentials, nil
}

func flattenNotificationHubsAPNSCredentials(input *notificationhubs.ApnsCredential) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	output := make(map[string]interface{}, 0)

	if applicationId := input.AppID; applicationId != nil {
		output["application_id"] = *applicationId
	}

	if name := input.AppName; name != nil {
		output["application_name"] = *name
	}

	if endpoint := input.Endpoint; endpoint != nil {
		applicationEndpoints := map[string]string{
			"https://api.push.apple.com:443/3/device":             "Production",
			"https://api.development.push.apple.com:443/3/device": "Sandbox",
		}
		applicationMode := applicationEndpoints[*endpoint]
		output["application_mode"] = applicationMode
	}

	if keyId := input.KeyID; keyId != nil {
		output["key_id"] = *keyId
	}

	if token := input.Token; token != nil {
		output["token"] = *token
	}

	return []interface{}{output}
}
