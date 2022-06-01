package logz

import (
	"github.com/Azure/azure-sdk-for-go/services/logz/mgmt/2020-10-01/logz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func SchemaUserInfo() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"email": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"first_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringLenBetween(1, 50),
				},

				"last_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringLenBetween(1, 50),
				},

				"phone_number": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringLenBetween(1, 40),
				},
			},
		},
	}
}

func expandUserInfo(input []interface{}) *logz.UserInfo {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &logz.UserInfo{
		FirstName:    utils.String(v["first_name"].(string)),
		LastName:     utils.String(v["last_name"].(string)),
		EmailAddress: utils.String(v["email"].(string)),
		PhoneNumber:  utils.String(v["phone_number"].(string)),
	}
}

func flattenUserInfo(input *logz.UserInfo) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var firstName string
	if input.FirstName != nil {
		firstName = *input.FirstName
	}

	var lastName string
	if input.LastName != nil {
		lastName = *input.LastName
	}

	var email string
	if input.EmailAddress != nil {
		email = *input.EmailAddress
	}

	var phoneNumber string
	if input.PhoneNumber != nil {
		phoneNumber = *input.PhoneNumber
	}

	return []interface{}{
		map[string]interface{}{
			"first_name":   firstName,
			"last_name":    lastName,
			"email":        email,
			"phone_number": phoneNumber,
		},
	}
}
