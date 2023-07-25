// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

// TODO 4.0 check if this can be removed
// Hacking the SDK model, together with the Create and Get method for working around issue: https://github.com/Azure/azure-rest-api-specs/issues/21487

type DataConnectorModel struct {
	autorest.Response `json:"-"`
	Value             securityinsight.BasicDataConnector `json:"value,omitempty"`
}

func (dcm *DataConnectorModel) UnmarshalJSON(body []byte) error {
	dc, err := unmarshalBasicDataConnector(body)
	if err != nil {
		return err
	}
	dcm.Value = dc

	return nil
}

func unmarshalBasicDataConnector(body []byte) (securityinsight.BasicDataConnector, error) {
	var m map[string]interface{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	switch m["kind"] {
	case string(securityinsight.KindBasicDataConnectorKindAzureActiveDirectory):
		var adc securityinsight.AADDataConnector
		err := json.Unmarshal(body, &adc)
		return adc, err
	case string(securityinsight.KindBasicDataConnectorKindMicrosoftThreatIntelligence):
		var mdc securityinsight.MSTIDataConnector
		err := json.Unmarshal(body, &mdc)
		return mdc, err
	case string(securityinsight.KindBasicDataConnectorKindMicrosoftThreatProtection):
		var mdc securityinsight.MTPDataConnector
		err := json.Unmarshal(body, &mdc)
		return mdc, err
	case string(securityinsight.KindBasicDataConnectorKindAzureAdvancedThreatProtection):
		var adc securityinsight.AATPDataConnector
		err := json.Unmarshal(body, &adc)
		return adc, err
	case string(securityinsight.KindBasicDataConnectorKindAzureSecurityCenter):
		var adc securityinsight.ASCDataConnector
		err := json.Unmarshal(body, &adc)
		return adc, err
	case string(securityinsight.KindBasicDataConnectorKindAmazonWebServicesCloudTrail):
		var actdc securityinsight.AwsCloudTrailDataConnector
		err := json.Unmarshal(body, &actdc)
		return actdc, err
	case string(securityinsight.KindBasicDataConnectorKindAmazonWebServicesS3):
		var asdc securityinsight.AwsS3DataConnector
		err := json.Unmarshal(body, &asdc)
		return asdc, err
	case string(securityinsight.KindBasicDataConnectorKindMicrosoftCloudAppSecurity):
		var mdc securityinsight.MCASDataConnector
		err := json.Unmarshal(body, &mdc)
		return mdc, err
	case string(securityinsight.KindBasicDataConnectorKindDynamics365):
		var d3dc securityinsight.Dynamics365DataConnector
		err := json.Unmarshal(body, &d3dc)
		return d3dc, err
	case string(securityinsight.KindBasicDataConnectorKindOfficeATP):
		var oadc securityinsight.OfficeATPDataConnector
		err := json.Unmarshal(body, &oadc)
		return oadc, err
	case string(securityinsight.KindBasicDataConnectorKindOffice365Project):
		var o3pdc securityinsight.Office365ProjectDataConnector
		err := json.Unmarshal(body, &o3pdc)
		return o3pdc, err
	case string(securityinsight.KindBasicDataConnectorKindOfficePowerBI):
		var opbdc securityinsight.OfficePowerBIDataConnector
		err := json.Unmarshal(body, &opbdc)
		return opbdc, err
	case string(securityinsight.KindBasicDataConnectorKindOfficeIRM):
		var oidc securityinsight.OfficeIRMDataConnector
		err := json.Unmarshal(body, &oidc)
		return oidc, err
	case string(securityinsight.KindBasicDataConnectorKindMicrosoftDefenderAdvancedThreatProtection):
		var mdc securityinsight.MDATPDataConnector
		err := json.Unmarshal(body, &mdc)
		return mdc, err
	case string(securityinsight.KindBasicDataConnectorKindOffice365):
		var odc securityinsight.OfficeDataConnector
		err := json.Unmarshal(body, &odc)
		return odc, err
	case string(securityinsight.KindBasicDataConnectorKindThreatIntelligence):
		var tdc TIDataConnector
		err := json.Unmarshal(body, &tdc)
		return tdc, err
	case string(securityinsight.KindBasicDataConnectorKindThreatIntelligenceTaxii):
		var ttdc TiTaxiiDataConnector // using the hacked one
		err := json.Unmarshal(body, &ttdc)
		return ttdc, err
	case string(securityinsight.KindBasicDataConnectorKindIOT):
		var itdc securityinsight.IoTDataConnector
		err := json.Unmarshal(body, &itdc)
		return itdc, err
	case string(securityinsight.KindBasicDataConnectorKindGenericUI):
		var cudc securityinsight.CodelessUIDataConnector
		err := json.Unmarshal(body, &cudc)
		return cudc, err
	case string(securityinsight.KindBasicDataConnectorKindAPIPolling):
		var capdc securityinsight.CodelessAPIPollingDataConnector
		err := json.Unmarshal(body, &capdc)
		return capdc, err
	default:
		var dc securityinsight.DataConnector
		err := json.Unmarshal(body, &dc)
		return dc, err
	}
}

var _ securityinsight.BasicDataConnector = TiTaxiiDataConnector{}

type TiTaxiiDataConnector struct {
	*TiTaxiiDataConnectorProperties `json:"properties,omitempty"`
	Kind                            securityinsight.KindBasicDataConnector `json:"kind,omitempty"`
	Etag                            *string                                `json:"etag,omitempty"`
	ID                              *string                                `json:"id,omitempty"`
	Name                            *string                                `json:"name,omitempty"`
	Type                            *string                                `json:"type,omitempty"`
	SystemData                      *securityinsight.SystemData            `json:"systemData,omitempty"`
}

func (t TiTaxiiDataConnector) AsAADDataConnector() (*securityinsight.AADDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsMSTIDataConnector() (*securityinsight.MSTIDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsMTPDataConnector() (*securityinsight.MTPDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsAATPDataConnector() (*securityinsight.AATPDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsASCDataConnector() (*securityinsight.ASCDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsAwsCloudTrailDataConnector() (*securityinsight.AwsCloudTrailDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsAwsS3DataConnector() (*securityinsight.AwsS3DataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsMCASDataConnector() (*securityinsight.MCASDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsDynamics365DataConnector() (*securityinsight.Dynamics365DataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsOfficeATPDataConnector() (*securityinsight.OfficeATPDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsOffice365ProjectDataConnector() (*securityinsight.Office365ProjectDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsOfficePowerBIDataConnector() (*securityinsight.OfficePowerBIDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsOfficeIRMDataConnector() (*securityinsight.OfficeIRMDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsMDATPDataConnector() (*securityinsight.MDATPDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsOfficeDataConnector() (*securityinsight.OfficeDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsTIDataConnector() (*securityinsight.TIDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsTiTaxiiDataConnector() (*securityinsight.TiTaxiiDataConnector, bool) {
	// This method is not used at all, only for implementing the interface.
	return nil, false
}

func (t TiTaxiiDataConnector) AsIoTDataConnector() (*securityinsight.IoTDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsCodelessUIDataConnector() (*securityinsight.CodelessUIDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsCodelessAPIPollingDataConnector() (*securityinsight.CodelessAPIPollingDataConnector, bool) {
	return nil, false
}

func (t TiTaxiiDataConnector) AsDataConnector() (*securityinsight.DataConnector, bool) {
	return nil, false
}

func (ttdc TiTaxiiDataConnector) MarshalJSON() ([]byte, error) {
	ttdc.Kind = securityinsight.KindBasicDataConnectorKindThreatIntelligenceTaxii
	objectMap := make(map[string]interface{})
	if ttdc.TiTaxiiDataConnectorProperties != nil {
		objectMap["properties"] = ttdc.TiTaxiiDataConnectorProperties
	}
	if ttdc.Kind != "" {
		objectMap["kind"] = ttdc.Kind
	}
	if ttdc.Etag != nil {
		objectMap["etag"] = ttdc.Etag
	}
	return json.Marshal(objectMap)
}

func (ttdc *TiTaxiiDataConnector) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "properties":
			if v != nil {
				var tiTaxiiDataConnectorProperties TiTaxiiDataConnectorProperties
				err = json.Unmarshal(*v, &tiTaxiiDataConnectorProperties)
				if err != nil {
					return err
				}
				ttdc.TiTaxiiDataConnectorProperties = &tiTaxiiDataConnectorProperties
			}
		case "kind":
			if v != nil {
				var kind securityinsight.KindBasicDataConnector
				err = json.Unmarshal(*v, &kind)
				if err != nil {
					return err
				}
				ttdc.Kind = kind
			}
		case "etag":
			if v != nil {
				var etag string
				err = json.Unmarshal(*v, &etag)
				if err != nil {
					return err
				}
				ttdc.Etag = &etag
			}
		case "id":
			if v != nil {
				var ID string
				err = json.Unmarshal(*v, &ID)
				if err != nil {
					return err
				}
				ttdc.ID = &ID
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				ttdc.Name = &name
			}
		case "type":
			if v != nil {
				var typeVar string
				err = json.Unmarshal(*v, &typeVar)
				if err != nil {
					return err
				}
				ttdc.Type = &typeVar
			}
		case "systemData":
			if v != nil {
				var systemData securityinsight.SystemData
				err = json.Unmarshal(*v, &systemData)
				if err != nil {
					return err
				}
				ttdc.SystemData = &systemData
			}
		}
	}

	return nil
}

var _ securityinsight.BasicDataConnector = TIDataConnector{}

type TIDataConnector struct {
	*TIDataConnectorProperties `json:"properties,omitempty"`
	Kind                       securityinsight.KindBasicDataConnector `json:"kind,omitempty"`
	Etag                       *string                                `json:"etag,omitempty"`
	ID                         *string                                `json:"id,omitempty"`
	Name                       *string                                `json:"name,omitempty"`
	Type                       *string                                `json:"type,omitempty"`
	SystemData                 *securityinsight.SystemData            `json:"systemData,omitempty"`
}

func (TIDataConnector) AsAADDataConnector() (*securityinsight.AADDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsMSTIDataConnector() (*securityinsight.MSTIDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsMTPDataConnector() (*securityinsight.MTPDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsAATPDataConnector() (*securityinsight.AATPDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsASCDataConnector() (*securityinsight.ASCDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsAwsCloudTrailDataConnector() (*securityinsight.AwsCloudTrailDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsAwsS3DataConnector() (*securityinsight.AwsS3DataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsMCASDataConnector() (*securityinsight.MCASDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsDynamics365DataConnector() (*securityinsight.Dynamics365DataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsOfficeATPDataConnector() (*securityinsight.OfficeATPDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsOffice365ProjectDataConnector() (*securityinsight.Office365ProjectDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsOfficePowerBIDataConnector() (*securityinsight.OfficePowerBIDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsOfficeIRMDataConnector() (*securityinsight.OfficeIRMDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsMDATPDataConnector() (*securityinsight.MDATPDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsOfficeDataConnector() (*securityinsight.OfficeDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsTIDataConnector() (*securityinsight.TIDataConnector, bool) {
	// This method is not used at all, only for implementing the interface.
	return nil, false
}

func (TIDataConnector) AsTiTaxiiDataConnector() (*securityinsight.TiTaxiiDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsIoTDataConnector() (*securityinsight.IoTDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsCodelessUIDataConnector() (*securityinsight.CodelessUIDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsCodelessAPIPollingDataConnector() (*securityinsight.CodelessAPIPollingDataConnector, bool) {
	return nil, false
}

func (TIDataConnector) AsDataConnector() (*securityinsight.DataConnector, bool) {
	return nil, false
}

type TIDataConnectorProperties struct {
	TipLookbackPeriod *Time                                     `json:"tipLookbackPeriod,omitempty"`
	DataTypes         *securityinsight.TIDataConnectorDataTypes `json:"dataTypes,omitempty"`
	TenantID          *string                                   `json:"tenantId,omitempty"`
}

type PollingFrequency string

func (freq *PollingFrequency) UnmarshalJSON(body []byte) error {
	switch string(body) {
	case "0", string(PollingFrequencyOnceAMinute):
		*freq = PollingFrequencyOnceAMinute
	case "1", string(PollingFrequencyOnceAnHour):
		*freq = PollingFrequencyOnceAnHour
	case "2", string(PollingFrequencyOnceADay):
		*freq = PollingFrequencyOnceADay
	default:
		return fmt.Errorf("unknown enum for pollingFrequency %s", string(body))
	}
	return nil
}

const (
	PollingFrequencyOnceAMinute PollingFrequency = "OnceAMinute" // API returns 0
	PollingFrequencyOnceAnHour  PollingFrequency = "OnceAnHour"  // API returns 1
	PollingFrequencyOnceADay    PollingFrequency = "OnceADay"    // API returns 2
)

type Time date.Time

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	// Firstly, try to parse the date time via RFC3339, which is the expected format defined by Swagger.
	// However, since the service issue (#21487), it currently doesn't return in this format.
	// In order not to break the code once the service fix it, we keep this try at first.
	if time, err := time.Parse(time.RFC3339, string(data)); err == nil {
		t.Time = time
		return nil
	}

	// This is the format that the service returns at this moment, which is not the expected format (RFC3339).
	t.Time, err = time.Parse(`"01/02/2006 15:04:05"`, string(data))
	return err
}

type TiTaxiiDataConnectorProperties struct {
	WorkspaceID         *string                                        `json:"workspaceId,omitempty"`
	FriendlyName        *string                                        `json:"friendlyName,omitempty"`
	TaxiiServer         *string                                        `json:"taxiiServer,omitempty"`
	CollectionID        *string                                        `json:"collectionId,omitempty"`
	UserName            *string                                        `json:"userName,omitempty"`
	Password            *string                                        `json:"password,omitempty"`
	TaxiiLookbackPeriod *Time                                          `json:"taxiiLookbackPeriod,omitempty"`
	PollingFrequency    PollingFrequency                               `json:"pollingFrequency,omitempty"`
	DataTypes           *securityinsight.TiTaxiiDataConnectorDataTypes `json:"dataTypes,omitempty"`
	TenantID            *string                                        `json:"tenantId,omitempty"`
}
