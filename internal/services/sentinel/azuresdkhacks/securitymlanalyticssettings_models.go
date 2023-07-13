// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/gofrs/uuid"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

// TODO 4.0 check if this can be removed
// Hacking the SDK model, adding definition for customizableObservations, tracked on https://github.com/Azure/azure-rest-api-specs/issues/22503

type SecurityMLAnalyticsSettingsListIterator struct {
	i    int
	page SecurityMLAnalyticsSettingsListPage
}

func (iter *SecurityMLAnalyticsSettingsListIterator) NextWithContext(ctx context.Context) (err error) {
	iter.i++
	if iter.i < len(iter.page.Values()) {
		return nil
	}
	err = iter.page.NextWithContext(ctx)
	if err != nil {
		iter.i--
		return err
	}
	iter.i = 0
	return nil
}

func (iter *SecurityMLAnalyticsSettingsListIterator) Next() error {
	return iter.NextWithContext(context.Background())
}

func (iter SecurityMLAnalyticsSettingsListIterator) NotDone() bool {
	return iter.page.NotDone() && iter.i < len(iter.page.Values())
}

func (iter SecurityMLAnalyticsSettingsListIterator) Response() SecurityMLAnalyticsSettingsList {
	return iter.page.Response()
}

func (iter SecurityMLAnalyticsSettingsListIterator) Value() BasicSecurityMLAnalyticsSetting {
	if !iter.page.NotDone() {
		return SecurityMLAnalyticsSetting{}
	}
	return iter.page.Values()[iter.i]
}

type SecurityMLAnalyticsSettingsListPage struct {
	fn    func(context.Context, SecurityMLAnalyticsSettingsList) (SecurityMLAnalyticsSettingsList, error)
	smasl SecurityMLAnalyticsSettingsList
}

func (page *SecurityMLAnalyticsSettingsListPage) NextWithContext(ctx context.Context) (err error) {
	for {
		next, err := page.fn(ctx, page.smasl)
		if err != nil {
			return err
		}
		page.smasl = next
		if !next.hasNextLink() || !next.IsEmpty() {
			break
		}
	}
	return nil
}

func (page *SecurityMLAnalyticsSettingsListPage) Next() error {
	return page.NextWithContext(context.Background())
}

func (page SecurityMLAnalyticsSettingsListPage) NotDone() bool {
	return !page.smasl.IsEmpty()
}

func (page SecurityMLAnalyticsSettingsListPage) Response() SecurityMLAnalyticsSettingsList {
	return page.smasl
}

func (page SecurityMLAnalyticsSettingsListPage) Values() []BasicSecurityMLAnalyticsSetting {
	if page.smasl.IsEmpty() {
		return nil
	}
	return *page.smasl.Value
}

type SecurityMLAnalyticsSettingsList struct {
	autorest.Response `json:"-"`
	// NextLink - READ-ONLY; URL to fetch the next set of SecurityMLAnalyticsSettings.
	NextLink *string `json:"nextLink,omitempty"`
	// Value - Array of SecurityMLAnalyticsSettings
	Value *[]BasicSecurityMLAnalyticsSetting `json:"value,omitempty"`
}

func (smasl SecurityMLAnalyticsSettingsList) securityMLAnalyticsSettingsListPreparer(ctx context.Context) (*http.Request, error) {
	if !smasl.hasNextLink() {
		return nil, nil
	}
	return autorest.Prepare((&http.Request{}).WithContext(ctx),
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(smasl.NextLink)))
}

// IsEmpty returns true if the ListResult contains no values.
func (smasl SecurityMLAnalyticsSettingsList) IsEmpty() bool {
	return smasl.Value == nil || len(*smasl.Value) == 0
}

// hasNextLink returns true if the NextLink is not empty.
func (smasl SecurityMLAnalyticsSettingsList) hasNextLink() bool {
	return smasl.NextLink != nil && len(*smasl.NextLink) != 0
}

// MarshalJSON is the custom marshaler for SecurityMLAnalyticsSettingsList.
func (smasl SecurityMLAnalyticsSettingsList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if smasl.Value != nil {
		objectMap["value"] = smasl.Value
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON is the custom unmarshaler for SecurityMLAnalyticsSettingsList struct.
func (smasl *SecurityMLAnalyticsSettingsList) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "nextLink":
			if v != nil {
				var nextLink string
				err = json.Unmarshal(*v, &nextLink)
				if err != nil {
					return err
				}
				smasl.NextLink = &nextLink
			}
		case "value":
			if v != nil {
				value, err := unmarshalBasicSecurityMLAnalyticsSettingArray(*v)
				if err != nil {
					return err
				}
				smasl.Value = &value
			}
		}
	}

	return nil
}

