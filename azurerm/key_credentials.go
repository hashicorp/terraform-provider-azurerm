package azurerm

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func keyCredentialsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key_id": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validateUUID,
				},

				"start_date": {
					Type:             schema.TypeString,
					Optional:         true,
					Computed:         true,
					DiffSuppressFunc: compareDataAsUTCSuppressFunc,
					ValidateFunc:     validateRFC3339Date,
				},

				"end_date": {
					Type:             schema.TypeString,
					Optional:         true,
					Computed:         true,
					DiffSuppressFunc: compareDataAsUTCSuppressFunc,
					ValidateFunc:     validateRFC3339Date,
				},

				"type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"AsymmetricX509Cert",
						"Symmetric",
					}, true),
				},

				"usage": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"Verify",
						"Sign",
					}, true),
				},

				"value": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
			},
		},
		Set: resourceKeyCredentialHash,
	}
}

func resourceKeyCredentialHash(v interface{}) int {
	var buf bytes.Buffer

	if v != nil {
		m := v.(map[string]interface{})
		buf.WriteString(fmt.Sprintf("%s-", m["key_id"].(string)))
	}

	return hashcode.String(buf.String())
}

func flattenAzureRmKeyCredentials(creds *[]graphrbac.KeyCredential) []interface{} {
	result := make([]interface{}, 0, len(*creds))
	for _, cred := range *creds {
		l := make(map[string]interface{})
		l["key_id"] = *cred.KeyID
		l["start_date"] = string((*cred.StartDate).Format(time.RFC3339))
		l["end_date"] = string((*cred.EndDate).Format(time.RFC3339))
		l["type"] = *cred.Type
		l["usage"] = *cred.Usage

		result = append(result, l)
	}
	return result
}

func expandAzureRmKeyCredentials(d *schema.ResourceData) (*[]graphrbac.KeyCredential, error) {
	creds := d.Get("key_credential").(*schema.Set).List()
	keyCreds := make([]graphrbac.KeyCredential, 0, len(creds))
	for _, credsConfig := range creds {
		config := credsConfig.(map[string]interface{})

		keyId := config["key_id"].(string)
		startDate := config["start_date"].(string)
		endDate := config["end_date"].(string)
		keyType := config["type"].(string)
		usage := config["usage"].(string)
		value := config["value"].(string)

		if keyType == "AsymmetricX509Cert" && usage == "Sign" {
			return nil, fmt.Errorf("Usage cannot be set to %s when %s is set as type for a Key Credential", usage, keyType)
		}

		keyCred := graphrbac.KeyCredential{
			KeyID: &keyId,
			Type:  &keyType,
			Usage: &usage,
			Value: &value,
		}

		if startDate != "" {
			starttime, sterr := time.Parse(time.RFC3339, startDate)
			if sterr != nil {
				return nil, fmt.Errorf("Cannot parse start_date: %q", startDate)
			}
			stdt := date.Time{Time: starttime}
			keyCred.StartDate = &stdt
		}

		if endDate != "" {
			endtime, eterr := time.Parse(time.RFC3339, endDate)
			if eterr != nil {
				return nil, fmt.Errorf("Cannot parse end_date: %q", endDate)
			}
			etdt := date.Time{Time: endtime}
			keyCred.EndDate = &etdt
		}

		keyCreds = append(keyCreds, keyCred)
	}

	return &keyCreds, nil
}
