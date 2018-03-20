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
func resourceArmStorageAccountSharedAccessSignature() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageAccountSasCreate,
		Read:   resourceArmStorageAccountSasRead,
		Delete: resourceArmStorageAccountSasDelete,

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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageAccountSasResourceTypes,
			},

			"services": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"sas": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func resourceArmStorageAccountSasCreate(d *schema.ResourceData, meta interface{}) error {

	connString := d.Get("connection_string").(string)
	httpsOnly := d.Get("https_only").(bool)
	resourceTypes := d.Get("resource_types").(string)
	services := d.Get("services").(string)
	start := d.Get("start").(string)
	expiry := d.Get("expiry").(string)
	permissions := d.Get("permissions").(string)

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

func resourceArmStorageAccountSasRead(d *schema.ResourceData, meta interface{}) error {
	// There really isn't anything to read
	// The computed SAS element stores what it needs.
	return nil
}

func resourceArmStorageAccountSasDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to delete...
	return nil
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
	permissions string, services string,
	resourceTypes string, start string, expiry string,
	signedProtocol string, signedIp string, signedVersion string) (string, error) {

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

	sasToken := "?sv=" + url.QueryEscape(signedVersion)
	sasToken += "&ss=" + url.QueryEscape(services)
	sasToken += "&srt=" + url.QueryEscape(resourceTypes)
	sasToken += "&sp=" + url.QueryEscape(permissions)
	sasToken += "&se=" + (expiry)
	sasToken += "&st=" + (start)
	sasToken += "&spr=" + url.QueryEscape(signedProtocol)

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
