package azuresdkhacks

//
//import (
//	"encoding/json"
//	"fmt"
//
//	"github.com/Azure/go-autorest/autorest/date"
//	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
//	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertrules"
//	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
//)
//
//// TODO 4.0 check if this can be removed
//// tracked on https://github.com/Azure/azure-rest-api-specs/issues/16615
//
//var _ alertrules.AlertRule = ThreatIntelligenceAlertRule{}
//
//type ThreatIntelligenceAlertRule struct {
//	Properties *ThreatIntelligenceAlertRuleProperties `json:"properties,omitempty"`
//
//	// Fields inherited from AlertRule
//	Etag       *string                `json:"etag,omitempty"`
//	Id         *string                `json:"id,omitempty"`
//	Name       *string                `json:"name,omitempty"`
//	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
//	Type       *string                `json:"type,omitempty"`
//}
//
//var _ json.Marshaler = ThreatIntelligenceAlertRule{}
//
//func (s ThreatIntelligenceAlertRule) MarshalJSON() ([]byte, error) {
//	type wrapper ThreatIntelligenceAlertRule
//	wrapped := wrapper(s)
//	encoded, err := json.Marshal(wrapped)
//	if err != nil {
//		return nil, fmt.Errorf("marshaling ThreatIntelligenceAlertRule: %+v", err)
//	}
//
//	var decoded map[string]interface{}
//	if err := json.Unmarshal(encoded, &decoded); err != nil {
//		return nil, fmt.Errorf("unmarshaling ThreatIntelligenceAlertRule: %+v", err)
//	}
//	decoded["kind"] = "ThreatIntelligence"
//
//	encoded, err = json.Marshal(decoded)
//	if err != nil {
//		return nil, fmt.Errorf("re-marshaling ThreatIntelligenceAlertRule: %+v", err)
//	}
//
//	return encoded, nil
//}
//
//func (tiar ThreatIntelligenceAlertRule) MarshalJSON() ([]byte, error) {
//	tiar.Kind = securityinsight.KindBasicAlertRuleKindThreatIntelligence
//	objectMap := make(map[string]interface{})
//	if tiar.ThreatIntelligenceAlertRuleProperties != nil {
//		objectMap["properties"] = tiar.ThreatIntelligenceAlertRuleProperties
//	}
//	if tiar.Kind != "" {
//		objectMap["kind"] = tiar.Kind
//	}
//	if tiar.Etag != nil {
//		objectMap["etag"] = tiar.Etag
//	}
//	return json.Marshal(objectMap)
//}
//
//// MarshalJSON is the custom marshaler for ThreatIntelligenceAlertRuleProperties.
//func (tiarp ThreatIntelligenceAlertRuleProperties) MarshalJSON() ([]byte, error) {
//	objectMap := make(map[string]interface{})
//	if tiarp.AlertRuleTemplateName != nil {
//		objectMap["alertRuleTemplateName"] = tiarp.AlertRuleTemplateName
//	}
//	if tiarp.Enabled != nil {
//		objectMap["enabled"] = tiarp.Enabled
//	}
//	objectMap["severity"] = tiarp.Severity
//	if tiarp.DisplayName != nil {
//		objectMap["displayName"] = tiarp.DisplayName
//	}
//	if tiarp.Description != nil {
//		objectMap["description"] = tiarp.Description
//	}
//	if tiarp.Tactics != nil {
//		objectMap["tactics"] = tiarp.Tactics
//	}
//	return json.Marshal(objectMap)
//}
