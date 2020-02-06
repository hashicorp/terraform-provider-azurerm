package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/hashicorp/go-azure-helpers/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func dataSourceArmStorageAccountBlobContainerSharedAccessSignature() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmStorageContainerSasRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"connection_string": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"container_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"https_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.SharedAccessSignatureIP,
			},

			"start": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ISO8601DateTime,
			},

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

						"add": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"create": {
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
					},
				},
			},

			"cache_control": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_disposition": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_language": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"sas": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceArmStorageContainerSasRead(d *schema.ResourceData, _ interface{}) error {
	connString := d.Get("connection_string").(string)
	containerName := d.Get("container_name").(string)
	httpsOnly := d.Get("https_only").(bool)
	ip := d.Get("ip_address").(string)
	start := d.Get("start").(string)
	expiry := d.Get("expiry").(string)
	permissionsIface := d.Get("permissions").([]interface{})

	// response headers
	cacheControl := d.Get("cache_control").(string)
	contentDisposition := d.Get("content_disposition").(string)
	contentEncoding := d.Get("content_encoding").(string)
	contentLanguage := d.Get("content_language").(string)
	contentType := d.Get("content_type").(string)

	permissions := BuildContainerPermissionsString(permissionsIface[0].(map[string]interface{}))

	// Parse the connection string
	kvp, err := storage.ParseAccountSASConnectionString(connString)
	if err != nil {
		return err
	}

	// Create the string to sign with the key...
	accountName := kvp[connStringAccountNameKey]
	accountKey := kvp[connStringAccountKeyKey]
	var signedProtocol = "https,http"
	if httpsOnly {
		signedProtocol = "https"
	}
	signedIp := ip
	signedIdentifier := ""
	signedSnapshotTime := ""

	sasToken, err := storage.ComputeContainerSASToken(permissions, start, expiry, accountName, accountKey,
		containerName, signedIdentifier, signedIp, signedProtocol, signedSnapshotTime, cacheControl,
		contentDisposition, contentEncoding, contentLanguage, contentType)
	if err != nil {
		return err
	}

	d.Set("sas", sasToken)
	tokenHash := sha256.Sum256([]byte(sasToken))
	d.SetId(hex.EncodeToString(tokenHash[:]))

	return nil
}

func BuildContainerPermissionsString(perms map[string]interface{}) string {
	retVal := ""

	if val, pres := perms["read"].(bool); pres && val {
		retVal += "r"
	}

	if val, pres := perms["add"].(bool); pres && val {
		retVal += "a"
	}

	if val, pres := perms["create"].(bool); pres && val {
		retVal += "c"
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

	return retVal
}
