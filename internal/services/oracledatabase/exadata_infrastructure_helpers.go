package oracledatabase

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudexadatainfrastructures"
)

func ConvertCustomerContactsToInternalModel(customerContactsList *[]cloudexadatainfrastructures.CustomerContact) []string {
	var customerContacts []string
	if customerContactsList != nil {
		for _, customerContact := range *customerContactsList {
			customerContacts = append(customerContacts, customerContact.Email)
		}
	}
	return customerContacts
}

func ConvertEstimatedPatchingTimesToInternalModel(estimatedPatchingTime *cloudexadatainfrastructures.EstimatedPatchingTime) []EstimatedPatchingTimeModel {
	if estimatedPatchingTime != nil {
		return []EstimatedPatchingTimeModel{
			{
				EstimatedDbServerPatchingTime:        estimatedPatchingTime.EstimatedDbServerPatchingTime,
				EstimatedNetworkSwitchesPatchingTime: estimatedPatchingTime.EstimatedNetworkSwitchesPatchingTime,
				EstimatedStorageServerPatchingTime:   estimatedPatchingTime.EstimatedStorageServerPatchingTime,
				TotalEstimatedPatchingTime:           estimatedPatchingTime.TotalEstimatedPatchingTime,
			},
		}
	}
	return nil
}

func ConvertMaintenanceWindowToInternalModel(maintenanceWindow *cloudexadatainfrastructures.MaintenanceWindow) []MaintenanceWindowModel {
	if maintenanceWindow != nil {
		return []MaintenanceWindowModel{
			{
				CustomActionTimeoutInMins:    pointer.From(maintenanceWindow.CustomActionTimeoutInMins),
				DaysOfWeek:                   ConvertDayOfWeekToInternalModel(maintenanceWindow.DaysOfWeek),
				HoursOfDay:                   pointer.From(maintenanceWindow.HoursOfDay),
				IsCustomActionTimeoutEnabled: pointer.From(maintenanceWindow.IsCustomActionTimeoutEnabled),
				IsMonthlyPatchingEnabled:     pointer.From(maintenanceWindow.IsMonthlyPatchingEnabled),
				LeadTimeInWeeks:              pointer.From(maintenanceWindow.LeadTimeInWeeks),
				Months:                       ConvertMonthsToInternalModel(maintenanceWindow.Months),
				PatchingMode:                 string(pointer.From(maintenanceWindow.PatchingMode)),
				Preference:                   string(pointer.From(maintenanceWindow.Preference)),
				WeeksOfMonth:                 pointer.From(maintenanceWindow.WeeksOfMonth),
			},
		}
	}
	return nil
}

func ConvertDayOfWeekToInternalModel(dayOfWeeks *[]cloudexadatainfrastructures.DayOfWeek) []string {
	var dayOfWeeksArray []string
	if dayOfWeeks != nil {
		for _, dayOfWeek := range *dayOfWeeks {
			dayOfWeeksArray = append(dayOfWeeksArray, string(dayOfWeek.Name))
		}
	}
	return dayOfWeeksArray
}

func ConvertMonthsToInternalModel(months *[]cloudexadatainfrastructures.Month) []string {
	var monthsArray []string
	if months != nil {
		for _, month := range *months {
			monthsArray = append(monthsArray, string(month.Name))
		}
	}
	return monthsArray
}

func ConvertCustomerContactsToSDK(customerContactsList []string) []cloudexadatainfrastructures.CustomerContact {
	var customerContacts []cloudexadatainfrastructures.CustomerContact
	if customerContactsList != nil {
		for _, customerContact := range customerContactsList {
			customerContacts = append(customerContacts, cloudexadatainfrastructures.CustomerContact{
				Email: customerContact,
			})
		}
	}
	return customerContacts
}

func ConvertDayOfWeekToSDK(daysOfWeek []string) []cloudexadatainfrastructures.DayOfWeek {
	var daysOfWeekConverted []cloudexadatainfrastructures.DayOfWeek
	for _, day := range daysOfWeek {
		daysOfWeekConverted = append(daysOfWeekConverted, cloudexadatainfrastructures.DayOfWeek{
			Name: cloudexadatainfrastructures.DayOfWeekName(day),
		})
	}
	return daysOfWeekConverted
}

func ConvertMonthsToSDK(months []string) []cloudexadatainfrastructures.Month {
	var monthsConverted []cloudexadatainfrastructures.Month
	for _, month := range months {
		monthsConverted = append(monthsConverted, cloudexadatainfrastructures.Month{
			Name: cloudexadatainfrastructures.MonthName(month),
		})
	}
	return monthsConverted
}
