package preflight

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	preflightvalidation "github.com/hashicorp/terraform-provider-azurerm/internal/preflight/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type ValidationRequest struct {
	Location   *string                                               `json:"location"`
	Provider   string                                                `json:"provider"`
	ResourceId resourceids.ResourceId                                `json:"resourceId"`
	Type       string                                                `json:"type"`
	Resource   preflightvalidation.ResourceValidationRequestResource `json:"resource"`
}

func (v ValidationRequest) ValidateResource(ctx context.Context, metadata sdk.ResourceMetaData) error {
	client := metadata.Client.Preflight.PreflightClient

	input := preflightvalidation.ResourceValidationRequest{
		Location:       v.Location,
		Provider:       v.Provider,
		Resources:      []preflightvalidation.ResourceValidationRequestResource{v.Resource},
		Scope:          v.decodeScope(),
		Type:           v.Type,
		ValidationType: pointer.To(preflightvalidation.ResourceValidationTypeArmFull),
	}

	resp, err := client.ValidateResources(ctx, input)
	if err != nil {
		if errorDetail := extractErrorDetail(resp.HttpResponse); errorDetail != nil {
			return errors.New(*errorDetail)
		}

		return err
	}

	if resp.Model == nil {
		return errors.New("missing model in validate response")
	}

	model := resp.Model

	if len(model.Properties.ValidatedResources) < 1 {
		return errors.New("validation did not return an error but there were no validated resources")
	}

	return nil
}

func (v ValidationRequest) decodeScope() string {
	// TODO - make this less hacky, plus it will only work for resources that are members of an RG
	resultFmt := "/subscriptions/%s/resourceGroups/%s"
	subscriptionID := ""
	resourceGroupName := ""
	segments := v.ResourceId.Segments()
	id := resourceids.NewParserFromResourceIdType(v.ResourceId)
	parsed, _ := id.Parse(v.ResourceId.ID(), true)

	for _, segment := range segments {
		switch segment.Type {
		case resourceids.SubscriptionIdSegmentType:
			s, _ := parsed.SegmentNamed(segment.Name, true)
			subscriptionID = pointer.From(s)
		case resourceids.ResourceGroupSegmentType:
			r, _ := parsed.SegmentNamed(segment.Name, true)
			resourceGroupName = pointer.From(r)
		}
	}

	return fmt.Sprintf(resultFmt, subscriptionID, resourceGroupName)
}

func extractErrorDetail(resp *http.Response) *string {
	type errorResponse struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Details []struct {
				Code    string `json:"code"`
				Target  string `json:"target"`
				Message string `json:"message"`
			}
		} `json:"error"`
	}

	if resp == nil {
		return nil
	}

	var errorResp errorResponse

	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		return nil // if there's no detail provided (some services don't) just bubble the error
	}

	if len(errorResp.Error.Details) > 0 {
		return &errorResp.Error.Details[0].Message
	}

	return nil
}
