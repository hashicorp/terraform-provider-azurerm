package azuread

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/graph"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
)

func dataServicePrincipal() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceActiveDirectoryServicePrincipalRead,

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validate.UUID,
				ConflictsWith: []string{"display_name", "application_id"},
			},

			"display_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validate.NoEmptyStrings,
				ConflictsWith: []string{"object_id", "application_id"},
			},

			"application_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validate.UUID,
				ConflictsWith: []string{"object_id", "display_name"},
			},

			"app_roles": graph.SchemaAppRoles(),

			"oauth2_permissions": graph.SchemaOauth2Permissions(),
		},
	}
}

func dataSourceActiveDirectoryServicePrincipalRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	var sp *graphrbac.ServicePrincipal

	if v, ok := d.GetOk("object_id"); ok {

		//use the object_id to find the Azure AD service principal
		objectId := v.(string)
		app, err := client.Get(ctx, objectId)
		if err != nil {
			if ar.ResponseWasNotFound(app.Response) {
				return fmt.Errorf("Service Principal with Object ID %q was not found!", objectId)
			}

			return fmt.Errorf("Error retrieving Service Principal ID %q: %+v", objectId, err)
		}

		sp = &app

	} else if _, ok := d.GetOk("display_name"); ok {

		// use the display_name to find the Azure AD service principal
		displayName := d.Get("display_name").(string)
		filter := fmt.Sprintf("displayName eq '%s'", displayName)

		apps, err := client.ListComplete(ctx, filter)
		if err != nil {
			return fmt.Errorf("Error listing Service Principals: %+v", err)
		}

		for _, app := range *apps.Response().Value {
			if app.DisplayName == nil {
				continue
			}

			if *app.DisplayName == displayName {
				sp = &app
				break
			}
		}

		if sp == nil {
			return fmt.Errorf("A Service Principal with the Display Name %q was not found", displayName)
		}

	} else {

		// use the application_id to find the Azure AD service principal
		applicationId := d.Get("application_id").(string)
		filter := fmt.Sprintf("appId eq '%s'", applicationId)

		apps, err := client.ListComplete(ctx, filter)
		if err != nil {
			return fmt.Errorf("Error listing Service Principals: %+v", err)
		}

		for _, app := range *apps.Response().Value {
			if app.AppID == nil {
				continue
			}

			if *app.AppID == applicationId {
				sp = &app
				break
			}
		}

		if sp == nil {
			return fmt.Errorf("A Service Principal for Application ID %q was not found", applicationId)
		}

	}

	if sp.ObjectID == nil {
		return fmt.Errorf("Service Principal objectId is nil")
	}
	d.SetId(*sp.ObjectID)

	d.Set("application_id", sp.AppID)
	d.Set("display_name", sp.DisplayName)
	d.Set("object_id", sp.ObjectID)

	if err := d.Set("app_roles", graph.FlattenAppRoles(sp.AppRoles)); err != nil {
		return fmt.Errorf("Error setting `app_roles`: %+v", err)
	}

	if err := d.Set("oauth2_permissions", graph.FlattenOauth2Permissions(sp.Oauth2Permissions)); err != nil {
		return fmt.Errorf("Error setting `oauth2_permissions`: %+v", err)
	}

	return nil
}
