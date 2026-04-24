package preflight

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	Scope      string                                                `json:"scope"`
}

// NewValidationRequest constructs a new ValidationRequest. It dynamically extracts
// Provider, Type, and Resource.Type/Resource.Name from the provided resourceId.
func NewValidationRequest(location *string, id resourceids.ResourceId, apiVersion string, properties any) (ValidationRequest, error) {
	scope, provider, resourceType, resourceName, err := parseResourceId(id)
	if err != nil {
		return ValidationRequest{}, fmt.Errorf("parsing resource ID for preflight validation: %w", err)
	}

	return ValidationRequest{
		Location:   location,
		Provider:   provider,
		ResourceId: id,
		Type:       strings.TrimPrefix(resourceType, provider+"/"),
		Scope:      scope,
		Resource: preflightvalidation.ResourceValidationRequestResource{
			ApiVersion: apiVersion,
			Name:       resourceName,
			Type:       resourceType,
			Properties: properties,
		},
	}, nil
}

func (v ValidationRequest) ValidateResource(ctx context.Context, metadata sdk.ResourceMetaData) error {
	client := metadata.Client.Preflight.PreflightClient

	input := preflightvalidation.ResourceValidationRequest{
		Location:       v.Location,
		Provider:       v.Provider,
		Resources:      []preflightvalidation.ResourceValidationRequestResource{v.Resource},
		Scope:          v.Scope,
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

// parseResourceId breaks down an ARM resource ID into components needed for validation using the resourceids package.
func parseResourceId(id resourceids.ResourceId) (scope, provider, resourceType, resourceName string, err error) {
	parser := resourceids.NewParserFromResourceIdType(id)
	parsed, err := parser.Parse(id.ID(), true)
	if err != nil {
		return "", "", "", "", fmt.Errorf("parsing resource ID: %w", err)
	}

	segments := id.Segments()
	providerIdx := -1
	for i, s := range segments {
		if s.Type == resourceids.ResourceProviderSegmentType {
			providerIdx = i
			provider = *s.FixedValue
			break
		}
	}

	if providerIdx == -1 {
		return "", "", "", "", fmt.Errorf("resource ID is missing a resource provider segment")
	}

	var typeSegs []string
	var nameSegs []string

	for i := providerIdx + 1; i < len(segments); i++ {
		s := segments[i]
		if s.Type == resourceids.ConstantSegmentType || s.Type == resourceids.StaticSegmentType {
			if val, ok := parsed.SegmentNamed(s.Name, true); ok && val != nil {
				typeSegs = append(typeSegs, *val)
			} else if s.FixedValue != nil {
				typeSegs = append(typeSegs, *s.FixedValue)
			} else if s.PossibleValues != nil && len(*s.PossibleValues) > 0 {
				typeSegs = append(typeSegs, (*s.PossibleValues)[0])
			}
		} else if s.Type == resourceids.UserSpecifiedSegmentType {
			if val, ok := parsed.SegmentNamed(s.Name, true); ok && val != nil {
				nameSegs = append(nameSegs, *val)
			}
		}
	}

	resourceType = fmt.Sprintf("%s/%s", provider, strings.Join(typeSegs, "/"))
	resourceName = strings.Join(nameSegs, "/")

	var scopeSegments []string
	cutOffIndex := providerIdx - 1
	if len(typeSegs) > 1 {
		for i := len(segments) - 1; i >= 0; i-- {
			if segments[i].Type == resourceids.ConstantSegmentType || segments[i].Type == resourceids.StaticSegmentType {
				cutOffIndex = i
				break
			}
		}
	}

	for i := 0; i < cutOffIndex; i++ {
		s := segments[i]
		if val, ok := parsed.SegmentNamed(s.Name, true); ok && val != nil {
			scopeSegments = append(scopeSegments, *val)
		} else if s.FixedValue != nil {
			scopeSegments = append(scopeSegments, *s.FixedValue)
		} else if s.PossibleValues != nil && len(*s.PossibleValues) > 0 {
			scopeSegments = append(scopeSegments, (*s.PossibleValues)[0])
		}
	}

	if len(scopeSegments) > 0 {
		scope = "/" + strings.Join(scopeSegments, "/")
	}

	return scope, provider, resourceType, resourceName, nil
}

func extractErrorDetail(resp *http.Response) *string {
	if resp == nil || resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	type errorResponse struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Details []struct {
				Code    string `json:"code"`
				Target  string `json:"target"`
				Message string `json:"message"`
			} `json:"details"`
		} `json:"error"`
	}

	var errorResp errorResponse
	if err := json.Unmarshal(bodyBytes, &errorResp); err != nil {
		return nil // if there's no detail provided just bubble the original error
	}

	var messages []string
	if errorResp.Error.Message != "" {
		messages = append(messages, errorResp.Error.Message)
	}

	for _, d := range errorResp.Error.Details {
		if d.Message != "" {
			messages = append(messages, d.Message)
		}
	}

	if len(messages) > 0 {
		msg := strings.Join(messages, " | ")
		return &msg
	}

	return nil
}
