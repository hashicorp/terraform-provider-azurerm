package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAdApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAdApplicationCreate,
		Read:   resourceArmAdApplicationRead,
		Update: resourceArmAdApplicationUpdate,
		Delete: resourceArmAdApplicationDelete,
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

func resourceArmAdApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("display_name").(string)
	multitenant := d.Get("available_to_other_tenants").(bool)

	properties := graphrbac.ApplicationCreateParameters{
		DisplayName:             &name,
		Homepage:                expandAzureRmAdApplicationHomepage(d, name),
		IdentifierUris:          expandAzureRmAdApplicationIdentifierUris(d, name),
		ReplyUrls:               expandAzureRmAdApplicationReplyUrls(d, name),
		AvailableToOtherTenants: utils.Bool(multitenant),
	}

	if v, ok := d.GetOk("oauth2_allow_implicit_flow"); ok {
		properties.Oauth2AllowImplicitFlow = utils.Bool(v.(bool))
	}

	if _, ok := d.GetOk("key_credential"); ok {
		keyCreds, err := expandAzureRmKeyCredentials(d)
		if err != nil {
			return err
		}
		properties.KeyCredentials = keyCreds
	}

	if _, ok := d.GetOk("password_credential"); ok {
		passCreds, err := expandAzureRmPasswordCredentials(d)
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

	return resourceArmAdApplicationRead(d, meta)
}

func resourceArmAdApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, d.Id())
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Application ID %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error loading Application %q: %+v", d.Id(), err)
	}

	d.Set("name", resp.DisplayName)
	d.Set("app_id", resp.AppID)
	d.Set("object_id", resp.ObjectID)
	d.Set("homepage", resp.Homepage)
	d.Set("identifier_uris", resp.IdentifierUris)
	d.Set("reply_urls", resp.ReplyUrls)
	d.Set("available_to_other_tenants", resp.AvailableToOtherTenants)
	d.Set("oauth2_allow_implicit_flow", resp.Oauth2AllowImplicitFlow)

	respKeyCreds, err := client.ListKeyCredentials(ctx, d.Id())
	if err != nil {
		return fmt.Errorf("Error loading Application Key Credentials %q: %+v", d.Id(), err)
	}

	if err := d.Set("key_credential", flattenAzureRmKeyCredentials(respKeyCreds.Value)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Application Key Credentials error: %#v", err)
	}

	respPassCreds, err := client.ListPasswordCredentials(ctx, d.Id())
	if err != nil {
		return fmt.Errorf("Error loading Application Password Credentials %q: %+v", d.Id(), err)
	}

	if err := d.Set("password_credential", flattenAzureRmPasswordCredentials(respPassCreds.Value)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Application Password Credentials error: %#v", err)
	}

	return nil
}

func resourceArmAdApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("display_name").(string)
	multitenant := d.Get("available_to_other_tenants").(bool)

	properties := graphrbac.ApplicationUpdateParameters{
		DisplayName:             &name,
		Homepage:                expandAzureRmAdApplicationHomepage(d, name),
		IdentifierUris:          expandAzureRmAdApplicationIdentifierUris(d, name),
		ReplyUrls:               expandAzureRmAdApplicationReplyUrls(d, name),
		AvailableToOtherTenants: utils.Bool(multitenant),
	}

	if v, ok := d.GetOk("oauth2_allow_implicit_flow"); ok {
		properties.Oauth2AllowImplicitFlow = utils.Bool(v.(bool))
	}

	keyCreds, keyErr := expandAzureRmKeyCredentials(d)
	if keyErr != nil {
		return keyErr
	}

	passCreds, passErr := expandAzureRmPasswordCredentials(d)
	if passErr != nil {
		return passErr
	}

	_, err := client.Patch(ctx, d.Id(), properties)
	if err != nil {
		return err
	}

	if d.HasChange("key_credential") {
		keyUpdate := graphrbac.KeyCredentialsUpdateParameters{
			Value: keyCreds,
		}

		_, err := client.UpdateKeyCredentials(ctx, d.Id(), keyUpdate)
		if err != nil {
			return err
		}
	}

	if d.HasChange("password_credential") {
		passUpdate := graphrbac.PasswordCredentialsUpdateParameters{
			Value: passCreds,
		}

		_, err := client.UpdatePasswordCredentials(ctx, d.Id(), passUpdate)
		if err != nil {
			return err
		}
	}

	return resourceArmAdApplicationRead(d, meta)
}

func resourceArmAdApplicationDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandAzureRmAdApplicationHomepage(d *schema.ResourceData, name string) *string {
	if v, ok := d.GetOk("homepage"); ok {
		return utils.String(v.(string))
	}

	return utils.String(fmt.Sprintf("http://%s", name))
}

func expandAzureRmAdApplicationIdentifierUris(d *schema.ResourceData, name string) *[]string {
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

func expandAzureRmAdApplicationReplyUrls(d *schema.ResourceData, name string) *[]string {
	replyUrls := d.Get("reply_urls").([]interface{})
	urls := []string{}
	for _, url := range replyUrls {
		urls = append(urls, url.(string))
	}

	return &urls
}
