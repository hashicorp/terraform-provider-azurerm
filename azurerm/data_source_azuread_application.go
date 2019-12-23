package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmAzureADApplication() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: `The Azure Active Directory resources have been split out into their own Provider.

Information on migrating to the new AzureAD Provider can be found here: https://terraform.io/docs/providers/azurerm/guides/migrating-to-azuread.html

As such the Azure Active Directory resources within the AzureRM Provider are now deprecated and will be removed in v2.0 of the AzureRM Provider.
`,
		Read: dataSourceArmAzureADApplicationRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

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
	client := meta.(*clients.Client).Graph.ApplicationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var application graphrbac.Application

	if oId, ok := d.GetOk("object_id"); ok {
		// use the object_id to find the Azure AD application
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
		filter := fmt.Sprintf("displayName eq '%s'", name)
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

	if err := d.Set("identifier_uris", application.IdentifierUris); err != nil {
		return fmt.Errorf("Error setting `identifier_uris`: %+v", err)
	}

	if err := d.Set("reply_urls", application.ReplyUrls); err != nil {
		return fmt.Errorf("Error setting `reply_urls`: %+v", err)
	}

	return nil
}
