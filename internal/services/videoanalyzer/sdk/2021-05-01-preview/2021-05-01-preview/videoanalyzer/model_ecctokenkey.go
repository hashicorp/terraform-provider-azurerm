package videoanalyzer

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/formatting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

type EccTokenKey struct {
	Alg AccessPolicyEccAlgo `json:"alg"`
	Kid string `json:"kid"`
	X string `json:"x"`
	Y string `json:"y"`
}



func (c *EccTokenKey) UnmarshalJSON(input []byte) error {
	type intermediateType struct {
	Alg AccessPolicyEccAlgo `json:"alg"`
	Kid string `json:"kid"`
	X string `json:"x"`
	Y string `json:"y"`
	Type json.RawMessage `json:"@type"`
	}
	var intermediate intermediateType
	if err := json.Unmarshal(input, &intermediate); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	c.Alg = intermediate.Alg
	c.Kid = intermediate.Kid
	c.X = intermediate.X
	c.Y = intermediate.Y

	@type, err := unmarshalstring(intermediate.Type)
	if err != nil {
		return fmt.Errorf("unmarshaling @type: %+v", err)
	}
	c.Type = @type


	return nil
}


var _ json.Marshaler = EccTokenKey{}

func (o EccTokenKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
"alg": o.Alg,
"kid": o.Kid,
"@type": "#Microsoft.VideoAnalyzer.EccTokenKey",
"x": o.X,
"y": o.Y,
	})
}

