package azuread

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/graph"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/p"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
)

const resourceApplicationName = "azuread_application"

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"available_to_other_tenants": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"group_membership_claims": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"All", "None", "SecurityGroup", "DirectoryRole", "DistributionGroup"},
					false,
				),
			},

			"homepage": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.URLIsHTTPS,
			},

			"identifier_uris": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.URLIsAppURI,
				},
			},

			"oauth2_allow_implicit_flow": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"public_client": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"reply_urls": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.NoEmptyStrings,
				},
			},

			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"webapp/api", "native"}, false),
				Default:      "webapp/api",
			},

			"app_role": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"allowed_member_types": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice(
									[]string{"User", "Application"},
									false,
								),
							},
						},

						"description": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"display_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"is_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"required_resource_access": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_app_id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"resource_access": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.UUID,
									},

									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice(
											[]string{"Scope", "Role"},
											false, // force case sensitivity
										),
									},
								},
							},
						},
					},
				},
			},

			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"oauth2_permissions": graph.SchemaOauth2Permissions(),
		},
	}
}

func resourceApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	appType := d.Get("type")
	identUrls, hasIdentUrls := d.GetOk("identifier_uris")
	if appType == "native" {
		if hasIdentUrls {
			return fmt.Errorf("identifier_uris is not required for a native application")
		}
	}

	properties := graphrbac.ApplicationCreateParameters{
		DisplayName:             &name,
		IdentifierUris:          tf.ExpandStringSlicePtr(identUrls.([]interface{})),
		ReplyUrls:               tf.ExpandStringSlicePtr(d.Get("reply_urls").(*schema.Set).List()),
		AvailableToOtherTenants: p.Bool(d.Get("available_to_other_tenants").(bool)),
		RequiredResourceAccess:  expandADApplicationRequiredResourceAccess(d),
		AppRoles:                expandADApplicationAppRoles(d.Get("app_role")),
	}

	if v, ok := d.GetOk("homepage"); ok {
		properties.Homepage = p.String(v.(string))
	} else {
		// continue to automatically set the homepage with the type is not native
		if appType != "native" {
			properties.Homepage = p.String(fmt.Sprintf("https://%s", name))

		}
	}

	if v, ok := d.GetOk("oauth2_allow_implicit_flow"); ok {
		properties.Oauth2AllowImplicitFlow = p.Bool(v.(bool))
	}

	if v, ok := d.GetOk("public_client"); ok {
		properties.PublicClient = p.Bool(v.(bool))
	}

	if v, ok := d.GetOk("group_membership_claims"); ok {
		properties.GroupMembershipClaims = v
	}

	app, err := client.Create(ctx, properties)
	if err != nil {
		return err
	}
	if app.ObjectID == nil {
		return fmt.Errorf("Application objectId is nil")
	}
	d.SetId(*app.ObjectID)

	_, err = graph.WaitForReplication(func() (interface{}, error) {
		return client.Get(ctx, *app.ObjectID)
	})
	if err != nil {
		return fmt.Errorf("Error waiting for Application with ObjectId %q: %+v", *app.ObjectID, err)
	}

	// follow suggested hack for azure-cli
	// AAD graph doesn't have the API to create a native app, aka public client, the recommended hack is
	// to create a web app first, then convert to a native one
	if appType == "native" {

		properties := graphrbac.ApplicationUpdateParameters{
			Homepage:       nil,
			IdentifierUris: &[]string{},
			PublicClient:   p.Bool(true),
		}
		if _, err := client.Patch(ctx, *app.ObjectID, properties); err != nil {
			return err
		}
	}

	return resourceApplicationRead(d, meta)
}

func resourceApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)

	var properties graphrbac.ApplicationUpdateParameters

	if d.HasChange("name") {
		properties.DisplayName = &name
	}

	if d.HasChange("homepage") {
		properties.Homepage = p.String(d.Get("homepage").(string))
	}

	if d.HasChange("identifier_uris") {
		properties.IdentifierUris = tf.ExpandStringSlicePtr(d.Get("identifier_uris").([]interface{}))
	}

	if d.HasChange("reply_urls") {
		properties.ReplyUrls = tf.ExpandStringSlicePtr(d.Get("reply_urls").(*schema.Set).List())
	}

	if d.HasChange("available_to_other_tenants") {
		properties.AvailableToOtherTenants = p.Bool(d.Get("available_to_other_tenants").(bool))
	}

	if d.HasChange("oauth2_allow_implicit_flow") {
		properties.Oauth2AllowImplicitFlow = p.Bool(d.Get("oauth2_allow_implicit_flow").(bool))
	}

	if d.HasChange("public_client") {
		properties.PublicClient = p.Bool(d.Get("public_client").(bool))
	}

	if d.HasChange("required_resource_access") {
		properties.RequiredResourceAccess = expandADApplicationRequiredResourceAccess(d)
	}

	if d.HasChange("app_role") {
		// if the app role already exists then it must be disabled
		// with no other changes before it can be edited or deleted
		var app graphrbac.Application
		var appRolesProperties graphrbac.ApplicationUpdateParameters
		resp, err := client.Get(ctx, d.Id())
		if err != nil {
			if ar.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error: AzureAD Application with ID %q was not found", d.Id())
			}

			return fmt.Errorf("Error making Read request on AzureAD Application with ID %q: %+v", d.Id(), err)
		}
		app = resp
		for _, appRole := range *app.AppRoles {
			*appRole.IsEnabled = false
		}
		appRolesProperties.AppRoles = app.AppRoles
		if _, err := client.Patch(ctx, d.Id(), appRolesProperties); err != nil {
			return fmt.Errorf("Error disabling App Roles for Azure AD Application with ID %q: %+v", d.Id(), err)
		}

		// now we can set the new state of the app role
		properties.AppRoles = expandADApplicationAppRoles(d.Get("app_role"))
	}

	if d.HasChange("group_membership_claims") {
		properties.GroupMembershipClaims = d.Get("group_membership_claims")
	}

	if d.HasChange("type") {
		switch appType := d.Get("type"); appType {
		case "webapp/api":
			properties.PublicClient = p.Bool(false)
			properties.IdentifierUris = tf.ExpandStringSlicePtr(d.Get("identifier_uris").([]interface{}))
		case "native":
			properties.PublicClient = p.Bool(true)
			properties.IdentifierUris = &[]string{}
		default:
			return fmt.Errorf("Error paching Azure AD Application with ID %q: Unknow application type %v. Supported types are [webapp/api, native]", d.Id(), appType)
		}
	}

	if _, err := client.Patch(ctx, d.Id(), properties); err != nil {
		return fmt.Errorf("Error patching Azure AD Application with ID %q: %+v", d.Id(), err)
	}

	return resourceApplicationRead(d, meta)
}

func resourceApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	app, err := client.Get(ctx, d.Id())
	if err != nil {
		if ar.ResponseWasNotFound(app.Response) {
			log.Printf("[DEBUG] Azure AD Application with ID %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Azure AD Application with ID %q: %+v", d.Id(), err)
	}

	d.Set("name", app.DisplayName)
	d.Set("application_id", app.AppID)
	d.Set("homepage", app.Homepage)
	d.Set("available_to_other_tenants", app.AvailableToOtherTenants)
	d.Set("oauth2_allow_implicit_flow", app.Oauth2AllowImplicitFlow)
	d.Set("public_client", app.PublicClient)
	d.Set("object_id", app.ObjectID)

	if v := app.PublicClient; v != nil && *v {
		d.Set("type", "native")
	} else {
		d.Set("type", "webapp/api")
	}

	if err := d.Set("group_membership_claims", app.GroupMembershipClaims); err != nil {
		return fmt.Errorf("Error setting `group_membership_claims`: %+v", err)
	}

	if err := d.Set("identifier_uris", tf.FlattenStringSlicePtr(app.IdentifierUris)); err != nil {
		return fmt.Errorf("Error setting `identifier_uris`: %+v", err)
	}

	if err := d.Set("reply_urls", tf.FlattenStringSlicePtr(app.ReplyUrls)); err != nil {
		return fmt.Errorf("Error setting `reply_urls`: %+v", err)
	}

	if err := d.Set("required_resource_access", flattenADApplicationRequiredResourceAccess(app.RequiredResourceAccess)); err != nil {
		return fmt.Errorf("Error setting `required_resource_access`: %+v", err)
	}

	if err := d.Set("app_role", graph.FlattenAppRoles(app.AppRoles)); err != nil {
		return fmt.Errorf("Error setting `app_role`: %+v", err)
	}

	if err := d.Set("oauth2_permissions", graph.FlattenOauth2Permissions(app.Oauth2Permissions)); err != nil {
		return fmt.Errorf("Error setting `oauth2_permissions`: %+v", err)
	}

	return nil
}

func resourceApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	// in order to delete an application which is available to other tenants, we first have to disable this setting
	availableToOtherTenants := d.Get("available_to_other_tenants").(bool)
	if availableToOtherTenants {
		log.Printf("[DEBUG] Azure AD Application is available to other tenants - disabling that feature before deleting.")
		properties := graphrbac.ApplicationUpdateParameters{
			AvailableToOtherTenants: p.Bool(false),
		}

		if _, err := client.Patch(ctx, d.Id(), properties); err != nil {
			return fmt.Errorf("Error patching Azure AD Application with ID %q: %+v", d.Id(), err)
		}
	}

	resp, err := client.Delete(ctx, d.Id())
	if err != nil {
		if !ar.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error Deleting Azure AD Application with ID %q: %+v", d.Id(), err)
		}
	}

	return nil
}

func expandADApplicationRequiredResourceAccess(d *schema.ResourceData) *[]graphrbac.RequiredResourceAccess {
	requiredResourcesAccesses := d.Get("required_resource_access").(*schema.Set).List()
	result := make([]graphrbac.RequiredResourceAccess, 0)

	for _, raw := range requiredResourcesAccesses {
		requiredResourceAccess := raw.(map[string]interface{})
		resource_app_id := requiredResourceAccess["resource_app_id"].(string)

		result = append(result,
			graphrbac.RequiredResourceAccess{
				ResourceAppID: &resource_app_id,
				ResourceAccess: expandADApplicationResourceAccess(
					requiredResourceAccess["resource_access"].([]interface{}),
				),
			},
		)
	}
	return &result
}

func expandADApplicationResourceAccess(in []interface{}) *[]graphrbac.ResourceAccess {
	var resourceAccesses []graphrbac.ResourceAccess
	for _, resource_access_raw := range in {
		resource_access := resource_access_raw.(map[string]interface{})

		resourceId := resource_access["id"].(string)
		resourceType := resource_access["type"].(string)

		resourceAccesses = append(resourceAccesses,
			graphrbac.ResourceAccess{
				ID:   &resourceId,
				Type: &resourceType,
			},
		)
	}

	return &resourceAccesses
}

func flattenADApplicationRequiredResourceAccess(in *[]graphrbac.RequiredResourceAccess) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0, len(*in))
	for _, requiredResourceAccess := range *in {
		resource := make(map[string]interface{})
		if requiredResourceAccess.ResourceAppID != nil {
			resource["resource_app_id"] = *requiredResourceAccess.ResourceAppID
		}

		resource["resource_access"] = flattenADApplicationResourceAccess(requiredResourceAccess.ResourceAccess)

		result = append(result, resource)
	}

	return result
}

func flattenADApplicationResourceAccess(in *[]graphrbac.ResourceAccess) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	accesses := make([]interface{}, 0)
	for _, resourceAccess := range *in {
		access := make(map[string]interface{})
		if resourceAccess.ID != nil {
			access["id"] = *resourceAccess.ID
		}
		if resourceAccess.Type != nil {
			access["type"] = *resourceAccess.Type
		}
		accesses = append(accesses, access)
	}

	return accesses
}

func expandADApplicationAppRoles(i interface{}) *[]graphrbac.AppRole {
	input := i.(*schema.Set).List()
	if len(input) == 0 {
		return nil
	}

	var output []graphrbac.AppRole

	for _, appRoleRaw := range input {
		appRole := appRoleRaw.(map[string]interface{})

		appRoleID := appRole["id"].(string)
		if appRoleID == "" {
			appRoleID = uuid.New().String()
		}

		var appRoleAllowedMemberTypes []string
		for _, appRoleAllowedMemberType := range appRole["allowed_member_types"].(*schema.Set).List() {
			appRoleAllowedMemberTypes = append(appRoleAllowedMemberTypes, appRoleAllowedMemberType.(string))
		}

		appRoleDescription := appRole["description"].(string)
		appRoleDisplayName := appRole["display_name"].(string)
		appRoleIsEnabled := appRole["is_enabled"].(bool)
		appRoleValue := appRole["value"].(string)

		output = append(output,
			graphrbac.AppRole{
				ID:                 &appRoleID,
				AllowedMemberTypes: &appRoleAllowedMemberTypes,
				Description:        &appRoleDescription,
				DisplayName:        &appRoleDisplayName,
				IsEnabled:          &appRoleIsEnabled,
				Value:              &appRoleValue,
			},
		)
	}

	return &output
}
