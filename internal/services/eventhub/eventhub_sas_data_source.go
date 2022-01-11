package eventhub

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/hashicorp/go-azure-helpers/eventhub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const (
	connStringSharedAccessKeyKey     = "SharedAccessKey"
	connStringSharedAccessKeyNameKey = "SharedAccessKeyName"
	connStringEndpointKey            = "Endpoint"
	connStringEntityPathKey          = "EntityPath"
)

func dataSourceEventHubSharedAccessSignature() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceEventHubSasRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"connection_string": {
				Type:      pluginsdk.TypeString,
				Required:  true,
				Sensitive: true,
			},

			// Always in UTC and must be ISO-8601 format
			"expiry": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ISO8601DateTime,
			},

			"sas": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceEventHubSasRead(d *pluginsdk.ResourceData, _ interface{}) error {
	connString := d.Get("connection_string").(string)
	expiry := d.Get("expiry").(string)

	// Parse the connection string
	kvp, err := eventhub.ParseEventHubSASConnectionString(connString)
	if err != nil {
		return err
	}

	sharedAccessKeyName := kvp[connStringSharedAccessKeyNameKey]
	sharedAccessKey := kvp[connStringSharedAccessKeyKey]
	endpoint := kvp[connStringEndpointKey]
	entityPath := kvp[connStringEntityPathKey]
	endpointUrl, err := eventhub.ComputeEventHubSASConnectionUrl(endpoint, entityPath)
	if err != nil {
		return err
	}

	sasToken, err := eventhub.ComputeEventHubSASToken(sharedAccessKeyName, sharedAccessKey, expiry, *endpointUrl)
	if err != nil {
		return err
	}

	sasConnectionString := eventhub.ComputeEventHubSASConnectionString(sasToken)

	d.Set("sas", sasConnectionString)
	tokenHash := sha256.Sum256([]byte(sasToken))
	d.SetId(hex.EncodeToString(tokenHash[:]))

	return nil
}
