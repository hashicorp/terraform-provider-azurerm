package azuread

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/graph"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/p"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_principal_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringIsEmailAddress,
			},

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"mail_nickname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"account_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 256), //currently the max length for AAD passwords is 256
			},

			"force_password_change": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"mail": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"onpremises_sam_account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"onpremises_user_principal_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"immutable_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "This must be specified if you are using a federated domain for the user's userPrincipalName (UPN) property when creating a new user account. " +
					"It is used to associate an on-premises Active Directory user account with their Azure AD user object.",
			},

			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"usage_location": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "A two letter country code (ISO standard 3166). " +
					"Required for users that will be assigned licenses due to legal requirement to check for availability of services in countries. " +
					"Examples include: `NO`, `JP`, and `GB`. Not nullable.",
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).usersClient
	ctx := meta.(*ArmClient).StopContext

	upn := d.Get("user_principal_name").(string)
	mailNickName := d.Get("mail_nickname").(string)

	//default mail nickname to the first part of the UPN (matches the portal)
	if mailNickName == "" {
		mailNickName = strings.Split(upn, "@")[0]
	}

	userCreateParameters := graphrbac.UserCreateParameters{
		AccountEnabled: p.BoolI(d.Get("account_enabled")),
		DisplayName:    p.StringI(d.Get("display_name")),
		MailNickname:   &mailNickName,
		PasswordProfile: &graphrbac.PasswordProfile{
			ForceChangePasswordNextLogin: p.BoolI(d.Get("force_password_change")),
			Password:                     p.StringI(d.Get("password")),
		},
		UserPrincipalName: &upn,
	}

	if v, ok := d.GetOk("usage_location"); ok {
		userCreateParameters.UsageLocation = p.StringI(v)
	}

	if v, ok := d.GetOk("immutable_id"); ok {
		userCreateParameters.ImmutableID = p.StringI(v)
	}

	user, err := client.Create(ctx, userCreateParameters)
	if err != nil {
		return fmt.Errorf("Error creating User (%q): %+v", upn, err)
	}
	if user.ObjectID == nil {
		return fmt.Errorf("nil User ID for %q: %+v", upn, err)
	}
	d.SetId(*user.ObjectID)

	_, err = graph.WaitForCreationReplication(func() (interface{}, error) {
		return client.Get(ctx, *user.ObjectID)
	})
	if err != nil {
		return fmt.Errorf("Error waiting for User (%s) with ObjectId %q: %+v", upn, *user.ObjectID, err)
	}

	return resourceUserRead(d, meta)
}

func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).usersClient
	ctx := meta.(*ArmClient).StopContext

	var userUpdateParameters graphrbac.UserUpdateParameters

	if d.HasChange("display_name") {
		userUpdateParameters.DisplayName = p.StringI(d.Get("display_name"))
	}

	if d.HasChange("mail_nickname") {
		userUpdateParameters.MailNickname = p.StringI(d.Get("mail_nickname"))
	}

	if d.HasChange("account_enabled") {
		userUpdateParameters.AccountEnabled = p.BoolI(d.Get("account_enabled"))
	}

	if d.HasChange("password") {
		userUpdateParameters.PasswordProfile = &graphrbac.PasswordProfile{
			ForceChangePasswordNextLogin: p.BoolI(d.Get("force_password_change")),
			Password:                     p.StringI(d.Get("password")),
		}
	}

	if d.HasChange("usage_location") {
		userUpdateParameters.UsageLocation = p.StringI(d.Get("usage_location"))
	}

	if d.HasChange("immutable_id") {
		userUpdateParameters.ImmutableID = p.StringI(d.Get("immutable_id"))
	}

	if _, err := client.Update(ctx, d.Id(), userUpdateParameters); err != nil {
		return fmt.Errorf("Error updating User with ID %q: %+v", d.Id(), err)
	}

	return resourceUserRead(d, meta)
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).usersClient
	ctx := meta.(*ArmClient).StopContext

	objectId := d.Id()

	user, err := client.Get(ctx, objectId)
	if err != nil {
		if ar.ResponseWasNotFound(user.Response) {
			log.Printf("[DEBUG] User with Object ID %q was not found - removing from state!", objectId)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving User with ID %q: %+v", objectId, err)
	}

	d.Set("user_principal_name", user.UserPrincipalName)
	d.Set("display_name", user.DisplayName)
	d.Set("mail", user.Mail)
	d.Set("mail_nickname", user.MailNickname)
	d.Set("account_enabled", user.AccountEnabled)
	d.Set("object_id", user.ObjectID)
	d.Set("usage_location", user.UsageLocation)
	d.Set("immutable_id", user.ImmutableID)

	d.Set("onpremises_sam_account_name", user.AdditionalProperties["onPremisesSamAccountName"])
	d.Set("onpremises_user_principal_name", user.AdditionalProperties["onPremisesUserPrincipalName"])

	return nil
}

func resourceUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).usersClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Delete(ctx, d.Id())
	if err != nil {
		if !ar.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error Deleting User with ID %q: %+v", d.Id(), err)
		}
	}

	return nil
}
