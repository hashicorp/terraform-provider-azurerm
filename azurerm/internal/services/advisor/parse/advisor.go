package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AdvisorRecommendationId struct {
	Name        string
	ResourceUri string
}

type AdvisorSuppressionId struct {
	Name               string
	RecommendationName string
	ResourceUri        string
}

func AdvisorRecommendationID(input string) (*AdvisorRecommendationId, error) {
	//recommendation ID is resourceuri/providers/Microsoft.Advisor/recommendations/recommendationID
	inputSplit := strings.Split(input, "/providers/Microsoft.Advisor/recommendations/")
	if len(inputSplit) != 2 || inputSplit[1] == "" {
		return nil, fmt.Errorf("advisor Recommendation ID was invalid")
	}

	if _, err := azure.ParseAzureResourceID(inputSplit[0]); err != nil {
		return nil, err
	}

	return &AdvisorRecommendationId{
		ResourceUri: inputSplit[0],
		Name:        inputSplit[1],
	}, nil
}

func AdvisorSuppressionID(input string) (*AdvisorSuppressionId, error) {
	//suppression ID is /resourceUri/providers/Microsoft.Advisor/recommendations/recommendationId/suppressions/suppressionName1
	inputSplit := strings.Split(input, "/providers/Microsoft.Advisor/recommendations/")
	if len(inputSplit) != 2 || inputSplit[1] == "" {
		return nil, fmt.Errorf("advisor Suppression ID was invalid")
	}

	if _, err := azure.ParseAzureResourceID(inputSplit[0]); err != nil {
		return nil, err
	}

	suppressionSplit := strings.Split(inputSplit[1], "/suppressions/")
	if len(suppressionSplit) != 2 || suppressionSplit[1] == "" {
		return nil, fmt.Errorf("advisor Suppression ID was invalid")
	}

	return &AdvisorSuppressionId{
		ResourceUri:        inputSplit[0],
		RecommendationName: suppressionSplit[0],
		Name:               suppressionSplit[1],
	}, nil
}
