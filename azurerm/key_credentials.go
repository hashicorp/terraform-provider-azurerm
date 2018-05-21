package azurerm

import (
	"bytes"
	"fmt"
	"regexp"
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
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key_id": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validateUUID,
				},

				"start_date": {
					Type:             schema.TypeString,
					Required:         true,
					DiffSuppressFunc: compareDataAsUTCSuppressFunc,
					ValidateFunc:     validateRFC3339Date,
				},

				"end_date": {
					Type:             schema.TypeString,
					Required:         true,
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

func customizeDiffKeyCredential(diff *schema.ResourceDiff, v interface{}) error {
	o, n := diff.GetChange("key_credential")
	if o == nil {
		o = new(schema.Set)
	}
	if n == nil {
		n = new(schema.Set)
	}
	os := o.(*schema.Set)
	ns := n.(*schema.Set)

	// Detect if the user changed a key credential property without changing the
	// KeyID associated with the credential. Changing the property changes the hash,
	// which causes Terraform to pick it up as a new Set item. This will cause
	// a conflict due to the unique KeyID requirement.
	m := make(map[string]string)
	for _, v := range os.Difference(ns).List() {
		x := v.(map[string]interface{})
		m[x["key_id"].(string)] = x["key_id"].(string)
	}
	for _, v := range ns.Difference(os).List() {
		x := v.(map[string]interface{})
		if _, ok := m[x["key_id"].(string)]; ok {
			return fmt.Errorf("Error: changing Key Credential properties on existing KeyID %s requires generating a new unique KeyID.", x["key_id"].(string))
		}
	}

	return nil
}

func resourceKeyCredentialHash(v interface{}) int {
	var buf bytes.Buffer

	if v != nil {
		m := v.(map[string]interface{})
		buf.WriteString(fmt.Sprintf("%s-", m["key_id"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["usage"].(string)))

		// We parse the DateTimes and then convert them back to a string
		// in order to have a consistent format for the hash.

		if st, err := time.Parse(time.RFC3339, m["start_date"].(string)); err == nil {
			buf.WriteString(fmt.Sprintf("%s-", string((st).Format(time.RFC3339))))
		}

		if et, err := time.Parse(time.RFC3339, m["end_date"].(string)); err == nil {
			buf.WriteString(fmt.Sprintf("%s-", string((et).Format(time.RFC3339))))
		}
	}

	return hashcode.String(buf.String())
}

func flattenAzureRmKeyCredential(cred *graphrbac.KeyCredential) interface{} {
	l := make(map[string]interface{})
	l["key_id"] = *cred.KeyID
	l["type"] = *cred.Type
	l["usage"] = *cred.Usage
	l["start_date"] = string((*cred.StartDate).Format(time.RFC3339))
	l["end_date"] = string((*cred.EndDate).Format(time.RFC3339))

	return l
}

func flattenAzureRmKeyCredentials(creds *[]graphrbac.KeyCredential) []interface{} {
	result := make([]interface{}, 0, len(*creds))
	for _, cred := range *creds {
		result = append(result, flattenAzureRmKeyCredential(&cred))
	}
	return result
}

func expandAzureRmKeyCredential(d *map[string]interface{}) (*graphrbac.KeyCredential, error) {
	keyId := (*d)["key_id"].(string)
	startDate := (*d)["start_date"].(string)
	endDate := (*d)["end_date"].(string)
	keyType := (*d)["type"].(string)
	usage := (*d)["usage"].(string)
	value := (*d)["value"].(string)

	if keyType == "AsymmetricX509Cert" && usage == "Sign" {
		return nil, fmt.Errorf("Usage cannot be set to %s when %s is set as type for a Key Credential", usage, keyType)
	}

	// Match against the prefix/suffix and '\n' of certificate values.
	// The API doesn't accept them so they need to be removed.
	pkre := regexp.MustCompile(`(-{5}.+?-{5})|(\n)`)
	value = pkre.ReplaceAllString(value, ``)

	kc := graphrbac.KeyCredential{
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
		stdt := date.Time{Time: starttime.Truncate(time.Second)}
		kc.StartDate = &stdt
	}

	if endDate != "" {
		endtime, eterr := time.Parse(time.RFC3339, endDate)
		if eterr != nil {
			return nil, fmt.Errorf("Cannot parse end_date: %q", endDate)
		}
		etdt := date.Time{Time: endtime.Truncate(time.Second)}
		kc.EndDate = &etdt
	}

	return &kc, nil
}

func expandAzureRmKeyCredentials(d *schema.ResourceData, o *schema.Set) (*[]graphrbac.KeyCredential, error) {
	creds := d.Get("key_credential").(*schema.Set).List()
	keyCreds := make([]graphrbac.KeyCredential, 0, len(creds))

	for _, v := range creds {
		cfg := v.(map[string]interface{})
		cred, err := expandAzureRmKeyCredential(&cfg)
		if err != nil {
			return nil, err
		}

		// Azure only allows an in-place update of the Key Credentials list.
		// Existing keys, matched by their KeyID, must be sent back with their
		// Value attribute set to nil. New keys need to provide a Value.
		// By referencing the existing schema (o), we can determine which
		// entries in the list are existing keys.
		if o != nil && o.Contains(v) {
			cred.Value = nil
		}

		keyCreds = append(keyCreds, *cred)
	}

	return &keyCreds, nil
}
