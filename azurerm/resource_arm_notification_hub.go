package azurerm

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var notificationHubResourceName = "azurerm_notification_hub"

const apnsProductionName = "Production"
const apnsProductionEndpoint = "https://api.push.apple.com:443/3/device"
const apnsSandboxName = "Sandbox"
const apnsSandboxEndpoint = "https://api.development.push.apple.com:443/3/device"

func resourceArmNotificationHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNotificationHubCreateUpdate,
		Read:   resourceArmNotificationHubRead,
		Update: resourceArmNotificationHubCreateUpdate,
		Delete: resourceArmNotificationHubDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			// NOTE: the ForceNew is to workaround a bug in the Azure SDK where nil-values aren't sent to the API.
			// Bug: https://github.com/Azure/azure-sdk-for-go/issues/2246

			oAPNS, nAPNS := diff.GetChange("apns_credential.#")
			oAPNSi := oAPNS.(int)
			nAPNSi := nAPNS.(int)
			if nAPNSi < oAPNSi {
				diff.ForceNew("apns_credential")
			}

			oGCM, nGCM := diff.GetChange("gcm_credential.#")
			oGCMi := oGCM.(int)
			nGCMi := nGCM.(int)
			if nGCMi < oGCMi {
				diff.ForceNew("gcm_credential")
			}

			return nil
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"apns_credential": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// NOTE: APNS supports two modes, certificate auth (v1) and token auth (v2)
						// certificate authentication/v1 is marked for deprecation; as such we're not
						// supporting it at this time.
						"application_mode": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								apnsProductionName,
								apnsSandboxName,
							}, false),
						},
						"bundle_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						// Team ID (within Apple & the Portal) == "AppID" (within the API)
						"team_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"token": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
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
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
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
	client := meta.(*ArmClient).notificationHubs.HubsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Notification Hub %q (Namespace %q / Resource Group %q): %+v", name, namespaceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_notification_hub", *existing.ID)
		}
	}

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

	if _, err = client.CreateOrUpdate(ctx, resourceGroup, namespaceName, name, parameters); err != nil {
		return fmt.Errorf("Error creating Notification Hub %q (Namespace %q / Resource Group %q): %+v", name, namespaceName, resourceGroup, err)
	}

	// Notification Hubs are eventually consistent
	log.Printf("[DEBUG] Waiting for Notification Hub %q to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   notificationHubStateRefreshFunc(ctx, client, resourceGroup, namespaceName, name),
		Timeout:                   10 * time.Minute,
		MinTimeout:                15 * time.Second,
		ContinuousTargetOccurence: 10,
	}
	if _, err2 := stateConf.WaitForState(); err2 != nil {
		return fmt.Errorf("Error waiting for Notification Hub %q to become available: %s", name, err2)
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

func notificationHubStateRefreshFunc(ctx context.Context, client *notificationhubs.Client, resourceGroup, namespaceName, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return nil, "404", nil
			}

			return nil, "", fmt.Errorf("Error retrieving Notification Hub %q (Namespace %q / Resource Group %q): %+v", name, namespaceName, resourceGroup, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func resourceArmNotificationHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).notificationHubs.HubsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := credentials.PnsCredentialsProperties; props != nil {
		apns := flattenNotificationHubsAPNSCredentials(props.ApnsCredential)
		if setErr := d.Set("apns_credential", apns); setErr != nil {
			return fmt.Errorf("Error setting `apns_credential`: %+v", setErr)
		}

		gcm := flattenNotificationHubsGCMCredentials(props.GcmCredential)
		if setErr := d.Set("gcm_credential", gcm); setErr != nil {
			return fmt.Errorf("Error setting `gcm_credential`: %+v", setErr)
		}
	}

	return nil
}

func resourceArmNotificationHubDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).notificationHubs.HubsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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
	bundleId := input["bundle_id"].(string)
	keyId := input["key_id"].(string)
	teamId := input["team_id"].(string)
	token := input["token"].(string)

	applicationEndpoints := map[string]string{
		apnsProductionName: apnsProductionEndpoint,
		apnsSandboxName:    apnsSandboxEndpoint,
	}
	endpoint := applicationEndpoints[applicationMode]

	credentials := notificationhubs.ApnsCredential{
		ApnsCredentialProperties: &notificationhubs.ApnsCredentialProperties{
			AppID:    utils.String(teamId),
			AppName:  utils.String(bundleId),
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

	output := make(map[string]interface{})

	if bundleId := input.AppName; bundleId != nil {
		output["bundle_id"] = *bundleId
	}

	if endpoint := input.Endpoint; endpoint != nil {
		applicationEndpoints := map[string]string{
			apnsProductionEndpoint: apnsProductionName,
			apnsSandboxEndpoint:    apnsSandboxName,
		}
		applicationMode := applicationEndpoints[*endpoint]
		output["application_mode"] = applicationMode
	}

	if keyId := input.KeyID; keyId != nil {
		output["key_id"] = *keyId
	}

	if teamId := input.AppID; teamId != nil {
		output["team_id"] = *teamId
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

	output := make(map[string]interface{})
	if props := input.GcmCredentialProperties; props != nil {
		if apiKey := props.GoogleAPIKey; apiKey != nil {
			output["api_key"] = *apiKey
		}
	}

	return []interface{}{output}
}
