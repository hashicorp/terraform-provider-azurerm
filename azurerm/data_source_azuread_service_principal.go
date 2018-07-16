package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmAzureADServicePrincipal() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAzureADServicePrincipalRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		// TODO: customizeDiff for validation of either name or object_id.

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"object_id"},
			},

			"application_id": {
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
		},
	}
}

func dataSourceArmAzureADServicePrincipalRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	var serviceprincipal graphrbac.ServicePrincipal

	if oId, ok := d.GetOk("object_id"); ok {
		objectId := oId.(string)
		resp, err := client.Get(ctx, objectId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error: AzureAD Application with ID %q was not found", objectId)
			}

			return fmt.Errorf("Error making Read request on AzureAD Application with ID %q: %+v", objectId, err)
		}

		serviceprincipal = resp
	} else {
		resp, err := client.ListComplete(ctx, "")
		if err != nil {
			return fmt.Errorf("Error listing Azure AD Service Principals: %+v", err)
		}

		name := d.Get("name").(string)

		var sp *graphrbac.ServicePrincipal
		for _, v := range *resp.Response().Value {
			if v.DisplayName != nil {
				if *v.DisplayName == name {
					sp = &v
					break
				}
			}
		}

		if sp == nil {
			return fmt.Errorf("Couldn't locate an Azure AD Service Principal with a name of %q", name)
		}

		serviceprincipal = *sp
	}

	d.SetId(*serviceprincipal.ObjectID)

	d.Set("object_id", serviceprincipal.ObjectID)
	d.Set("name", serviceprincipal.DisplayName)
	d.Set("application_id", serviceprincipal.AppID)

	servicePrincipalNames := flattenAzureADDataSourceServicePrincipalNames(serviceprincipal.ServicePrincipalNames)
	if err := d.Set("service_principal_names", servicePrincipalNames); err != nil {
		return fmt.Errorf("Error setting `service_principal_names`: %+v", err)
	}

	return nil
}

func flattenAzureADDataSourceServicePrincipalNames(input *[]string) []string {
	output := make([]string, 0)

	if input != nil {
		for _, v := range *input {
			output = append(output, v)
		}
	}

	return output
}
