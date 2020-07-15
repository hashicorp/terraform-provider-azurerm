package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/nickmhankins/go-azure-helpers/storage"
	// "github.com/hashicorp/go-azure-helpers/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func dataSourceArmStorageAccountTableSharedAccessSignature() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmStorageTableSasRead,

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

			"table_name": {
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

						"update": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"delete": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			"start_partition_key": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"end_partition_key": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"start_row_key": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"end_row_key": {
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

func dataSourceArmStorageTableSasRead(d *schema.ResourceData, _ interface{}) error {
	connString := d.Get("connection_string").(string)
	tableName := d.Get("table_name").(string)
	httpsOnly := d.Get("https_only").(bool)
	ip := d.Get("ip_address").(string)
	start := d.Get("start").(string)
	expiry := d.Get("expiry").(string)
	permissionsIface := d.Get("permissions").([]interface{})

	startPartitionKey := d.Get("start_partition_key").(string)
	startRowKey := d.Get("start_row_key").(string)
	endPartitionKey := d.Get("end_partition_key").(string)
	endRowKey := d.Get("end_row_key").(string)

	permissions := BuildTablePermissionsString(permissionsIface[0].(map[string]interface{}))

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

	sasToken, err := storage.ComputeTableSASToken(permissions, start, expiry, accountName, accountKey,
		tableName, signedIdentifier, signedIp, signedProtocol, startPartitionKey, startRowKey,
		endPartitionKey, endRowKey)
	if err != nil {
		return err
	}

	d.Set("sas", sasToken)
	tokenHash := sha256.Sum256([]byte(sasToken))
	d.SetId(hex.EncodeToString(tokenHash[:]))

	return nil
}

func BuildTablePermissionsString(perms map[string]interface{}) string {
	retVal := ""

	if val, pres := perms["read"].(bool); pres && val {
		retVal += "r"
	}

	if val, pres := perms["add"].(bool); pres && val {
		retVal += "a"
	}

	if val, pres := perms["update"].(bool); pres && val {
		retVal += "u"
	}

	if val, pres := perms["delete"].(bool); pres && val {
		retVal += "d"
	}

	return retVal
}
