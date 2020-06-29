package azuread

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/graph"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
)

func dataUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUsersRead,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"object_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MinItems:      1,
				ConflictsWith: []string{"user_principal_names"},
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.UUID,
				},
			},

			"user_principal_names": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MinItems:      1,
				ConflictsWith: []string{"object_ids"},
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.NoEmptyStrings,
				},
			},

			"mail_nicknames": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MinItems:      1,
				ConflictsWith: []string{"object_ids", "user_principal_names"},
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.NoEmptyStrings,
				},
			},
		},
	}
}

func dataSourceUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).usersClient
	ctx := meta.(*ArmClient).StopContext

	var users []graphrbac.User
	expectedCount := 0

	if upns, ok := d.Get("user_principal_names").([]interface{}); ok && len(upns) > 0 {
		expectedCount = len(upns)
		for _, v := range upns {
			resp, err := client.Get(ctx, v.(string))
			if err != nil {
				return fmt.Errorf("Error making Read request on AzureAD User with ID %q: %+v", v.(string), err)
			}

			users = append(users, resp)
		}
	} else if oids, ok := d.Get("object_ids").([]interface{}); ok && len(oids) > 0 {
		expectedCount = len(oids)
		for _, v := range oids {
			u, err := graph.UserGetByObjectId(&client, ctx, v.(string))
			if err != nil {
				return fmt.Errorf("Error finding Azure AD User with object ID %q: %+v", v.(string), err)
			}
			users = append(users, *u)
		}
	} else if mailNicknames, ok := d.Get("mail_nicknames").([]interface{}); ok && len(mailNicknames) > 0 {
		expectedCount = len(mailNicknames)
		for _, v := range mailNicknames {
			u, err := graph.UserGetByMailNickname(&client, ctx, v.(string))
			if err != nil {
				return fmt.Errorf("Error finding Azure AD User with email alias %q: %+v", v.(string), err)
			}
			users = append(users, *u)
		}
	} else {
		return fmt.Errorf("one of `object_ids`, `user_principal_names` or `mail_nicknames` must be supplied")
	}

	if len(users) != expectedCount {
		return fmt.Errorf("Unexpected number of users returned (%d != %d)", len(users), expectedCount)
	}

	upns := make([]string, 0, len(users))
	oids := make([]string, 0, len(users))
	mailNicknames := make([]string, 0, len(users))
	for _, u := range users {
		if u.ObjectID == nil || u.UserPrincipalName == nil {
			return fmt.Errorf("User with nil ObjectId or UPN was found: %v", u)
		}

		oids = append(oids, *u.ObjectID)
		upns = append(upns, *u.UserPrincipalName)
		mailNicknames = append(mailNicknames, *u.MailNickname)
	}

	h := sha1.New()
	if _, err := h.Write([]byte(strings.Join(upns, "-"))); err != nil {
		return fmt.Errorf("Unable to compute hash for upns: %v", err)
	}

	d.SetId("users#" + base64.URLEncoding.EncodeToString(h.Sum(nil)))
	d.Set("object_ids", oids)
	d.Set("user_principal_names", upns)
	d.Set("mail_nicknames", mailNicknames)
	return nil
}
