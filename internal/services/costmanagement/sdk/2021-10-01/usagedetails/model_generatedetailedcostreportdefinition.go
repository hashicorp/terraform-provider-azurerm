package usagedetails

type GenerateDetailedCostReportDefinition struct {
	BillingPeriod *string                               `json:"billingPeriod,omitempty"`
	CustomerId    *string                               `json:"customerId,omitempty"`
	InvoiceId     *string                               `json:"invoiceId,omitempty"`
	Metric        *GenerateDetailedCostReportMetricType `json:"metric,omitempty"`
	TimePeriod    *GenerateDetailedCostReportTimePeriod `json:"timePeriod,omitempty"`
}
