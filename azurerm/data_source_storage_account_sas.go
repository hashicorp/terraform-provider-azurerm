package azurerm

import "github.com/hashicorp/terraform/helper/schema"

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"crypto/hmac"
	"crypto/sha256"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"https_only": {
				Type:     schema.TypeBool,
				Required: true,
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
				ForceNew: true,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// Always in UTC and must be ISO-8601 format
			"expiry": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceArmStorageAccountSasRead(d *schema.ResourceData, meta interface{}) error {

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

	_, svcErr := validateArmStorageAccountSasResourceTypes(services, "")
	if svcErr != nil {
		return svcErr[0]
	}

	// Parse the connection string
	kvp, err := parseAzureStorageAccountConnectionString(connString)
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

	sasToken, err := computeAzureStorageAccountSas(accountName, accountKey, permissions, services, resourceTypes,
		start, expiry, signedProtocol, signedIp, signedVersion)
	if err != nil {
		return err
	}

	d.Set("sas", sasToken)

	return nil
}

func buildPermissionsString(perms map[string]interface{}) string {
	retVal := ""

	if val := perms["read"].(bool); val {
		retVal += "r"
	}

	if val := perms["write"].(bool); val {
		retVal += "w"
	}

	if val := perms["delete"].(bool); val {
		retVal += "d"
	}

	if val := perms["list"].(bool); val {
		retVal += "l"
	}

	if val := perms["add"].(bool); val {
		retVal += "a"
	}

	if val := perms["create"].(bool); val {
		retVal += "c"
	}

	if val := perms["update"].(bool); val {
		retVal += "u"
	}

	if val := perms["process"].(bool); val {
		retVal += "p"
	}

	return retVal
}

func buildServicesString(services map[string]interface{}) string {
	retVal := ""

	if val := services["blob"].(bool); val {
		retVal += "b"
	}

	if val := services["queue"].(bool); val {
		retVal += "q"
	}

	if val := services["table"].(bool); val {
		retVal += "t"
	}

	if val := services["file"].(bool); val {
		retVal += "f"
	}

	return retVal
}

func buildResourceTypesString(resTypes map[string]interface{}) string {
	retVal := ""

	if val := resTypes["service"].(bool); val {
		retVal += "s"
	}

	if val := resTypes["container"].(bool); val {
		retVal += "c"
	}

	if val := resTypes["object"].(bool); val {
		retVal += "o"
	}

	return retVal
}

func validateArmStorageAccountSasResourceTypes(v interface{}, k string) (ws []string, es []error) {
	input := v.(string)

	if !regexp.MustCompile(`\A([cos]{1,3})\z`).MatchString(input) {
		es = append(es, fmt.Errorf("resource Types can only consist of 's', 'c', 'o', and must be between 1 and 3 characters long"))
	}

	return
}

func computeAzureStorageAccountSas(accountName string,
	accountKey string,
	permissions string,
	services string,
	resourceTypes string,
	start string,
	expiry string,
	signedProtocol string,
	signedIp string,
	signedVersion string) (string, error) {

	// UTF-8 by default...
	stringToSign := accountName + "\n"
	stringToSign += permissions + "\n"
	stringToSign += services + "\n"
	stringToSign += resourceTypes + "\n"
	stringToSign += start + "\n"
	stringToSign += expiry + "\n"
	stringToSign += signedIp + "\n"
	stringToSign += signedProtocol + "\n"
	stringToSign += signedVersion + "\n"

	binaryKey, err := base64.StdEncoding.DecodeString(accountKey)
	if err != nil {
		return "", err
	}
	hasher := hmac.New(sha256.New, binaryKey)
	hasher.Write([]byte(stringToSign))
	signature := hasher.Sum(nil)

	// Trial and error to determine which fields the Azure portal
	// URL encodes for a query string and which it does not.
	sasToken := "?sv=" + url.QueryEscape(signedVersion)
	sasToken += "&ss=" + url.QueryEscape(services)
	sasToken += "&srt=" + url.QueryEscape(resourceTypes)
	sasToken += "&sp=" + url.QueryEscape(permissions)
	sasToken += "&se=" + (expiry)
	sasToken += "&st=" + (start)
	sasToken += "&spr=" + (signedProtocol)

	// this is consistent with how the Azure portal builds these.
	if len(signedIp) > 0 {
		sasToken += "&sip=" + signedIp
	}

	sasToken += "&sig=" + url.QueryEscape(base64.StdEncoding.EncodeToString(signature))

	return sasToken, nil
}

// This connection string was for a real storage account which has been deleted
// so its safe to include here for reference to understand the format.
// DefaultEndpointsProtocol=https;AccountName=azurermtestsa0;AccountKey=2vJrjEyL4re2nxCEg590wJUUC7PiqqrDHjAN5RU304FNUQieiEwS2bfp83O0v28iSfWjvYhkGmjYQAdd9x+6nw==;EndpointSuffix=core.windows.net

func parseAzureStorageAccountConnectionString(connString string) (kvp map[string]string, err error) {
	validKeys := map[string]bool{"DefaultEndpointsProtocol": true, "BlobEndpoint": true,
		"AccountName": true, "AccountKey": true, "EndpointSuffix": true}
	// The k-v pairs are separated with semi-colons
	tokens := strings.Split(connString, ";")

	for _, atoken := range tokens {
		// The individual k-v are separated by an equals sign.
		kv := strings.SplitN(atoken, "=", 2)
		key := kv[0]
		val := kv[1]
		if _, present := validKeys[key]; !present {
			return nil, fmt.Errorf("[ERROR] Unknown Key: %s", key)
		}
		kvp[key] = val
	}

	if _, present := kvp[connStringAccountKeyKey]; !present {
		return nil, fmt.Errorf("[ERROR] Storage Account Key not found in connection string: %s", connString)
	}

	return kvp, nil
}
