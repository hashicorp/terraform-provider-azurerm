package graph

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/p"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
)

// valid types are `application` and `service_principal`
func PasswordResourceSchema(object_type string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		object_type + "_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.UUID,
		},

		"key_id": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validate.UUID,
		},

		"value": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringLenBetween(1, 863), // Encrypted secret cannot be empty and can be at most 1024 bytes.
		},

		"start_date": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"end_date": {
			Type:          schema.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"end_date_relative"},
			ValidateFunc:  validation.IsRFC3339Time,
		},

		"end_date_relative": {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"end_date"},
			ValidateFunc:  validate.NoEmptyStrings,
		},
	}
}

type PasswordCredentialId struct {
	ObjectId string
	KeyId    string
}

func (id PasswordCredentialId) String() string {
	return id.ObjectId + "/" + id.KeyId
}

func ParsePasswordCredentialId(id string) (PasswordCredentialId, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return PasswordCredentialId{}, fmt.Errorf("Password Credential ID should be in the format {objectId}/{keyId} - but got %q", id)
	}

	if _, err := uuid.ParseUUID(parts[0]); err != nil {
		return PasswordCredentialId{}, fmt.Errorf("Object ID isn't a valid UUID (%q): %+v", id[0], err)
	}

	if _, err := uuid.ParseUUID(parts[1]); err != nil {
		return PasswordCredentialId{}, fmt.Errorf("Credential ID isn't a valid UUID (%q): %+v", id[1], err)
	}

	return PasswordCredentialId{
		ObjectId: parts[0],
		KeyId:    parts[1],
	}, nil
}

func PasswordCredentialIdFrom(objectId, keyId string) PasswordCredentialId {
	return PasswordCredentialId{
		ObjectId: objectId,
		KeyId:    keyId,
	}
}

func PasswordCredentialForResource(d *schema.ResourceData) (*graphrbac.PasswordCredential, error) {
	value := d.Get("value").(string)

	// errors should be handled by the validation
	var keyId string
	if v, ok := d.GetOk("key_id"); ok {
		keyId = v.(string)
	} else {
		kid, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}

		keyId = kid
	}

	var endDate time.Time
	if v := d.Get("end_date").(string); v != "" {
		endDate, _ = time.Parse(time.RFC3339, v)
	} else if v := d.Get("end_date_relative").(string); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("unable to parse `end_date_relative` (%s) as a duration", v)
		}
		endDate = time.Now().Add(d)
	} else {
		return nil, fmt.Errorf("one of `end_date` or `end_date_relative` must be specified")
	}

	credential := graphrbac.PasswordCredential{
		KeyID:   p.String(keyId),
		Value:   p.String(value),
		EndDate: &date.Time{Time: endDate},
	}

	if v, ok := d.GetOk("start_date"); ok {
		// errors will be handled by the validation
		startDate, _ := time.Parse(time.RFC3339, v.(string))
		credential.StartDate = &date.Time{Time: startDate}
	}

	return &credential, nil
}

func PasswordCredentialResultFindByKeyId(creds graphrbac.PasswordCredentialListResult, keyId string) *graphrbac.PasswordCredential {
	var cred *graphrbac.PasswordCredential

	if creds.Value != nil {
		for _, c := range *creds.Value {
			if c.KeyID == nil {
				continue
			}

			if *c.KeyID == keyId {
				cred = &c
				break
			}
		}
	}

	return cred
}

func PasswordCredentialResultAdd(existing graphrbac.PasswordCredentialListResult, cred *graphrbac.PasswordCredential, errorOnDuplicate bool) (*[]graphrbac.PasswordCredential, error) {
	newCreds := make([]graphrbac.PasswordCredential, 0)

	if existing.Value != nil {
		if errorOnDuplicate {
			for _, v := range *existing.Value {
				if v.KeyID == nil {
					continue
				}

				if *v.KeyID == *cred.KeyID {
					return nil, fmt.Errorf("credential already exists found")
				}
			}
		}

		newCreds = *existing.Value
	}
	newCreds = append(newCreds, *cred)

	return &newCreds, nil
}

func PasswordCredentialResultRemoveByKeyId(existing graphrbac.PasswordCredentialListResult, keyId string) *[]graphrbac.PasswordCredential {
	newCreds := make([]graphrbac.PasswordCredential, 0)

	if existing.Value != nil {
		for _, v := range *existing.Value {
			if v.KeyID == nil {
				continue
			}

			if *v.KeyID == keyId {
				continue
			}

			newCreds = append(newCreds, v)
		}
	}

	return &newCreds
}

func WaitForPasswordCredentialReplication(keyId string, f func() (graphrbac.PasswordCredentialListResult, error)) (interface{}, error) {
	return (&resource.StateChangeConf{
		Pending:                   []string{"404", "BadCast", "NotFound"},
		Target:                    []string{"Found"},
		Timeout:                   5 * time.Minute,
		MinTimeout:                1 * time.Second,
		ContinuousTargetOccurence: 10,
		Refresh: func() (interface{}, string, error) {
			creds, err := f()
			if err != nil {
				if ar.ResponseWasNotFound(creds.Response) {
					return creds, "404", nil
				}
				return creds, "Error", fmt.Errorf("Error calling f, response was not 404 (%d): %v", creds.Response.StatusCode, err)
			}

			credential := PasswordCredentialResultFindByKeyId(creds, keyId)
			if credential == nil {
				return creds, "NotFound", nil
			}

			return creds, "Found", nil
		},
	}).WaitForState()
}
