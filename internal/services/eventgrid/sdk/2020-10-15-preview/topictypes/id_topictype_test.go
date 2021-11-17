package topictypes

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = TopicTypeId{}

func TestNewTopicTypeID(t *testing.T) {
	id := NewTopicTypeID("topicTypeValue")

	if id.TopicTypeName != "topicTypeValue" {
		t.Fatalf("Expected %q but got %q for Segment 'TopicTypeName'", id.TopicTypeName, "topicTypeValue")
	}
}

func TestFormatTopicTypeID(t *testing.T) {
	actual := NewTopicTypeID("topicTypeValue").ID()
	expected := "/providers/Microsoft.EventGrid/topicTypes/topicTypeValue"
	if actual != expected {
		t.Fatalf("Expected the Formatted ID to be %q but got %q", actual, expected)
	}
}

func TestParseTopicTypeID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *TopicTypeId
	}{
		{
			// Incomplete URI
			Input: "",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/providers",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/providers/Microsoft.EventGrid",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/providers/Microsoft.EventGrid/topicTypes",
			Error: true,
		},
		{
			// Valid URI
			Input: "/providers/Microsoft.EventGrid/topicTypes/topicTypeValue",
			Expected: &TopicTypeId{
				TopicTypeName: "topicTypeValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/providers/Microsoft.EventGrid/topicTypes/topicTypeValue/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseTopicTypeID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %+v", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.TopicTypeName != v.Expected.TopicTypeName {
			t.Fatalf("Expected %q but got %q for TopicTypeName", v.Expected.TopicTypeName, actual.TopicTypeName)
		}

	}
}

func TestParseTopicTypeIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *TopicTypeId
	}{
		{
			// Incomplete URI
			Input: "",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/providers",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/pRoViDeRs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/providers/Microsoft.EventGrid",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/pRoViDeRs/mIcRoSoFt.eVeNtGrId",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/providers/Microsoft.EventGrid/topicTypes",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/pRoViDeRs/mIcRoSoFt.eVeNtGrId/tOpIcTyPeS",
			Error: true,
		},
		{
			// Valid URI
			Input: "/providers/Microsoft.EventGrid/topicTypes/topicTypeValue",
			Expected: &TopicTypeId{
				TopicTypeName: "topicTypeValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/providers/Microsoft.EventGrid/topicTypes/topicTypeValue/extra",
			Error: true,
		},
		{
			// Valid URI (mIxEd CaSe since this is insensitive)
			Input: "/pRoViDeRs/mIcRoSoFt.eVeNtGrId/tOpIcTyPeS/tOpIcTyPeVaLuE",
			Expected: &TopicTypeId{
				TopicTypeName: "tOpIcTyPeVaLuE",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment - mIxEd CaSe since this is insensitive)
			Input: "/pRoViDeRs/mIcRoSoFt.eVeNtGrId/tOpIcTyPeS/tOpIcTyPeVaLuE/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseTopicTypeIDInsensitively(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %+v", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.TopicTypeName != v.Expected.TopicTypeName {
			t.Fatalf("Expected %q but got %q for TopicTypeName", v.Expected.TopicTypeName, actual.TopicTypeName)
		}

	}
}

func TestSegmentsForTopicTypeId(t *testing.T) {
	segments := TopicTypeId{}.Segments()
	if len(segments) == 0 {
		t.Fatalf("TopicTypeId has no segments")
	}

	uniqueNames := make(map[string]struct{}, 0)
	for _, segment := range segments {
		uniqueNames[segment.Name] = struct{}{}
	}
	if len(uniqueNames) != len(segments) {
		t.Fatalf("Expected the Segments to be unique but got %q unique segments and %d total segments", len(uniqueNames), len(segments))
	}
}
