package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
)

func dataSourceArmAzureADApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAzureADApplicationRead,

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

			"homepage": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"identifier_uris": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"reply_urls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"available_to_other_tenants": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"oauth2_allow_implicit_flow": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmAzureADApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	var application graphrbac.Application

	if oId, ok := d.GetOk("object_id"); ok {
		objectId := oId.(string)
		resp, err := client.Get(ctx, objectId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error: AzureAD Application with ID %q was not found", objectId)
			}

			return fmt.Errorf("Error making Read request on AzureAD Application with ID %q: %+v", objectId, err)
		}

		application = resp
	} else {

		// use the name to find the Azure AD application
		name := d.Get("name").(string)
		filter := "displayName eq '" + name + "'"
		log.Printf("[DEBUG] [data_source_azuread_application] Using filter %q", filter)

		resp, err := client.ListComplete(ctx, filter)

		if err != nil {
			return fmt.Errorf("Error listing Azure AD Applications: %+v", err)
		}

		var app *graphrbac.Application
		for _, v := range *resp.Response().Value {
			if v.DisplayName != nil {
				if *v.DisplayName == name {
					app = &v
					break
				}
			}
		}

		if app == nil {
			return fmt.Errorf("Couldn't locate an Azure AD Application with a name of %q", name)
		}

		application = *app
	}

	d.SetId(*application.ObjectID)

	d.Set("object_id", application.ObjectID)
	d.Set("name", application.DisplayName)
	d.Set("application_id", application.AppID)
	d.Set("homepage", application.Homepage)
	d.Set("available_to_other_tenants", application.AvailableToOtherTenants)
	d.Set("oauth2_allow_implicit_flow", application.Oauth2AllowImplicitFlow)

	identifierUris := flattenAzureADDataSourceApplicationIdentifierUris(application.IdentifierUris)
	if err := d.Set("identifier_uris", identifierUris); err != nil {
		return fmt.Errorf("Error setting `identifier_uris`: %+v", err)
	}

	replyUrls := flattenAzureADDataSourceApplicationReplyUrls(application.ReplyUrls)
	if err := d.Set("reply_urls", replyUrls); err != nil {
		return fmt.Errorf("Error setting `reply_urls`: %+v", err)
	}

	return nil
}

func flattenAzureADDataSourceApplicationIdentifierUris(input *[]string) []string {
	output := make([]string, 0)

	if input != nil {
		for _, v := range *input {
			output = append(output, v)
		}
	}

	return output
}

func flattenAzureADDataSourceApplicationReplyUrls(input *[]string) []string {
	output := make([]string, 0)

	if input != nil {
		for _, v := range *input {
			output = append(output, v)
		}
	}

	return output
}
