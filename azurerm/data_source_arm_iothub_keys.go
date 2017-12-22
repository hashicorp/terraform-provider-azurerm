package azurerm

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/dataplane/keyvault"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceArmIotHubKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmIotHubSkuRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_name": resourceGroupNameForDataSourceSchema(),
			"key_type": {
				Type:     schema.TypeString,
				Required: true,
				// turns out Azure's *really* sensitive about the casing of these
				// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
				ValidateFunc: validation.StringInSlice([]string{
					// TODO: add `oct` back in once this is fixed
					// https://github.com/Azure/azure-rest-api-specs/issues/1739#issuecomment-332236257
					string(keyvault.EC),
					string(keyvault.RSA),
					string(keyvault.RSAHSM),
				}, false),
			},
			"key_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"key_opts": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					// turns out Azure's *really* sensitive about the casing of these
					// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
					ValidateFunc: validation.StringInSlice([]string{
						string(keyvault.Decrypt),
						string(keyvault.Encrypt),
						string(keyvault.Sign),
						string(keyvault.UnwrapKey),
						string(keyvault.Verify),
						string(keyvault.WrapKey),
					}, false),
				},
			},
			// Computed
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"n": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"e": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func dataSourceArmIotHubKeysRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)

	log.Printf("[INFO] Acquiring Keys for Authentication")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	keyName := d.Get("key_name").(string)
	resourceId := &ResourceID{
		SubscriptionID: armClient.subscriptionId,
		ResourceGroup:  resourceGroup,
	}
	resourceIdString, err := composeAzureResourceID(resourceId)

	if err != nil {
		return err
	}

	d.SetId(resourceIdString)

	if err := resourceArmResourceGroupRead(d, meta); err != nil {
		return err
	}
	// query := &iothub.SharedAccessSignatureAuthorizationRule{
	// 	KeyName: &keyName,
	// }

	// subscriptionID := armClient.subscriptionId
	// resourceId := &ResourceID{
	// 	SubscriptionID: armClient.iothubResourceClient.SubscriptionID,
	// 	ResourceGroup:  name,
	// }
	// desc := iothub.Description{
	// 	Resourcegroup:  &resourceGroup,
	// 	Name:           &name,
	// 	Subscriptionid: &subscriptionID,
	// }
	// resourceIdString, err := composeAzureResourceID(resourceId)

	// if err != nil {
	// 	return err
	// }

	// d.SetId(resourceIdString)

	return nil
}
