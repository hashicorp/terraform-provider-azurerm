package videoanalyzer

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/formatting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

type JwtAuthentication struct {
	Audiences *[]string `json:"audiences,omitempty"`
	Claims *[]TokenClaim `json:"claims,omitempty"`
	Issuers *[]string `json:"issuers,omitempty"`
	Keys *[]TokenKey `json:"keys,omitempty"`
}



func (c *JwtAuthentication) UnmarshalJSON(input []byte) error {
	type intermediateType struct {
	Audiences *[]string `json:"audiences,omitempty"`
	Claims *[]TokenClaim `json:"claims,omitempty"`
	Issuers *[]string `json:"issuers,omitempty"`
	Keys *[]TokenKey `json:"keys,omitempty"`
	Type json.RawMessage `json:"@type"`
	}
	var intermediate intermediateType
	if err := json.Unmarshal(input, &intermediate); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	c.Audiences = intermediate.Audiences
	c.Claims = intermediate.Claims
	c.Issuers = intermediate.Issuers
	c.Keys = intermediate.Keys

	@type, err := unmarshalstring(intermediate.Type)
	if err != nil {
		return fmt.Errorf("unmarshaling @type: %+v", err)
	}
	c.Type = @type


	return nil
}


var _ json.Marshaler = JwtAuthentication{}

func (o JwtAuthentication) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
"audiences": o.Audiences,
"claims": o.Claims,
"issuers": o.Issuers,
"keys": o.Keys,
"@type": "#Microsoft.VideoAnalyzer.JwtAuthentication",
	})
}

