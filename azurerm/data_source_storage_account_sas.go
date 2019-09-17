package azurerm

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/hashicorp/go-azure-helpers/storage"
	"github.com/hashicorp/terraform/helper/schema"
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

		Schema: map[string]*schema.Schema{
			"connection_string": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"https_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"resource_types": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},

						"container": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},

						"object": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"services": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},

						"queue": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},

						"table": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},

						"file": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			// Always in UTC and must be ISO-8601 format
			"start": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ISO8601DateTime,
			},

			// Always in UTC and must be ISO-8601 format
			"expiry": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ISO8601DateTime,
			},

			"permissions": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},

						"write": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},

						"delete": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},

						"list": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},

						"add": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},

						"create": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},

						"update": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},

						"process": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
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
	resourceTypesIface := d.Get("resource_types").([]interface{})
	servicesIface := d.Get("services").([]interface{})
	start := d.Get("start").(string)
	expiry := d.Get("expiry").(string)
	permissionsIface := d.Get("permissions").([]interface{})

	resourceTypes := buildResourceTypesString(resourceTypesIface[0].(map[string]interface{}))
	services := buildServicesString(servicesIface[0].(map[string]interface{}))
	permissions := buildPermissionsString(permissionsIface[0].(map[string]interface{}))

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
	var signedProtocol = "https,http"
	if httpsOnly {
		signedProtocol = "https"
	}
	signedIp := ""
	signedVersion := sasSignedVersion

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

func buildPermissionsString(perms map[string]interface{}) string {
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

func buildServicesString(services map[string]interface{}) string {
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

func buildResourceTypesString(resTypes map[string]interface{}) string {
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
