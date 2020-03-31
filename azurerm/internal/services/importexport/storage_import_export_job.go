package importexport

import (
	"github.com/Azure/azure-sdk-for-go/services/storageimportexport/mgmt/2016-11-01/storageimportexport"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const (
	ImportJobType string = "Import"
	ExportJobType string = "Export"
)

func expandArmJobReturnAddress(input []interface{}) *storageimportexport.ReturnAddress {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &storageimportexport.ReturnAddress{
		RecipientName:   utils.String(v["recipient_name"].(string)),
		StreetAddress1:  utils.String(v["street_address1"].(string)),
		StreetAddress2:  utils.String(v["street_address2"].(string)),
		City:            utils.String(v["city"].(string)),
		StateOrProvince: utils.String(v["state_or_province"].(string)),
		PostalCode:      utils.String(v["postal_code"].(string)),
		CountryOrRegion: utils.String(v["country_or_region"].(string)),
		Phone:           utils.String(v["phone"].(string)),
		Email:           utils.String(v["email"].(string)),
	}
}

func expandArmJobReturnShipping(input []interface{}) *storageimportexport.ReturnShipping {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &storageimportexport.ReturnShipping{
		CarrierName:          utils.String(v["carrier_name"].(string)),
		CarrierAccountNumber: utils.String(v["carrier_account_number"].(string)),
	}
}

func flattenArmJobReturnAddress(input *storageimportexport.ReturnAddress) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var city string
	if input.City != nil {
		city = *input.City
	}
	var countryOrRegion string
	if input.CountryOrRegion != nil {
		countryOrRegion = *input.CountryOrRegion
	}
	var email string
	if input.Email != nil {
		email = *input.Email
	}
	var phone string
	if input.Phone != nil {
		phone = *input.Phone
	}
	var postalCode string
	if input.PostalCode != nil {
		postalCode = *input.PostalCode
	}
	var recipientName string
	if input.RecipientName != nil {
		recipientName = *input.RecipientName
	}
	var streetAddress1 string
	if input.StreetAddress1 != nil {
		streetAddress1 = *input.StreetAddress1
	}
	var stateOrProvince string
	if input.StateOrProvince != nil {
		stateOrProvince = *input.StateOrProvince
	}
	var streetAddress2 string
	if input.StreetAddress2 != nil {
		streetAddress2 = *input.StreetAddress2
	}
	return []interface{}{
		map[string]interface{}{
			"city":              city,
			"country_or_region": countryOrRegion,
			"email":             email,
			"phone":             phone,
			"postal_code":       postalCode,
			"recipient_name":    recipientName,
			"street_address1":   streetAddress1,
			"street_address2":   streetAddress2,
			"state_or_province": stateOrProvince,
		},
	}
}

func flattenArmJobReturnShipping(input *storageimportexport.ReturnShipping) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var carrierAccountNumber string
	if input.CarrierAccountNumber != nil {
		carrierAccountNumber = *input.CarrierAccountNumber
	}
	var carrierName string
	if input.CarrierName != nil {
		carrierName = *input.CarrierName
	}
	return []interface{}{
		map[string]interface{}{
			"carrier_account_number": carrierAccountNumber,
			"carrier_name":           carrierName,
		},
	}
}

func flattenArmJobShippingInformation(input *storageimportexport.ShippingInformation) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var city string
	if input.City != nil {
		city = *input.City
	}
	var countryOrRegion string
	if input.CountryOrRegion != nil {
		countryOrRegion = *input.CountryOrRegion
	}
	var postalCode string
	if input.PostalCode != nil {
		postalCode = *input.PostalCode
	}
	var recipientName string
	if input.RecipientName != nil {
		recipientName = *input.RecipientName
	}
	var stateOrProvince string
	if input.StateOrProvince != nil {
		stateOrProvince = *input.StateOrProvince
	}
	var streetAddress1 string
	if input.StreetAddress1 != nil {
		streetAddress1 = *input.StreetAddress1
	}
	var phone string
	if input.Phone != nil {
		phone = *input.Phone
	}
	var streetAddress2 string
	if input.StreetAddress2 != nil {
		streetAddress2 = *input.StreetAddress2
	}
	return []interface{}{
		map[string]interface{}{
			"city":              city,
			"country_or_region": countryOrRegion,
			"postal_code":       postalCode,
			"recipient_name":    recipientName,
			"state_or_province": stateOrProvince,
			"street_address1":   streetAddress1,
			"phone":             phone,
			"street_address2":   streetAddress2,
		},
	}
}
