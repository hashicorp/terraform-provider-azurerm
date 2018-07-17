package azurerm

import (
	"fmt"

	"log"

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

			"gcm_credential": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			// NOTE: skipping tags as there's a bug in the API where the Keys for Tags are returned in lower-case
			// Azure Rest API Specs issue: https://github.com/Azure/azure-sdk-for-go/issues/2239
			//"tags": tagsSchema(),
		},
	}
}

func resourceArmNotificationHubCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).notificationHubsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))

	apnsRaw := d.Get("apns_credential").([]interface{})
	apnsCredential, err := expandNotificationHubsAPNSCredentials(apnsRaw)
	if err != nil {
		return err
	}

	gcmRaw := d.Get("gcm_credential").([]interface{})
	gcmCredentials, err := expandNotificationHubsGCMCredentials(gcmRaw)
	if err != nil {
		return err
	}

	parameters := notificationhubs.CreateOrUpdateParameters{
		Location: utils.String(location),
		Properties: &notificationhubs.Properties{
			ApnsCredential: apnsCredential,
			GcmCredential:  gcmCredentials,
		},
	}

	_, err = client.CreateOrUpdate(ctx, resourceGroup, namespaceName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Notification Hub %q (Namespace %q / Resource Group %q): %+v", name, namespaceName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Notification Hub %q (Namespace %q / Resource Group %q): %+v", name, namespaceName, resourceGroup, err)
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
			log.Printf("[DEBUG] Notification Hub %q was not found in Namespace %q / Resource Group %q", name, namespaceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Notification Hub %q (Namespace %q / Resource Group %q): %+v", name, namespaceName, resourceGroup, err)
	}

	credentials, err := client.GetPnsCredentials(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Credentials for Notification Hub %q (Namespace %q / Resource Group %q): %+v", name, namespaceName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("namespace_name", namespaceName)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := credentials.PnsCredentialsProperties; props != nil {
		apns := flattenNotificationHubsAPNSCredentials(props.ApnsCredential)
		if d.Set("apns_settings", apns); err != nil {
			return fmt.Errorf("Error flattening `apns_settings`: %+v", err)
		}

		gcm := flattenNotificationHubsGCMCredentials(props.GcmCredential)
		if d.Set("gcm_settings", gcm); err != nil {
			return fmt.Errorf("Error flattening `gcm_settings`: %+v", err)
		}
	}

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

func expandNotificationHubsAPNSCredentials(inputs []interface{}) (*notificationhubs.ApnsCredential, error) {
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

func expandNotificationHubsGCMCredentials(inputs []interface{}) (*notificationhubs.GcmCredential, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	input := inputs[0].(map[string]interface{})
	apiKey := input["api_key"].(string)
	credentials := notificationhubs.GcmCredential{
		GcmCredentialProperties: &notificationhubs.GcmCredentialProperties{
			GoogleAPIKey: utils.String(apiKey),
		},
	}
	return &credentials, nil
}

func flattenNotificationHubsGCMCredentials(input *notificationhubs.GcmCredential) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{}, 0)
	if props := input.GcmCredentialProperties; props != nil {
		if apiKey := props.GoogleAPIKey; apiKey != nil {
			output["api_key"] = *apiKey
		}
	}

	return []interface{}{output}
}
