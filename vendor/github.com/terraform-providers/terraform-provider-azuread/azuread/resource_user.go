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

			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).usersClient
	ctx := meta.(*ArmClient).StopContext

	upn := d.Get("user_principal_name").(string)
	displayName := d.Get("display_name").(string)
	mailNickName := d.Get("mail_nickname").(string)
	accountEnabled := d.Get("account_enabled").(bool)
	password := d.Get("password").(string)
	forcePasswordChange := d.Get("force_password_change").(bool)

	//default mail nickname to the first part of the UPN (matches the portal)
	if mailNickName == "" {
		mailNickName = strings.Split(upn, "@")[0]
	}

	userCreateParameters := graphrbac.UserCreateParameters{
		AccountEnabled: &accountEnabled,
		DisplayName:    &displayName,
		MailNickname:   &mailNickName,
		PasswordProfile: &graphrbac.PasswordProfile{
			ForceChangePasswordNextLogin: &forcePasswordChange,
			Password:                     &password,
		},
		UserPrincipalName: &upn,
	}

	user, err := client.Create(ctx, userCreateParameters)
	if err != nil {
		return fmt.Errorf("Error creating User (%q): %+v", upn, err)
	}
	if user.ObjectID == nil {
		return fmt.Errorf("nil User ID for %q: %+v", upn, err)
	}
	d.SetId(*user.ObjectID)

	_, err = graph.WaitForReplication(func() (interface{}, error) {
		return client.Get(ctx, *user.ObjectID)
	})
	if err != nil {
		return fmt.Errorf("Error waiting for User (%s) with ObjectId %q: %+v", upn, *user.ObjectID, err)
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
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).usersClient
	ctx := meta.(*ArmClient).StopContext

	var userUpdateParameters graphrbac.UserUpdateParameters

	if d.HasChange("display_name") {
		displayName := d.Get("display_name").(string)
		userUpdateParameters.DisplayName = p.String(displayName)
	}

	if d.HasChange("mail_nickname") {
		mailNickName := d.Get("mail_nickname").(string)
		userUpdateParameters.MailNickname = p.String(mailNickName)
	}

	if d.HasChange("account_enabled") {
		accountEnabled := d.Get("account_enabled").(bool)
		userUpdateParameters.AccountEnabled = p.Bool(accountEnabled)
	}

	if d.HasChange("password") {
		password := d.Get("password").(string)
		forcePasswordChange := d.Get("force_password_change").(bool)

		passwordProfile := &graphrbac.PasswordProfile{
			ForceChangePasswordNextLogin: &forcePasswordChange,
			Password:                     &password,
		}

		userUpdateParameters.PasswordProfile = passwordProfile
	}

	if _, err := client.Update(ctx, d.Id(), userUpdateParameters); err != nil {
		return fmt.Errorf("Error updating User with ID %q: %+v", d.Id(), err)
	}

	return resourceUserRead(d, meta)
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
