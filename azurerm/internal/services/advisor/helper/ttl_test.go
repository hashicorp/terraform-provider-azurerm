package helper

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAdvisorRecommendationID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *AdvisorTtl
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:  "TTL1",
			Input: "7",
			Expected: &AdvisorTtl{
				days: utils.Int32(7),
				times: &AdvisorTtlTime{
					hours:   utils.Int32(0),
					minutes: utils.Int32(0),
					seconds: utils.Int32(0)},
			},
		},
		{
			Name:  "TTL2",
			Input: "7.10",
			Expected: &AdvisorTtl{
				days: utils.Int32(7),
				times: &AdvisorTtlTime{
					hours:   utils.Int32(10),
					minutes: utils.Int32(0),
					seconds: utils.Int32(0)},
			},
		},
		{
			Name:  "TTL3",
			Input: "7.10:30",
			Expected: &AdvisorTtl{
				days: utils.Int32(7),
				times: &AdvisorTtlTime{
					hours:   utils.Int32(10),
					minutes: utils.Int32(30),
					seconds: utils.Int32(0)},
			},
		},
		{
			Name:  "TTL4",
			Input: "7.10:30:20",
			Expected: &AdvisorTtl{
				days: utils.Int32(7),
				times: &AdvisorTtlTime{
					hours:   utils.Int32(10),
					minutes: utils.Int32(30),
					seconds: utils.Int32(20)},
			},
		},
		{
			Name:  "TTL5",
			Input: "10:30",
			Expected: &AdvisorTtl{
				days: utils.Int32(0),
				times: &AdvisorTtlTime{
					hours:   utils.Int32(10),
					minutes: utils.Int32(30),
					seconds: utils.Int32(0)},
			},
		},
		{
			Name:  "TTL6",
			Input: "0.00:30",
			Expected: &AdvisorTtl{
				days: utils.Int32(0),
				times: &AdvisorTtlTime{
					hours:   utils.Int32(0),
					minutes: utils.Int32(30),
					seconds: utils.Int32(0)},
			},
		},
		{
			Name:  "TTL7",
			Input: "23.00:30",
			Expected: &AdvisorTtl{
				days: utils.Int32(23),
				times: &AdvisorTtlTime{
					hours:   utils.Int32(0),
					minutes: utils.Int32(30),
					seconds: utils.Int32(0)},
			},
		},
		{
			Name:     "Invalid hour",
			Input:    "0.:30",
			Expected: nil,
		},
		{
			Name:     "Invalid Time",
			Input:    "0::30",
			Expected: nil,
		},
		{
			Name:     "Invalid Time 2",
			Input:    "0:60:30",
			Expected: nil,
		},
		{
			Name:     "Invalid Time 3",
			Input:    "0:60:-1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ParseAdvisorSuppresionTTL(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if !actual.Equal(v.Expected) {
			t.Fatalf("Expected %s but got %s for Duration", v.Expected.toString(), actual.toString())
		}
	}
}
