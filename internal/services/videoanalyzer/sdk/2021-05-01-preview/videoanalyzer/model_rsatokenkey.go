package videoanalyzer

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/formatting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

type RsaTokenKey struct {
	Alg AccessPolicyRsaAlgo `json:"alg"`
	E string `json:"e"`
	Kid string `json:"kid"`
	N string `json:"n"`
}



func (c *RsaTokenKey) UnmarshalJSON(input []byte) error {
	type intermediateType struct {
	Alg AccessPolicyRsaAlgo `json:"alg"`
	E string `json:"e"`
	Kid string `json:"kid"`
	N string `json:"n"`
	Type json.RawMessage `json:"@type"`
	}
	var intermediate intermediateType
	if err := json.Unmarshal(input, &intermediate); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	c.Alg = intermediate.Alg
	c.E = intermediate.E
	c.Kid = intermediate.Kid
	c.N = intermediate.N

	@type, err := unmarshalstring(intermediate.Type)
	if err != nil {
		return fmt.Errorf("unmarshaling @type: %+v", err)
	}
	c.Type = @type


	return nil
}


var _ json.Marshaler = RsaTokenKey{}

func (o RsaTokenKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
"alg": o.Alg,
"e": o.E,
"kid": o.Kid,
"n": o.N,
"@type": "#Microsoft.VideoAnalyzer.RsaTokenKey",
	})
}

