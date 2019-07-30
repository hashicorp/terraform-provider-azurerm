package azuread

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validate.UUID,
				ConflictsWith: []string{"user_principal_name"},
			},

			"user_principal_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validate.NoEmptyStrings,
				ConflictsWith: []string{"object_id"},
			},

			"account_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mail": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mail_nickname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).usersClient
	ctx := meta.(*ArmClient).StopContext

	var user graphrbac.User

	if upn, ok := d.Get("user_principal_name").(string); ok && upn != "" {

		// use the object_id to find the Azure AD application
		resp, err := client.Get(ctx, upn)
		if err != nil {
			if ar.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error: AzureAD User with ID %q was not found", upn)
			}

			return fmt.Errorf("Error making Read request on AzureAD User with ID %q: %+v", upn, err)
		}

		user = resp
	} else if oId, ok := d.Get("object_id").(string); ok && oId != "" {
		filter := fmt.Sprintf("objectId eq '%s'", oId)

		resp, err := client.ListComplete(ctx, filter)
		if err != nil {
			return fmt.Errorf("Error listing Azure AD Users for filter %q: %+v", filter, err)
		}

		values := resp.Response().Value
		if values == nil {
			return fmt.Errorf("nil values for AD Users matching %q", filter)
		}
		if len(*values) == 0 {
			return fmt.Errorf("Found no AD Users matching %q", filter)
		}
		if len(*values) > 2 {
			return fmt.Errorf("Found multiple AD Users matching %q", filter)
		}

		user = (*values)[0]
		if user.DisplayName == nil {
			return fmt.Errorf("nil DisplayName for AD Users matching %q", filter)
		}
		if *user.ObjectID != oId {
			return fmt.Errorf("objectID for AD Users matching %q does is does not match(%q!=%q)", filter, *user.ObjectID, oId)
		}
	} else {
		return fmt.Errorf("one of `object_id` or `user_principal_name` must be supplied")
	}

	if user.ObjectID == nil {
		return fmt.Errorf("Group objectId is nil")
	}
	d.SetId(*user.ObjectID)

	d.SetId(*user.ObjectID)
	d.Set("object_id", user.ObjectID)
	d.Set("user_principal_name", user.UserPrincipalName)
	d.Set("account_enabled", user.AccountEnabled)
	d.Set("display_name", user.DisplayName)
	d.Set("mail", user.Mail)
	d.Set("mail_nickname", user.MailNickname)

	return nil
}
