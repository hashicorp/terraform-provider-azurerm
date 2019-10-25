package network

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func ValidatePrivateLinkEndpointSettings(d *schema.ResourceData) error {
	privateServiceConnections := d.Get("private_service_connection").([]interface{})

	for _, psc := range privateServiceConnections {
		privateServiceConnection := psc.(map[string]interface{})
		name := privateServiceConnection["name"].(string)

		// If this is not a manule connection and the message is set return an error since this does not make sense.
		if !privateServiceConnection["is_manual_connection"].(bool) && privateServiceConnection["request_message"].(string) != "" {
			return fmt.Errorf(`"private_service_connection":%q is invalid, the "request_message" attribute cannot be set if the "is_manual_connection" attribute is "false"`, name)
		}
	}

	return nil
}