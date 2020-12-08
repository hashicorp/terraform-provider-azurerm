package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/hashicorp/go-azure-helpers/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

const (
	connStringAccountKeyKey  = "AccountKey"
	connStringAccountNameKey = "AccountName"
	sasSignedVersion         = "2017-07-29"
)

// This is an ACCOUNT SAS : https://docs.microsoft.com/en-us/rest/api/storageservices/Constructing-an-Account-SAS
// not Service SAS
func dataSourceArmStorageAccountSharedAccessSignature() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmStorageAccountSasRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"connection_string": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"https_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"signed_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  sasSignedVersion,
			},

			"resource_types": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"container": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"object": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			"services": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"queue": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"table": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"file": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			// Always in UTC and must be ISO-8601 format
			"start": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ISO8601DateTime,
			},

			// Always in UTC and must be ISO-8601 format
			"expiry": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ISO8601DateTime,
			},

			"permissions": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"write": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"delete": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"list": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"add": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"create": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"update": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"process": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			"sas": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceArmStorageAccountSasRead(d *schema.ResourceData, _ interface{}) error {
	connString := d.Get("connection_string").(string)
	httpsOnly := d.Get("https_only").(bool)
	signedVersion := d.Get("signed_version").(string)
	resourceTypesIface := d.Get("resource_types").([]interface{})
	servicesIface := d.Get("services").([]interface{})
	start := d.Get("start").(string)
	expiry := d.Get("expiry").(string)
	permissionsIface := d.Get("permissions").([]interface{})

	resourceTypes := BuildResourceTypesString(resourceTypesIface[0].(map[string]interface{}))
	services := BuildServicesString(servicesIface[0].(map[string]interface{}))
	permissions := BuildPermissionsString(permissionsIface[0].(map[string]interface{}))

	// Parse the connection string
	kvp, err := storage.ParseAccountSASConnectionString(connString)
	if err != nil {
		return err
	}

	// Create the string to sign with the key...

	// Details on how to do this are here:
	// https://docs.microsoft.com/en-us/rest/api/storageservices/Constructing-an-Account-SAS
	accountName := kvp[connStringAccountNameKey]
	accountKey := kvp[connStringAccountKeyKey]
	signedProtocol := "https,http"
	if httpsOnly {
		signedProtocol = "https"
	}
	signedIp := ""

	sasToken, err := storage.ComputeAccountSASToken(accountName, accountKey, permissions, services, resourceTypes,
		start, expiry, signedProtocol, signedIp, signedVersion)
	if err != nil {
		return err
	}

	d.Set("sas", sasToken)
	tokenHash := sha256.Sum256([]byte(sasToken))
	d.SetId(hex.EncodeToString(tokenHash[:]))

	return nil
}

func BuildPermissionsString(perms map[string]interface{}) string {
	retVal := ""

	if val, pres := perms["read"].(bool); pres && val {
		retVal += "r"
	}

	if val, pres := perms["write"].(bool); pres && val {
		retVal += "w"
	}

	if val, pres := perms["delete"].(bool); pres && val {
		retVal += "d"
	}

	if val, pres := perms["list"].(bool); pres && val {
		retVal += "l"
	}

	if val, pres := perms["add"].(bool); pres && val {
		retVal += "a"
	}

	if val, pres := perms["create"].(bool); pres && val {
		retVal += "c"
	}

	if val, pres := perms["update"].(bool); pres && val {
		retVal += "u"
	}

	if val, pres := perms["process"].(bool); pres && val {
		retVal += "p"
	}

	return retVal
}

func BuildServicesString(services map[string]interface{}) string {
	retVal := ""

	if val, pres := services["blob"].(bool); pres && val {
		retVal += "b"
	}

	if val, pres := services["queue"].(bool); pres && val {
		retVal += "q"
	}

	if val, pres := services["table"].(bool); pres && val {
		retVal += "t"
	}

	if val, pres := services["file"].(bool); pres && val {
		retVal += "f"
	}

	return retVal
}

func BuildResourceTypesString(resTypes map[string]interface{}) string {
	retVal := ""

	if val, pres := resTypes["service"].(bool); pres && val {
		retVal += "s"
	}

	if val, pres := resTypes["container"].(bool); pres && val {
		retVal += "c"
	}

	if val, pres := resTypes["object"].(bool); pres && val {
		retVal += "o"
	}

	return retVal
}
