package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmActiveDirectoryServicePrincipalPassword() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: `The Azure Active Directory resources have been split out into their own Provider.

Information on migrating to the new AzureAD Provider can be found here: https://terraform.io/docs/providers/azurerm/guides/migrating-to-azuread.html

As such the Azure Active Directory resources within the AzureRM Provider are now deprecated and will be removed in v2.0 of the AzureRM Provider.
`,
		Create: resourceArmActiveDirectoryServicePrincipalPasswordCreate,
		Read:   resourceArmActiveDirectoryServicePrincipalPasswordRead,
		Delete: resourceArmActiveDirectoryServicePrincipalPasswordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service_principal_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.UUID,
			},

			"key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validate.UUID,
			},

			"value": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"start_date": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validate.RFC3339Time,
			},

			"end_date": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RFC3339Time,
			},
		},
	}
}

func resourceArmActiveDirectoryServicePrincipalPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).graph.ServicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	objectId := d.Get("service_principal_id").(string)
	value := d.Get("value").(string)
	// errors will be handled by the validation
	endDate, _ := time.Parse(time.RFC3339, d.Get("end_date").(string))

	var keyId string
	if v, ok := d.GetOk("key_id"); ok {
		keyId = v.(string)
	} else {
		kid, err := uuid.GenerateUUID()
		if err != nil {
			return err
		}

		keyId = kid
	}

	credential := graphrbac.PasswordCredential{
		KeyID:   utils.String(keyId),
		Value:   utils.String(value),
		EndDate: &date.Time{Time: endDate},
	}

	if v, ok := d.GetOk("start_date"); ok {
		// errors will be handled by the validation
		startDate, _ := time.Parse(time.RFC3339, v.(string))
		credential.StartDate = &date.Time{Time: startDate}
	}

	locks.ByName(objectId, servicePrincipalResourceName)
	defer locks.UnlockByName(objectId, servicePrincipalResourceName)

	existingCredentials, err := client.ListPasswordCredentials(ctx, objectId)
	if err != nil {
		return fmt.Errorf("Error Listing Password Credentials for Service Principal %q: %+v", objectId, err)
	}

	updatedCredentials := make([]graphrbac.PasswordCredential, 0)
	if existingCredentials.Value != nil {
		for _, v := range *existingCredentials.Value {
			if v.KeyID == nil {
				continue
			}

			if *v.KeyID == keyId {
				return tf.ImportAsExistsError("azurerm_azuread_service_principal_password", fmt.Sprintf("%s/%s", objectId, keyId))
			}
		}

		updatedCredentials = *existingCredentials.Value
	}

	updatedCredentials = append(updatedCredentials, credential)

	parameters := graphrbac.PasswordCredentialsUpdateParameters{
		Value: &updatedCredentials,
	}
	_, err = client.UpdatePasswordCredentials(ctx, objectId, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Password Credential %q for Service Principal %q: %+v", keyId, objectId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", objectId, keyId))

	return resourceArmActiveDirectoryServicePrincipalPasswordRead(d, meta)
}

func resourceArmActiveDirectoryServicePrincipalPasswordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).graph.ServicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	id := strings.Split(d.Id(), "/")
	if len(id) != 2 {
		return fmt.Errorf("ID should be in the format {objectId}/{keyId} - but got %q", d.Id())
	}

	objectId := id[0]
	keyId := id[1]

	// ensure the parent Service Principal exists
	servicePrincipal, err := client.Get(ctx, objectId)
	if err != nil {
		// the parent Service Principal has been removed - skip it
		if utils.ResponseWasNotFound(servicePrincipal.Response) {
			log.Printf("[DEBUG] Service Principal with Object ID %q was not found - removing from state!", objectId)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Service Principal ID %q: %+v", objectId, err)
	}

	credentials, err := client.ListPasswordCredentials(ctx, objectId)
	if err != nil {
		return fmt.Errorf("Error Listing Password Credentials for Service Principal with Object ID %q: %+v", objectId, err)
	}

	var credential *graphrbac.PasswordCredential
	for _, c := range *credentials.Value {
		if c.KeyID == nil {
			continue
		}

		if *c.KeyID == keyId {
			credential = &c
			break
		}
	}

	if credential == nil {
		log.Printf("[DEBUG] Service Principal Password %q (Object ID %q) was not found - removing from state!", keyId, objectId)
		d.SetId("")
		return nil
	}

	// value is available in the SDK but isn't returned from the API
	d.Set("key_id", credential.KeyID)
	d.Set("service_principal_id", objectId)

	if endDate := credential.EndDate; endDate != nil {
		d.Set("end_date", endDate.Format(time.RFC3339))
	}

	if startDate := credential.StartDate; startDate != nil {
		d.Set("start_date", startDate.Format(time.RFC3339))
	}

	return nil
}

func resourceArmActiveDirectoryServicePrincipalPasswordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).graph.ServicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	id := strings.Split(d.Id(), "/")
	if len(id) != 2 {
		return fmt.Errorf("ID should be in the format {objectId}/{keyId} - but got %q", d.Id())
	}

	objectId := id[0]
	keyId := id[1]

	locks.ByName(objectId, servicePrincipalResourceName)
	defer locks.UnlockByName(objectId, servicePrincipalResourceName)

	// ensure the parent Service Principal exists
	servicePrincipal, err := client.Get(ctx, objectId)
	if err != nil {
		// the parent Service Principal was removed - skip it
		if utils.ResponseWasNotFound(servicePrincipal.Response) {
			return nil
		}

		return fmt.Errorf("Error retrieving Service Principal ID %q: %+v", objectId, err)
	}

	existing, err := client.ListPasswordCredentials(ctx, objectId)
	if err != nil {
		return fmt.Errorf("Error Listing Password Credentials for Service Principal with Object ID %q: %+v", objectId, err)
	}

	updatedCredentials := make([]graphrbac.PasswordCredential, 0)
	for _, credential := range *existing.Value {
		if credential.KeyID == nil {
			continue
		}

		if *credential.KeyID != keyId {
			updatedCredentials = append(updatedCredentials, credential)
		}
	}

	parameters := graphrbac.PasswordCredentialsUpdateParameters{
		Value: &updatedCredentials,
	}
	_, err = client.UpdatePasswordCredentials(ctx, objectId, parameters)
	if err != nil {
		return fmt.Errorf("Error removing Password %q from Service Principal %q: %+v", keyId, objectId, err)
	}

	return nil
}
