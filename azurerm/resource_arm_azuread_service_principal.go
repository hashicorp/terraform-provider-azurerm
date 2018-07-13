package azurerm

import (
	"fmt"
	"log"
	//"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmActiveDirectoryServicePrincipal() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmActiveDirectoryServicePrincipalCreate,
		Read:   resourceArmActiveDirectoryServicePrincipalRead,
		Update: resourceArmActiveDirectoryServicePrincipalUpdate,
		Delete: resourceArmActiveDirectoryServicePrincipalDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"application_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_principal_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"access_credential": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				// MaxItems: 16,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credential_type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								string("public_key"),
								string("password"),
							}, true),
						},
						"start_date": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.RFC3339Time,
							ValidateFunc:     validate.RFC3339DateInFutureBy(time.Duration(5) * time.Minute),
						},
						"end_date": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validate.RFC3339Time,
						},
						"value": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmActiveDirectoryServicePrincipalCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	AppID := d.Get("application_id").(string)

	passwordCredentials, err := expandAzureRmActiveDirectoryServicePasswordCredentials(d)
	if err != nil {
		return fmt.Errorf("Error expanding list of password credentials: %+v", err)
	}

	properties := graphrbac.ServicePrincipalCreateParameters{
		AppID:               &AppID,
		AccountEnabled:      utils.Bool(true),
		PasswordCredentials: &passwordCredentials,
	}

	app, err := client.Create(ctx, properties)
	if err != nil {
		return err
	}

	d.SetId(*app.ObjectID)

	return resourceArmActiveDirectoryServicePrincipalRead(d, meta)
}

func resourceArmActiveDirectoryServicePrincipalUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	passwordCredentials, passwdExpandErr := expandAzureRmActiveDirectoryServicePasswordCredentials(d)
	if passwdExpandErr != nil {
		return fmt.Errorf("Error expanding list of password credentials: %+v", passwdExpandErr)
	}

	properties := graphrbac.PasswordCredentialsUpdateParameters{
		Value: &passwordCredentials,
	}

	_, err := client.UpdatePasswordCredentials(ctx, d.Id(), properties)

	if err != nil {
		return fmt.Errorf("Error patching Azure AD Service Principal with ID %q: %+v", d.Id(), err)
	}

	return resourceArmActiveDirectoryServicePrincipalRead(d, meta)
}

func resourceArmActiveDirectoryServicePrincipalRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, d.Id())
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Azure AD Serice Principal with ID %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Azure AD Service Principal with ID %q: %+v", d.Id(), err)
	}

	d.Set("name", resp.DisplayName)
	d.Set("application_id", resp.AppID)
	d.Set("object_id", resp.ObjectID)

	servicePrincipalNames := flattenAzureADDataSourceServicePrincipalNames(resp.ServicePrincipalNames)
	if err := d.Set("service_principal_names", servicePrincipalNames); err != nil {
		return fmt.Errorf("Error setting `service_principal_names`: %+v", err)
	}

	return nil
}

func resourceArmActiveDirectoryServicePrincipalDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Delete(ctx, d.Id())
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error Deleting Azure AD Service Principal with ID %q: %+v", d.Id(), err)
		}
	}

	return nil
}

func expandAzureRmActiveDirectoryServicePrincipalNames(d *schema.ResourceData) *[]string {
	servicePrincipalNames := d.Get("service_principal_names").([]interface{})
	names := make([]string, 0)

	for _, id := range servicePrincipalNames {
		names = append(names, id.(string))
	}

	return &names
}

func expandAzureRmActiveDirectoryServicePasswordCredentials(d *schema.ResourceData) ([]graphrbac.PasswordCredential, error) {
	credentials := d.Get("access_credential").([]interface{})
	passwordCredentials := make([]graphrbac.PasswordCredential, 0, len(credentials))

	for _, configRaw := range credentials {
		data := configRaw.(map[string]interface{})

		startDate, _ := time.Parse(time.RFC3339, (data["start_date"].(string)))
		endDate, _ := time.Parse(time.RFC3339, (data["end_date"].(string)))
		value := data["value"].(string)
		//index := strconv.Itoa(i)

		passwordCredential := graphrbac.PasswordCredential{
			StartDate: &date.Time{Time: startDate},
			EndDate:   &date.Time{Time: endDate},
			Value:     &value,
			//KeyID:     &index,
		}

		/* if v := data["next_hop_in_ip_address"].(string); v != "" {
			properties.NextHopIPAddress = &v
		} */

		passwordCredentials = append(passwordCredentials, passwordCredential)
	}

	return passwordCredentials, nil
}

/* func expandAzureRmActiveDirectoryServicePasswordCredentialsUpdates(d *schema.ResourceData) ([]graphrbac.PasswordCredentialsUpdateParameters, error) {
	credentials := d.Get("access_credential").([]interface{})
	passwordCredentials := make([]graphrbac.PasswordCredentialsUpdateParameters, 0, len(credentials))

	for _, configRaw := range credentials {
		data := configRaw.(map[string]interface{})

		startDate, _ := time.Parse(time.RFC3339, (data["start_date"].(string)))
		endDate, _ := time.Parse(time.RFC3339, (data["end_date"].(string)))
		value := data["value"].(string)
		//index := strconv.Itoa(i)

		passwordCredential := graphrbac.PasswordCredentialsUpdateParameters{
			Value: &value,
		}

		if v := data["next_hop_in_ip_address"].(string); v != "" {
			properties.NextHopIPAddress = &v
		}

		passwordCredentials = append(passwordCredentials, passwordCredential)
	}

	return passwordCredentials, nil
} */
