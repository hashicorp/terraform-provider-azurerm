package billing

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Customer -id=/providers/Microsoft.Billing/billingAccounts/123456/customers/123456
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=BillingProfile -id=/providers/Microsoft.Billing/billingAccounts/123456/billingProfiles/123456
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Department -id=/providers/Microsoft.Billing/billingAccounts/123456/departments/123456
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=BillingAccount -id=/providers/Microsoft.Billing/billingAccounts/123456
