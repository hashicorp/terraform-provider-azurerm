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
		Create:        resourceArmActiveDirectoryApplicationCreate,
		Read:          resourceArmActiveDirectoryApplicationRead,
		Update:        resourceArmActiveDirectoryApplicationUpdate,
		Delete:        resourceArmActiveDirectoryApplicationDelete,
		CustomizeDiff: customizeDiffActiveDirectoryApplication,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"display_name": {
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

			"key_credential": keyCredentialsSchema(),

			"password_credential": passwordCredentialsSchema(),

			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func customizeDiffActiveDirectoryApplication(diff *schema.ResourceDiff, v interface{}) error {

	if err := customizeDiffKeyCredential(diff, v); err != nil {
		return err
	}

	if err := customizeDiffPasswordCredential(diff, v); err != nil {
		return err
	}

	return nil
}

func resourceArmActiveDirectoryApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("display_name").(string)
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

	if _, ok := d.GetOk("key_credential"); ok {
		keyCreds, err := expandAzureRmKeyCredentials(d, nil)
		if err != nil {
			return err
		}
		properties.KeyCredentials = keyCreds
	}

	if _, ok := d.GetOk("password_credential"); ok {
		passCreds, err := expandAzureRmPasswordCredentials(d, nil)
		if err != nil {
			return err
		}
		properties.PasswordCredentials = passCreds
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

	d.Set("display_name", resp.DisplayName)
	d.Set("app_id", resp.AppID)
	d.Set("object_id", resp.ObjectID)
	d.Set("homepage", resp.Homepage)
	d.Set("identifier_uris", resp.IdentifierUris)
	d.Set("reply_urls", resp.ReplyUrls)
	d.Set("available_to_other_tenants", resp.AvailableToOtherTenants)
	d.Set("oauth2_allow_implicit_flow", resp.Oauth2AllowImplicitFlow)

	rkc, err := client.ListKeyCredentials(ctx, d.Id())
	if err != nil {
		return fmt.Errorf("Error loading Azure AD Application Key Credentials %q: %+v", d.Id(), err)
	}

	if err := d.Set("key_credential", flattenAzureRmKeyCredentials(rkc.Value)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Azure AD Application Key Credentials error: %#v", err)
	}

	rpc, err := client.ListPasswordCredentials(ctx, d.Id())
	if err != nil {
		return fmt.Errorf("Error loading Azure AD Application Password Credentials %q: %+v", d.Id(), err)
	}

	if err := d.Set("password_credential", flattenAzureRmPasswordCredentials(rpc.Value)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Azure AD Application Password Credentials error: %#v", err)
	}

	return nil
}

func resourceArmActiveDirectoryApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("display_name").(string)

	var properties graphrbac.ApplicationUpdateParameters

	if d.HasChange("display_name") {
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

	d.Partial(true)

	_, err := client.Patch(ctx, d.Id(), properties)
	if err != nil {
		return err
	}

	d.SetPartial("display_name")
	d.SetPartial("homepage")
	d.SetPartial("identifier_uris")
	d.SetPartial("reply_urls")
	d.SetPartial("available_to_other_tenants")
	d.SetPartial("oauth2_allow_implicit_flow")
	d.SetPartial("app_id")
	d.SetPartial("object_id")

	if d.HasChange("key_credential") {
		o, _ := d.GetChange("key_credential")

		kc, kcErr := expandAzureRmKeyCredentials(d, o.(*schema.Set))
		if kcErr != nil {
			return kcErr
		}

		keyUpdate := graphrbac.KeyCredentialsUpdateParameters{
			Value: kc,
		}

		_, err := client.UpdateKeyCredentials(ctx, d.Id(), keyUpdate)
		if err != nil {
			return err
		}

		d.SetPartial("key_credential")
	}

	if d.HasChange("password_credential") {
		o, _ := d.GetChange("password_credential")

		pc, pcErr := expandAzureRmPasswordCredentials(d, o.(*schema.Set))
		if pcErr != nil {
			return pcErr
		}

		passUpdate := graphrbac.PasswordCredentialsUpdateParameters{
			Value: pc,
		}

		_, err := client.UpdatePasswordCredentials(ctx, d.Id(), passUpdate)
		if err != nil {
			return err
		}

		d.SetPartial("password_credential")
	}

	d.Partial(false)

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