func unmarshalBasicSecurityMLAnalyticsSettingArray(body []byte) ([]BasicSecurityMLAnalyticsSetting, error) {
	var rawMessages []*json.RawMessage
	err := json.Unmarshal(body, &rawMessages)
	if err != nil {
		return nil, err
	}

	smasArray := make([]BasicSecurityMLAnalyticsSetting, len(rawMessages))

	for index, rawMessage := range rawMessages {
		smas, err := unmarshalBasicSecurityMLAnalyticsSetting(*rawMessage)
		if err != nil {
			return nil, err
		}
		smasArray[index] = smas
	}
	return smasArray, nil
}

func unmarshalBasicSecurityMLAnalyticsSetting(body []byte) (BasicSecurityMLAnalyticsSetting, error) {
	var m map[string]interface{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	switch m["kind"] {
	case string(securityinsight.KindBasicSecurityMLAnalyticsSettingKindAnomaly):
		var asmas AnomalySecurityMLAnalyticsSettings
		err := json.Unmarshal(body, &asmas)
		return asmas, err
	default:
		var smas SecurityMLAnalyticsSetting
		err := json.Unmarshal(body, &smas)
		return smas, err
	}
}

// BasicSecurityMLAnalyticsSetting security ML Analytics Setting
type BasicSecurityMLAnalyticsSetting interface {
	AsAnomalySecurityMLAnalyticsSettings() (*AnomalySecurityMLAnalyticsSettings, bool)
	AsSecurityMLAnalyticsSetting() (*SecurityMLAnalyticsSetting, bool)
}

// SecurityMLAnalyticsSetting security ML Analytics Setting
type SecurityMLAnalyticsSetting struct {
	autorest.Response `json:"-"`
	// Kind - Possible values include: 'KindBasicSecurityMLAnalyticsSettingKindSecurityMLAnalyticsSetting', 'KindBasicSecurityMLAnalyticsSettingKindAnomaly'
	Kind securityinsight.KindBasicSecurityMLAnalyticsSetting `json:"kind,omitempty"`
	// Etag - Etag of the azure resource
	Etag *string `json:"etag,omitempty"`
	// ID - READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty"`
	// SystemData - READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *securityinsight.SystemData `json:"systemData,omitempty"`
}

func (smas SecurityMLAnalyticsSetting) AsAnomalySecurityMLAnalyticsSettings() (*AnomalySecurityMLAnalyticsSettings, bool) {
	return nil, false
}
func (smas SecurityMLAnalyticsSetting) AsSecurityMLAnalyticsSetting() (*SecurityMLAnalyticsSetting, bool) {
	return &smas, true
}
func (smas SecurityMLAnalyticsSetting) AsBasicSecurityMLAnalyticsSetting() (BasicSecurityMLAnalyticsSetting, bool) {
	return &smas, true
}

type AnomalySecurityMLAnalyticsSettings struct {
	*AnomalySecurityMLAnalyticsSettingsProperties `json:"properties,omitempty"`
	Etag                                          *string                                             `json:"etag,omitempty"`
	ID                                            *string                                             `json:"id,omitempty"`
	Name                                          *string                                             `json:"name,omitempty"`
	Type                                          *string                                             `json:"type,omitempty"`
	SystemData                                    *securityinsight.SystemData                         `json:"systemData,omitempty"`
	Kind                                          securityinsight.KindBasicSecurityMLAnalyticsSetting `json:"kind,omitempty"`
}

func (asmas AnomalySecurityMLAnalyticsSettings) AsAnomalySecurityMLAnalyticsSettings() (*AnomalySecurityMLAnalyticsSettings, bool) {
	return &asmas, true
}
func (asmas AnomalySecurityMLAnalyticsSettings) AsSecurityMLAnalyticsSetting() (*SecurityMLAnalyticsSetting, bool) {
	return nil, false
}
func (asmas AnomalySecurityMLAnalyticsSettings) AsBasicSecurityMLAnalyticsSetting() (BasicSecurityMLAnalyticsSetting, bool) {
	return &asmas, true
}

type AnomalySecurityMLAnalyticsSettingsProperties struct {
	Description              *string                                                  `json:"description,omitempty"`
	DisplayName              *string                                                  `json:"displayName,omitempty"`
	Enabled                  *bool                                                    `json:"enabled,omitempty"`
	LastModifiedUtc          *date.Time                                               `json:"lastModifiedUtc,omitempty"`
	RequiredDataConnectors   *[]securityinsight.SecurityMLAnalyticsSettingsDataSource `json:"requiredDataConnectors,omitempty"`
	Tactics                  *[]securityinsight.AttackTactic                          `json:"tactics,omitempty"`
	Techniques               *[]string                                                `json:"techniques,omitempty"`
	AnomalyVersion           *string                                                  `json:"anomalyVersion,omitempty"`
	CustomizableObservations *AnomalySecurityMLAnalyticsCustomizableObservations      `json:"customizableObservations,omitempty"`
	Frequency                *string                                                  `json:"frequency,omitempty"`
	SettingsStatus           securityinsight.SettingsStatus                           `json:"settingsStatus,omitempty"`
	IsDefaultSettings        *bool                                                    `json:"isDefaultSettings,omitempty"`
	AnomalySettingsVersion   *int32                                                   `json:"anomalySettingsVersion,omitempty"`
	SettingsDefinitionID     *uuid.UUID                                               `json:"settingsDefinitionId,omitempty"`
}

type AnomalySecurityMLAnalyticsCustomizableObservations struct {
	MultiSelectObservations       *[]AnomalySecurityMLAnalyticsMultiSelectObservations       `json:"multiSelectObservations,omitempty"`
	SingleSelectObservations      *[]AnomalySecurityMLAnalyticsSingleSelectObservations      `json:"singleSelectObservations,omitempty"`
	PrioritizeExcludeObservations *[]AnomalySecurityMLAnalyticsPrioritizeExcludeObservations `json:"prioritizeExcludeObservations,omitempty"`
	ThresholdObservations         *[]AnomalySecurityMLAnalyticsThresholdObservations         `json:"thresholdObservations,omitempty"`
}

// unused properties are defined to interface{}.
type AnomalySecurityMLAnalyticsMultiSelectObservations struct {
	SupportValues      *[]string    `json:"supportedValues,omitempty"`
	Values             *[]string    `json:"values,omitempty"`
	Name               *string      `json:"name,omitempty"`
	Description        *string      `json:"description,omitempty"`
	SupportedValuesKql *interface{} `json:"supportedValuesKql,omitempty"`
	ValuesKql          *interface{} `json:"valuesKql,omitempty"`
	SequenceNumber     *interface{} `json:"sequenceNumber,omitempty"`
	Rerun              *interface{} `json:"rerun,omitempty"`
}

type AnomalySecurityMLAnalyticsSingleSelectObservations struct {
	SupportValues      *[]string    `json:"supportedValues,omitempty"`
	Value              *string      `json:"value,omitempty"`
	Name               *string      `json:"name,omitempty"`
	Description        *string      `json:"description,omitempty"`
	SupportedValuesKql *interface{} `json:"supportedValuesKql,omitempty"`
	SequenceNumber     *interface{} `json:"sequenceNumber,omitempty"`
	Rerun              *interface{} `json:"rerun,omitempty"`
}

type AnomalySecurityMLAnalyticsPrioritizeExcludeObservations struct {
	Name           *string      `json:"name,omitempty"`
	Description    *string      `json:"description,omitempty"`
	Prioritize     *string      `json:"prioritize,omitempty"`
	Exclude        *string      `json:"exclude,omitempty"`
	DataType       *interface{} `json:"dataType,omitempty"`
	SequenceNumber *interface{} `json:"sequenceNumber,omitempty"`
	Rerun          *interface{} `json:"rerun,omitempty"`
}

type AnomalySecurityMLAnalyticsThresholdObservations struct {
	Name           *string      `json:"name,omitempty"`
	Description    *string      `json:"description,omitempty"`
	Max            *string      `json:"maximum,omitempty"`
	Min            *string      `json:"minimum,omitempty"`
	Value          *string      `json:"value,omitempty"`
	SequenceNumber *interface{} `json:"sequenceNumber,omitempty"`
	Rerun          *interface{} `json:"rerun,omitempty"`
}
