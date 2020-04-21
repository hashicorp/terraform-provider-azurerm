package validate

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	iso8601 "github.com/btubbs/datetime"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/rickb777/date/period"
)

func ISO8601Duration(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if _, err := period.Parse(v); err != nil {
		errors = append(errors, err)
	}
	return warnings, errors
}

func ISO8601DurationBetween(min string, max string) func(i interface{}, k string) (warnings []string, errors []error) {
	minDuration := period.MustParse(min).DurationApprox()
	maxDuration := period.MustParse(max).DurationApprox()
	if minDuration >= maxDuration {
		panic(fmt.Sprintf("min duration (%v) >= max duration (%v)", minDuration, maxDuration))
	}
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
		}

		p, err := period.Parse(v)
		if err != nil {
			return nil, []error{err}
		}

		duration := p.DurationApprox()
		if duration < minDuration || duration > maxDuration {
			return nil, []error{fmt.Errorf("expected %s to be in the range (%v - %v), got %v", k, minDuration, maxDuration, duration)}
		}

		return nil, nil
	}
}

func ISO8601DateTime(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := iso8601.Parse(v, time.UTC); err != nil {
		errors = append(errors, fmt.Errorf("%q has the invalid ISO8601 date format %q: %+v", k, i, err))
	}

	return warnings, errors
}

// RFC3339 date is duration d or greater into the future
func RFC3339DateInFutureBy(d time.Duration) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
			return
		}

		t, err := date.ParseTime(time.RFC3339, v)
		if err != nil {
			errors = append(errors, fmt.Errorf("%q has the invalid RFC3339 date format %q: %+v", k, i, err))
			return
		}

		if time.Until(t) < d {
			errors = append(errors, fmt.Errorf("%q is %q and should be at least %q in the future", k, i, d))
		}

		return warnings, errors
	}
}

func AzureTimeZoneString() schema.SchemaValidateFunc {
	// List collected from https://support.microsoft.com/en-gb/help/973627/microsoft-time-zone-index-values
	// TODO look into programatic retrieval https://docs.microsoft.com/en-us/rest/api/maps/timezone/gettimezoneenumwindows
	validTimeZones := []string{
		"Afghanistan Standard Time",
		"Alaskan Standard Time",
		"Arab Standard Time",
		"Arabian Standard Time",
		"Arabic Standard Time",
		"Argentina Standard Time",
		"Atlantic Standard Time",
		"AUS Central Standard Time",
		"AUS Eastern Standard Time",
		"Azerbaijan Standard Time",
		"Azores Standard Time",
		"Bahia Standard Time",
		"Bangladesh Standard Time",
		"Belarus Standard Time",
		"Canada Central Standard Time",
		"Cape Verde Standard Time",
		"Caucasus Standard Time",
		"Cen. Australia Standard Time",
		"Central America Standard Time",
		"Central Asia Standard Time",
		"Central Brazilian Standard Time",
		"Central Europe Standard Time",
		"Central European Standard Time",
		"Central Pacific Standard Time",
		"Central Standard Time (Mexico)",
		"Central Standard Time",
		"China Standard Time",
		"Dateline Standard Time",
		"E. Africa Standard Time",
		"E. Australia Standard Time",
		"E. Europe Standard Time",
		"E. South America Standard Time",
		"Eastern Standard Time (Mexico)",
		"Eastern Standard Time",
		"Egypt Standard Time",
		"Ekaterinburg Standard Time",
		"Fiji Standard Time",
		"FLE Standard Time",
		"Georgian Standard Time",
		"GMT Standard Time",
		"Greenland Standard Time",
		"Greenwich Standard Time",
		"GTB Standard Time",
		"Hawaiian Standard Time",
		"India Standard Time",
		"Iran Standard Time",
		"Israel Standard Time",
		"Jordan Standard Time",
		"Kaliningrad Standard Time",
		"Korea Standard Time",
		"Libya Standard Time",
		"Line Islands Standard Time",
		"Magadan Standard Time",
		"Mauritius Standard Time",
		"Middle East Standard Time",
		"Montevideo Standard Time",
		"Morocco Standard Time",
		"Mountain Standard Time (Mexico)",
		"Mountain Standard Time",
		"Myanmar Standard Time",
		"N. Central Asia Standard Time",
		"Namibia Standard Time",
		"Nepal Standard Time",
		"New Zealand Standard Time",
		"Newfoundland Standard Time",
		"North Asia East Standard Time",
		"North Asia Standard Time",
		"Pacific SA Standard Time",
		"Pacific Standard Time (Mexico)",
		"Pacific Standard Time",
		"Pakistan Standard Time",
		"Paraguay Standard Time",
		"Romance Standard Time",
		"Russia Time Zone 10",
		"Russia Time Zone 11",
		"Russia Time Zone 3",
		"Russian Standard Time",
		"SA Eastern Standard Time",
		"SA Pacific Standard Time",
		"SA Western Standard Time",
		"Samoa Standard Time",
		"SE Asia Standard Time",
		"Singapore Standard Time",
		"South Africa Standard Time",
		"Sri Lanka Standard Time",
		"Syria Standard Time",
		"Taipei Standard Time",
		"Tasmania Standard Time",
		"Tokyo Standard Time",
		"Tonga Standard Time",
		"Turkey Standard Time",
		"Ulaanbaatar Standard Time",
		"US Eastern Standard Time",
		"US Mountain Standard Time",
		"UTC",
		"UTC+12",
		"UTC-02",
		"UTC-11",
		"Venezuela Standard Time",
		"Vladivostok Standard Time",
		"W. Australia Standard Time",
		"W. Central Africa Standard Time",
		"W. Europe Standard Time",
		"West Asia Standard Time",
		"West Pacific Standard Time",
		"Yakutsk Standard Time",
	}

	return validation.StringInSlice(validTimeZones, false)
}
