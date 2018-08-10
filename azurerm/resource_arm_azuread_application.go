package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"available_to_other_tenants": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"oauth2_allow_implicit_flow": {
				Type:     schema.TypeBool,
				Optional: true,
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

	// NOTE: name isn't the Resource ID here, so we don't check it exists
	name := d.Get("name").(string)
	availableToOtherTenants := d.Get("available_to_other_tenants").(bool)
	properties := graphrbac.ApplicationCreateParameters{
		DisplayName:             &name,
		Homepage:                expandAzureRmActiveDirectoryApplicationHomepage(d, name),
		IdentifierUris:          expandAzureRmActiveDirectoryApplicationIdentifierUris(d),
		ReplyUrls:               expandAzureRmActiveDirectoryApplicationReplyUrls(d),
		AvailableToOtherTenants: utils.Bool(availableToOtherTenants),
	}

	if v, ok := d.GetOk("oauth2_allow_implicit_flow"); ok {
		properties.Oauth2AllowImplicitFlow = utils.Bool(v.(bool))
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	app, err := client.Create(waitCtx, properties)
	if err != nil {
		return err
	}

	d.SetId(*app.ObjectID)

	return resourceArmActiveDirectoryApplicationRead(d, meta)
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
		properties.IdentifierUris = expandAzureRmActiveDirectoryApplicationIdentifierUris(d)
	}

	if d.HasChange("reply_urls") {
		properties.ReplyUrls = expandAzureRmActiveDirectoryApplicationReplyUrls(d)
	}

	if d.HasChange("available_to_other_tenants") {
		availableToOtherTenants := d.Get("available_to_other_tenants").(bool)
		properties.AvailableToOtherTenants = utils.Bool(availableToOtherTenants)
	}

	if d.HasChange("oauth2_allow_implicit_flow") {
		oauth := d.Get("oauth2_allow_implicit_flow").(bool)
		properties.Oauth2AllowImplicitFlow = utils.Bool(oauth)
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	_, err := client.Patch(waitCtx, d.Id(), properties)
	if err != nil {
		return fmt.Errorf("Error patching Azure AD Application with ID %q: %+v", d.Id(), err)
	}

	return resourceArmActiveDirectoryApplicationRead(d, meta)
}

func resourceArmActiveDirectoryApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, d.Id())
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Azure AD Application with ID %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Azure AD Application with ID %q: %+v", d.Id(), err)
	}

	d.Set("name", resp.DisplayName)
	d.Set("application_id", resp.AppID)
	d.Set("homepage", resp.Homepage)
	d.Set("available_to_other_tenants", resp.AvailableToOtherTenants)
	d.Set("oauth2_allow_implicit_flow", resp.Oauth2AllowImplicitFlow)

	identifierUris := flattenAzureADApplicationIdentifierUris(resp.IdentifierUris)
	if err := d.Set("identifier_uris", identifierUris); err != nil {
		return fmt.Errorf("Error setting `identifier_uris`: %+v", err)
	}

	replyUrls := flattenAzureADApplicationReplyUrls(resp.ReplyUrls)
	if err := d.Set("reply_urls", replyUrls); err != nil {
		return fmt.Errorf("Error setting `reply_urls`: %+v", err)
	}

	return nil
}

func resourceArmActiveDirectoryApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()

	// in order to delete an application which is available to other tenants, we first have to disable this setting
	availableToOtherTenants := d.Get("available_to_other_tenants").(bool)
	if availableToOtherTenants {
		log.Printf("[DEBUG] Azure AD Application is available to other tenants - disabling that feature before deleting.")
		properties := graphrbac.ApplicationUpdateParameters{
			AvailableToOtherTenants: utils.Bool(false),
		}
		_, err := client.Patch(waitCtx, d.Id(), properties)
		if err != nil {
			return fmt.Errorf("Error patching Azure AD Application with ID %q: %+v", d.Id(), err)
		}
	}

	resp, err := client.Delete(waitCtx, d.Id())
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error Deleting Azure AD Application with ID %q: %+v", d.Id(), err)
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

func expandAzureRmActiveDirectoryApplicationIdentifierUris(d *schema.ResourceData) *[]string {
	identifierUris := d.Get("identifier_uris").([]interface{})
	identifiers := make([]string, 0)

	for _, id := range identifierUris {
		identifiers = append(identifiers, id.(string))
	}

	return &identifiers
}

func expandAzureRmActiveDirectoryApplicationReplyUrls(d *schema.ResourceData) *[]string {
	replyUrls := d.Get("reply_urls").([]interface{})
	urls := make([]string, 0)

	for _, url := range replyUrls {
		urls = append(urls, url.(string))
	}

	return &urls
}

func flattenAzureADApplicationIdentifierUris(input *[]string) []string {
	output := make([]string, 0)

	if input != nil {
		for _, v := range *input {
			output = append(output, v)
		}
	}

	return output
}

func flattenAzureADApplicationReplyUrls(input *[]string) []string {
	output := make([]string, 0)

	if input != nil {
		for _, v := range *input {
			output = append(output, v)
		}
	}

	return output
}
