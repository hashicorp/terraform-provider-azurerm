package frontdoors

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/formatting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

type ForwardingConfiguration struct {
	BackendPool *SubResource `json:"backendPool,omitempty"`
	CacheConfiguration *CacheConfiguration `json:"cacheConfiguration,omitempty"`
	CustomForwardingPath *string `json:"customForwardingPath,omitempty"`
	ForwardingProtocol *FrontDoorForwardingProtocol `json:"forwardingProtocol,omitempty"`
}



func (c *ForwardingConfiguration) UnmarshalJSON(input []byte) error {
	type intermediateType struct {
	BackendPool *SubResource `json:"backendPool,omitempty"`
	CacheConfiguration *CacheConfiguration `json:"cacheConfiguration,omitempty"`
	CustomForwardingPath *string `json:"customForwardingPath,omitempty"`
	ForwardingProtocol *FrontDoorForwardingProtocol `json:"forwardingProtocol,omitempty"`
	OdataType json.RawMessage `json:"@odata.type"`
	}
	var intermediate intermediateType
	if err := json.Unmarshal(input, &intermediate); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	c.BackendPool = intermediate.BackendPool
	c.CacheConfiguration = intermediate.CacheConfiguration
	c.CustomForwardingPath = intermediate.CustomForwardingPath
	c.ForwardingProtocol = intermediate.ForwardingProtocol

	@odata.type, err := unmarshalstring(intermediate.OdataType)
	if err != nil {
		return fmt.Errorf("unmarshaling @odata.type: %+v", err)
	}
	c.OdataType = @odata.type


	return nil
}


var _ json.Marshaler = ForwardingConfiguration{}

func (o ForwardingConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
"backendPool": o.BackendPool,
"cacheConfiguration": o.CacheConfiguration,
"customForwardingPath": o.CustomForwardingPath,
"forwardingProtocol": o.ForwardingProtocol,
"@odata.type": "#Microsoft.Azure.FrontDoor.Models.FrontdoorForwardingConfiguration",
	})
}

