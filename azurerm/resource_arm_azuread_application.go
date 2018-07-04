package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmActiveDirectoryApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmActiveDirectoryApplicationCreate,
		Read:   resourceArmActiveDirectoryApplicationRead,
		Update: resourceArmActiveDirectoryApplicationUpdate,
		Delete: resourceArmActiveDirectoryApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"homepage": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"identifier_uris": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"reply_urls": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"available_to_other_tenants": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"oauth2_allow_implicit_flow": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmActiveDirectoryApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	multitenant := d.Get("available_to_other_tenants").(bool)

	properties := graphrbac.ApplicationCreateParameters{
		DisplayName:             &name,
		Homepage:                expandAzureRmActiveDirectoryApplicationHomepage(d, name),
		IdentifierUris:          expandAzureRmActiveDirectoryApplicationIdentifierUris(d, name),
		ReplyUrls:               expandAzureRmActiveDirectoryApplicationReplyUrls(d, name),
		AvailableToOtherTenants: utils.Bool(multitenant),
	}

	if v, ok := d.GetOk("oauth2_allow_implicit_flow"); ok {
		properties.Oauth2AllowImplicitFlow = utils.Bool(v.(bool))
	}

	app, err := client.Create(ctx, properties)
	if err != nil {
		return err
	}

	d.SetId(*app.ObjectID)

	return resourceArmActiveDirectoryApplicationRead(d, meta)
}

func resourceArmActiveDirectoryApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, d.Id())
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Azure AD Application ID %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error loading Azure AD Application %q: %+v", d.Id(), err)
	}

	d.Set("name", resp.DisplayName)
	d.Set("application_id", resp.AppID)
	d.Set("homepage", resp.Homepage)
	d.Set("identifier_uris", resp.IdentifierUris)
	d.Set("reply_urls", resp.ReplyUrls)
	d.Set("available_to_other_tenants", resp.AvailableToOtherTenants)
	d.Set("oauth2_allow_implicit_flow", resp.Oauth2AllowImplicitFlow)

	return nil
}

func resourceArmActiveDirectoryApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)

	var properties graphrbac.ApplicationUpdateParameters

	if d.HasChange("name") {
		properties.DisplayName = &name
	}

	if d.HasChange("homepage") {
		properties.Homepage = expandAzureRmActiveDirectoryApplicationHomepage(d, name)
	}

	if d.HasChange("identifier_uris") {
		properties.IdentifierUris = expandAzureRmActiveDirectoryApplicationIdentifierUris(d, name)
	}

	if d.HasChange("reply_urls") {
		properties.ReplyUrls = expandAzureRmActiveDirectoryApplicationReplyUrls(d, name)
	}

	if d.HasChange("available_to_other_tenants") {
		multitenant := d.Get("available_to_other_tenants").(bool)
		properties.AvailableToOtherTenants = utils.Bool(multitenant)
	}

	if d.HasChange("oauth2_allow_implicit_flow") {
		oauth := d.Get("oauth2_allow_implicit_flow").(bool)
		properties.Oauth2AllowImplicitFlow = utils.Bool(oauth)
	}

	_, err := client.Patch(ctx, d.Id(), properties)
	if err != nil {
		return err
	}

	return resourceArmActiveDirectoryApplicationRead(d, meta)
}

func resourceArmActiveDirectoryApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Delete(ctx, d.Id())
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}

func expandAzureRmActiveDirectoryApplicationHomepage(d *schema.ResourceData, name string) *string {
	if v, ok := d.GetOk("homepage"); ok {
		return utils.String(v.(string))
	}

	return utils.String(fmt.Sprintf("http://%s", name))
}

func expandAzureRmActiveDirectoryApplicationIdentifierUris(d *schema.ResourceData, name string) *[]string {
	identifierUris := d.Get("identifier_uris").([]interface{})
	identifiers := []string{}
	for _, id := range identifierUris {
		identifiers = append(identifiers, id.(string))
	}
	if len(identifiers) == 0 {
		identifiers = append(identifiers, fmt.Sprintf("http://%s", name))
	}

	return &identifiers
}

func expandAzureRmActiveDirectoryApplicationReplyUrls(d *schema.ResourceData, name string) *[]string {
	replyUrls := d.Get("reply_urls").([]interface{})
	urls := []string{}
	for _, url := range replyUrls {
		urls = append(urls, url.(string))
	}

	return &urls
}
