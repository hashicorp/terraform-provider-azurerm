package alerts

import "strings"

type AlertCategory string

const (
	AlertCategoryBilling AlertCategory = "Billing"
	AlertCategoryCost    AlertCategory = "Cost"
	AlertCategorySystem  AlertCategory = "System"
	AlertCategoryUsage   AlertCategory = "Usage"
)

func PossibleValuesForAlertCategory() []string {
	return []string{
		string(AlertCategoryBilling),
		string(AlertCategoryCost),
		string(AlertCategorySystem),
		string(AlertCategoryUsage),
	}
}

func parseAlertCategory(input string) (*AlertCategory, error) {
	vals := map[string]AlertCategory{
		"billing": AlertCategoryBilling,
		"cost":    AlertCategoryCost,
		"system":  AlertCategorySystem,
		"usage":   AlertCategoryUsage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertCategory(input)
	return &out, nil
}

type AlertCriteria string

const (
	AlertCriteriaCostThresholdExceeded          AlertCriteria = "CostThresholdExceeded"
	AlertCriteriaCreditThresholdApproaching     AlertCriteria = "CreditThresholdApproaching"
	AlertCriteriaCreditThresholdReached         AlertCriteria = "CreditThresholdReached"
	AlertCriteriaCrossCloudCollectionError      AlertCriteria = "CrossCloudCollectionError"
	AlertCriteriaCrossCloudNewDataAvailable     AlertCriteria = "CrossCloudNewDataAvailable"
	AlertCriteriaForecastCostThresholdExceeded  AlertCriteria = "ForecastCostThresholdExceeded"
	AlertCriteriaForecastUsageThresholdExceeded AlertCriteria = "ForecastUsageThresholdExceeded"
	AlertCriteriaGeneralThresholdError          AlertCriteria = "GeneralThresholdError"
	AlertCriteriaInvoiceDueDateApproaching      AlertCriteria = "InvoiceDueDateApproaching"
	AlertCriteriaInvoiceDueDateReached          AlertCriteria = "InvoiceDueDateReached"
	AlertCriteriaMultiCurrency                  AlertCriteria = "MultiCurrency"
	AlertCriteriaQuotaThresholdApproaching      AlertCriteria = "QuotaThresholdApproaching"
	AlertCriteriaQuotaThresholdReached          AlertCriteria = "QuotaThresholdReached"
	AlertCriteriaUsageThresholdExceeded         AlertCriteria = "UsageThresholdExceeded"
)

func PossibleValuesForAlertCriteria() []string {
	return []string{
		string(AlertCriteriaCostThresholdExceeded),
		string(AlertCriteriaCreditThresholdApproaching),
		string(AlertCriteriaCreditThresholdReached),
		string(AlertCriteriaCrossCloudCollectionError),
		string(AlertCriteriaCrossCloudNewDataAvailable),
		string(AlertCriteriaForecastCostThresholdExceeded),
		string(AlertCriteriaForecastUsageThresholdExceeded),
		string(AlertCriteriaGeneralThresholdError),
		string(AlertCriteriaInvoiceDueDateApproaching),
		string(AlertCriteriaInvoiceDueDateReached),
		string(AlertCriteriaMultiCurrency),
		string(AlertCriteriaQuotaThresholdApproaching),
		string(AlertCriteriaQuotaThresholdReached),
		string(AlertCriteriaUsageThresholdExceeded),
	}
}

func parseAlertCriteria(input string) (*AlertCriteria, error) {
	vals := map[string]AlertCriteria{
		"costthresholdexceeded":          AlertCriteriaCostThresholdExceeded,
		"creditthresholdapproaching":     AlertCriteriaCreditThresholdApproaching,
		"creditthresholdreached":         AlertCriteriaCreditThresholdReached,
		"crosscloudcollectionerror":      AlertCriteriaCrossCloudCollectionError,
		"crosscloudnewdataavailable":     AlertCriteriaCrossCloudNewDataAvailable,
		"forecastcostthresholdexceeded":  AlertCriteriaForecastCostThresholdExceeded,
		"forecastusagethresholdexceeded": AlertCriteriaForecastUsageThresholdExceeded,
		"generalthresholderror":          AlertCriteriaGeneralThresholdError,
		"invoiceduedateapproaching":      AlertCriteriaInvoiceDueDateApproaching,
		"invoiceduedatereached":          AlertCriteriaInvoiceDueDateReached,
		"multicurrency":                  AlertCriteriaMultiCurrency,
		"quotathresholdapproaching":      AlertCriteriaQuotaThresholdApproaching,
		"quotathresholdreached":          AlertCriteriaQuotaThresholdReached,
		"usagethresholdexceeded":         AlertCriteriaUsageThresholdExceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertCriteria(input)
	return &out, nil
}

type AlertOperator string

const (
	AlertOperatorEqualTo              AlertOperator = "EqualTo"
	AlertOperatorGreaterThan          AlertOperator = "GreaterThan"
	AlertOperatorGreaterThanOrEqualTo AlertOperator = "GreaterThanOrEqualTo"
	AlertOperatorLessThan             AlertOperator = "LessThan"
	AlertOperatorLessThanOrEqualTo    AlertOperator = "LessThanOrEqualTo"
	AlertOperatorNone                 AlertOperator = "None"
)

func PossibleValuesForAlertOperator() []string {
	return []string{
		string(AlertOperatorEqualTo),
		string(AlertOperatorGreaterThan),
		string(AlertOperatorGreaterThanOrEqualTo),
		string(AlertOperatorLessThan),
		string(AlertOperatorLessThanOrEqualTo),
		string(AlertOperatorNone),
	}
}

func parseAlertOperator(input string) (*AlertOperator, error) {
	vals := map[string]AlertOperator{
		"equalto":              AlertOperatorEqualTo,
		"greaterthan":          AlertOperatorGreaterThan,
		"greaterthanorequalto": AlertOperatorGreaterThanOrEqualTo,
		"lessthan":             AlertOperatorLessThan,
		"lessthanorequalto":    AlertOperatorLessThanOrEqualTo,
		"none":                 AlertOperatorNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertOperator(input)
	return &out, nil
}

type AlertSource string

const (
	AlertSourcePreset AlertSource = "Preset"
	AlertSourceUser   AlertSource = "User"
)

func PossibleValuesForAlertSource() []string {
	return []string{
		string(AlertSourcePreset),
		string(AlertSourceUser),
	}
}

func parseAlertSource(input string) (*AlertSource, error) {
	vals := map[string]AlertSource{
		"preset": AlertSourcePreset,
		"user":   AlertSourceUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertSource(input)
	return &out, nil
}

type AlertStatus string

const (
	AlertStatusActive     AlertStatus = "Active"
	AlertStatusDismissed  AlertStatus = "Dismissed"
	AlertStatusNone       AlertStatus = "None"
	AlertStatusOverridden AlertStatus = "Overridden"
	AlertStatusResolved   AlertStatus = "Resolved"
)

func PossibleValuesForAlertStatus() []string {
	return []string{
		string(AlertStatusActive),
		string(AlertStatusDismissed),
		string(AlertStatusNone),
		string(AlertStatusOverridden),
		string(AlertStatusResolved),
	}
}

func parseAlertStatus(input string) (*AlertStatus, error) {
	vals := map[string]AlertStatus{
		"active":     AlertStatusActive,
		"dismissed":  AlertStatusDismissed,
		"none":       AlertStatusNone,
		"overridden": AlertStatusOverridden,
		"resolved":   AlertStatusResolved,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertStatus(input)
	return &out, nil
}

type AlertTimeGrainType string

const (
	AlertTimeGrainTypeAnnually       AlertTimeGrainType = "Annually"
	AlertTimeGrainTypeBillingAnnual  AlertTimeGrainType = "BillingAnnual"
	AlertTimeGrainTypeBillingMonth   AlertTimeGrainType = "BillingMonth"
	AlertTimeGrainTypeBillingQuarter AlertTimeGrainType = "BillingQuarter"
	AlertTimeGrainTypeMonthly        AlertTimeGrainType = "Monthly"
	AlertTimeGrainTypeNone           AlertTimeGrainType = "None"
	AlertTimeGrainTypeQuarterly      AlertTimeGrainType = "Quarterly"
)

func PossibleValuesForAlertTimeGrainType() []string {
	return []string{
		string(AlertTimeGrainTypeAnnually),
		string(AlertTimeGrainTypeBillingAnnual),
		string(AlertTimeGrainTypeBillingMonth),
		string(AlertTimeGrainTypeBillingQuarter),
		string(AlertTimeGrainTypeMonthly),
		string(AlertTimeGrainTypeNone),
		string(AlertTimeGrainTypeQuarterly),
	}
}

func parseAlertTimeGrainType(input string) (*AlertTimeGrainType, error) {
	vals := map[string]AlertTimeGrainType{
		"annually":       AlertTimeGrainTypeAnnually,
		"billingannual":  AlertTimeGrainTypeBillingAnnual,
		"billingmonth":   AlertTimeGrainTypeBillingMonth,
		"billingquarter": AlertTimeGrainTypeBillingQuarter,
		"monthly":        AlertTimeGrainTypeMonthly,
		"none":           AlertTimeGrainTypeNone,
		"quarterly":      AlertTimeGrainTypeQuarterly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertTimeGrainType(input)
	return &out, nil
}

type AlertType string

const (
	AlertTypeBudget         AlertType = "Budget"
	AlertTypeBudgetForecast AlertType = "BudgetForecast"
	AlertTypeCredit         AlertType = "Credit"
	AlertTypeGeneral        AlertType = "General"
	AlertTypeInvoice        AlertType = "Invoice"
	AlertTypeQuota          AlertType = "Quota"
	AlertTypeXCloud         AlertType = "xCloud"
)

func PossibleValuesForAlertType() []string {
	return []string{
		string(AlertTypeBudget),
		string(AlertTypeBudgetForecast),
		string(AlertTypeCredit),
		string(AlertTypeGeneral),
		string(AlertTypeInvoice),
		string(AlertTypeQuota),
		string(AlertTypeXCloud),
	}
}

func parseAlertType(input string) (*AlertType, error) {
	vals := map[string]AlertType{
		"budget":         AlertTypeBudget,
		"budgetforecast": AlertTypeBudgetForecast,
		"credit":         AlertTypeCredit,
		"general":        AlertTypeGeneral,
		"invoice":        AlertTypeInvoice,
		"quota":          AlertTypeQuota,
		"xcloud":         AlertTypeXCloud,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertType(input)
	return &out, nil
}

type ExternalCloudProviderType string

const (
	ExternalCloudProviderTypeExternalBillingAccounts ExternalCloudProviderType = "externalBillingAccounts"
	ExternalCloudProviderTypeExternalSubscriptions   ExternalCloudProviderType = "externalSubscriptions"
)

func PossibleValuesForExternalCloudProviderType() []string {
	return []string{
		string(ExternalCloudProviderTypeExternalBillingAccounts),
		string(ExternalCloudProviderTypeExternalSubscriptions),
	}
}

func parseExternalCloudProviderType(input string) (*ExternalCloudProviderType, error) {
	vals := map[string]ExternalCloudProviderType{
		"externalbillingaccounts": ExternalCloudProviderTypeExternalBillingAccounts,
		"externalsubscriptions":   ExternalCloudProviderTypeExternalSubscriptions,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExternalCloudProviderType(input)
	return &out, nil
}
