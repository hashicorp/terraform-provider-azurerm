package azuresdkhacks

import (
	"encoding/json"

	"github.com/Azure/go-autorest/autorest/date"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

// TODO 4.0 check if this can be removed
// tracked on https://github.com/Azure/azure-rest-api-specs/issues/16615

var _ securityinsight.BasicAlertRule = ThreatIntelligenceAlertRule{}

type ThreatIntelligenceAlertRule struct {
	*ThreatIntelligenceAlertRuleProperties `json:"properties,omitempty"`
	Kind                                   securityinsight.KindBasicAlertRule `json:"kind,omitempty"`
	Etag                                   *string                            `json:"etag,omitempty"`
	ID                                     *string                            `json:"id,omitempty"`
	Name                                   *string                            `json:"name,omitempty"`
	Type                                   *string                            `json:"type,omitempty"`
	SystemData                             *securityinsight.SystemData        `json:"systemData,omitempty"`
}

func (tiar ThreatIntelligenceAlertRule) MarshalJSON() ([]byte, error) {
	tiar.Kind = securityinsight.KindBasicAlertRuleKindThreatIntelligence
	objectMap := make(map[string]interface{})
	if tiar.ThreatIntelligenceAlertRuleProperties != nil {
		objectMap["properties"] = tiar.ThreatIntelligenceAlertRuleProperties
	}
	if tiar.Kind != "" {
		objectMap["kind"] = tiar.Kind
	}
	if tiar.Etag != nil {
		objectMap["etag"] = tiar.Etag
	}
	return json.Marshal(objectMap)
}

func (t ThreatIntelligenceAlertRule) AsMLBehaviorAnalyticsAlertRule() (*securityinsight.MLBehaviorAnalyticsAlertRule, bool) {
	return nil, false
}

func (t ThreatIntelligenceAlertRule) AsFusionAlertRule() (*securityinsight.FusionAlertRule, bool) {
	return nil, false
}

func (t ThreatIntelligenceAlertRule) AsThreatIntelligenceAlertRule() (*securityinsight.ThreatIntelligenceAlertRule, bool) {
	return nil, false
}

func (t ThreatIntelligenceAlertRule) AsMicrosoftSecurityIncidentCreationAlertRule() (*securityinsight.MicrosoftSecurityIncidentCreationAlertRule, bool) {
	return nil, false
}

func (t ThreatIntelligenceAlertRule) AsScheduledAlertRule() (*securityinsight.ScheduledAlertRule, bool) {
	return nil, false
}

func (t ThreatIntelligenceAlertRule) AsNrtAlertRule() (*securityinsight.NrtAlertRule, bool) {
	return nil, false
}

func (t ThreatIntelligenceAlertRule) AsAlertRule() (*securityinsight.AlertRule, bool) {
	return nil, false
}

type ThreatIntelligenceAlertRuleProperties struct {
	AlertRuleTemplateName *string                         `json:"alertRuleTemplateName,omitempty"`
	Description           *string                         `json:"description,omitempty"`
	DisplayName           *string                         `json:"displayName,omitempty"`
	Enabled               *bool                           `json:"enabled,omitempty"`
	LastModifiedUtc       *date.Time                      `json:"lastModifiedUtc,omitempty"`
	Severity              securityinsight.AlertSeverity   `json:"severity,omitempty"`
	Tactics               *[]securityinsight.AttackTactic `json:"tactics,omitempty"`
	Techniques            *[]string                       `json:"techniques,omitempty"`
}

// MarshalJSON is the custom marshaler for ThreatIntelligenceAlertRuleProperties.
func (tiarp ThreatIntelligenceAlertRuleProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if tiarp.AlertRuleTemplateName != nil {
		objectMap["alertRuleTemplateName"] = tiarp.AlertRuleTemplateName
	}
	if tiarp.Enabled != nil {
		objectMap["enabled"] = tiarp.Enabled
	}
	objectMap["severity"] = tiarp.Severity
	if tiarp.DisplayName != nil {
		objectMap["displayName"] = tiarp.DisplayName
	}
	if tiarp.Description != nil {
		objectMap["description"] = tiarp.Description
	}
	if tiarp.Tactics != nil {
		objectMap["tactics"] = tiarp.Tactics
	}
	return json.Marshal(objectMap)
}
