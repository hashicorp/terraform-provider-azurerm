package validate

import (
	"fmt"
)

func DataboxEdgeCountry(v interface{}, k string) (warnings []string, errors []error) {
	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	validCountries := getDataboxEdgeCountries()

	for _, str := range validCountries {
		if value == str {
			return warnings, errors
		}
	}

	errors = append(errors, fmt.Errorf("expected %q to be one of [%s], got %q", k, prettyErrorString(validCountries), value))

	return warnings, errors
}

func getDataboxEdgeCountries() []string {
	return []string{
		"Algeria",
		"Argentina",
		"Australia",
		"Austria",
		"Bahamas",
		"Bahrain",
		"Bangladesh",
		"Barbados",
		"Belgium",
		"Bermuda",
		"Bolivia",
		"Bosnia and Herzegovina",
		"Brazil",
		"Bulgaria",
		"Canada",
		"Cayman Islands",
		"Chile",
		"Colombia",
		"Costa Rica",
		"Croatia",
		"Cyprus",
		"Czechia",
		"CÃ´te D'ivoire",
		"Denmark",
		"Dominican Republic",
		"Ecuador",
		"Egypt",
		"El Salvador",
		"Estonia",
		"Ethiopia",
		"Finland",
		"France",
		"Georgia",
		"Germany",
		"Ghana",
		"Greece",
		"Guatemala",
		"Honduras",
		"Hong Kong SAR",
		"Hungary",
		"Iceland",
		"India",
		"Indonesia",
		"Ireland",
		"Israel",
		"Italy",
		"Jamaica",
		"Japan",
		"Jordan",
		"Kazakhstan",
		"Kenya",
		"Kuwait",
		"Kyrgyzstan",
		"Latvia",
		"Libya",
		"Liechtenstein",
		"Lithuania",
		"Luxembourg",
		"Macao SAR",
		"Malaysia",
		"Malta",
		"Mauritius",
		"Mexico",
		"Moldova",
		"Monaco",
		"Mongolia",
		"Montenegro",
		"Morocco",
		"Namibia",
		"Nepal",
		"Netherlands",
		"New Zealand",
		"Nicaragua",
		"Nigeria",
		"Norway",
		"Oman",
		"Pakistan",
		"Palestinian Authority",
		"Panama",
		"Paraguay",
		"Peru",
		"Philippines",
		"Poland",
		"Portugal",
		"Puerto Rico",
		"Qatar",
		"Republic of Korea",
		"Romania",
		"Russia",
		"Rwanda",
		"Saint Kitts And Nevis",
		"Saudi Arabia",
		"Senegal",
		"Serbia",
		"Singapore",
		"Slovakia",
		"Slovenia",
		"South Africa",
		"Spain",
		"Sri Lanka",
		"Sweden",
		"Switzerland",
		"Taiwan",
		"Tajikistan",
		"Tanzania",
		"Thailand",
		"Trinidad And Tobago",
		"Tunisia",
		"Turkey",
		"Turkmenistan",
		"U.S. Virgin Islands",
		"Uganda",
		"Ukraine",
		"United Arab Emirates",
		"United Kingdom",
		"United States",
		"Uruguay",
		"Uzbekistan",
		"Venezuela",
		"Vietnam",
		"Yemen",
		"Zambia",
		"Zimbabwe",
	}
}
