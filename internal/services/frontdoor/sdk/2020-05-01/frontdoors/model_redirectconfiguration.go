package frontdoors

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/formatting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/identity"
)

type RedirectConfiguration struct {
	CustomFragment *string `json:"customFragment,omitempty"`
	CustomHost *string `json:"customHost,omitempty"`
	CustomPath *string `json:"customPath,omitempty"`
	CustomQueryString *string `json:"customQueryString,omitempty"`
	RedirectProtocol *FrontDoorRedirectProtocol `json:"redirectProtocol,omitempty"`
	RedirectType *FrontDoorRedirectType `json:"redirectType,omitempty"`
}



func (c *RedirectConfiguration) UnmarshalJSON(input []byte) error {
	type intermediateType struct {
	CustomFragment *string `json:"customFragment,omitempty"`
	CustomHost *string `json:"customHost,omitempty"`
	CustomPath *string `json:"customPath,omitempty"`
	CustomQueryString *string `json:"customQueryString,omitempty"`
	RedirectProtocol *FrontDoorRedirectProtocol `json:"redirectProtocol,omitempty"`
	RedirectType *FrontDoorRedirectType `json:"redirectType,omitempty"`
	OdataType json.RawMessage `json:"@odata.type"`
	}
	var intermediate intermediateType
	if err := json.Unmarshal(input, &intermediate); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	c.CustomFragment = intermediate.CustomFragment
	c.CustomHost = intermediate.CustomHost
	c.CustomPath = intermediate.CustomPath
	c.CustomQueryString = intermediate.CustomQueryString
	c.RedirectProtocol = intermediate.RedirectProtocol
	c.RedirectType = intermediate.RedirectType

	@odata.type, err := unmarshalstring(intermediate.OdataType)
	if err != nil {
		return fmt.Errorf("unmarshaling @odata.type: %+v", err)
	}
	c.OdataType = @odata.type


	return nil
}


var _ json.Marshaler = RedirectConfiguration{}

func (o RedirectConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
"customFragment": o.CustomFragment,
"customHost": o.CustomHost,
"customPath": o.CustomPath,
"customQueryString": o.CustomQueryString,
"@odata.type": "#Microsoft.Azure.FrontDoor.Models.FrontdoorRedirectConfiguration",
"redirectProtocol": o.RedirectProtocol,
"redirectType": o.RedirectType,
	})
}

