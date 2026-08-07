// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gen

import (
	"go/format"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
)

// untypedListIR builds a minimal subscription-scoped untyped list IR modelled on
// azurerm_eventgrid_domain, whose flatten issues an additional API call.
func untypedListIR() *ir.ResourceIR {
	return &ir.ResourceIR{
		ProviderName:         "azurerm",
		ProviderGithubOrg:    "hashicorp",
		ServicePackage:       "eventgrid",
		ServiceName:          "EventGrid",
		ClientField:          "Domains",
		Name:                 "EventGridDomain",
		TerraformType:        "azurerm_eventgrid_domain",
		SDKPackage:           "domains",
		SDKImportPath:        "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/domains",
		ReadModel:            "Domain",
		IDPackage:            "domains",
		IDImportPath:         "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/domains",
		IDParseFunc:          "ParseDomainID",
		IDTypeName:           "DomainId",
		ConstructorFunc:      "resourceEventGridDomain",
		FlattenFunc:          "resourceEventGridDomainFlatten",
		Untyped:              true,
		IsListable:           true,
		ListBySubscriptionOp: "ListBySubscription",
	}
}

func TestRenderList_UntypedFlattenNeedsContext(t *testing.T) {
	res := untypedListIR()
	res.FlattenNeedsContext = true
	res.FlattenClientType = "DomainsClient"

	got := RenderListResource(res)
	if _, err := format.Source([]byte(got)); err != nil {
		t.Fatalf("rendered list is not valid Go: %v\n\n%s", err, got)
	}

	// The deadline is captured before the iterator and a fresh context is built
	// from it inside the closure.
	mustContain(t, got, "deadline, ok := ctx.Deadline()")
	mustContain(t, got, `sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")`)
	mustContain(t, got, "ctx, cancel := context.WithDeadline(context.Background(), deadline)")
	mustContain(t, got, "defer cancel()")
	// The flatten is called with the live context and the client.
	mustContain(t, got, "resourceEventGridDomainFlatten(ctx, client, rd, id, &item)")
}

func TestRenderList_UntypedNoContextByDefault(t *testing.T) {
	res := untypedListIR() // FlattenNeedsContext defaults to false

	got := RenderListResource(res)
	if _, err := format.Source([]byte(got)); err != nil {
		t.Fatalf("rendered list is not valid Go: %v\n\n%s", err, got)
	}

	// Without an additional API call the iterator keeps the simple shape.
	if strings.Contains(got, "ctx.Deadline()") {
		t.Errorf("did not expect a deadline refresh when the flatten needs no context")
	}
	if strings.Contains(got, "context.WithDeadline") {
		t.Errorf("did not expect context.WithDeadline when the flatten needs no context")
	}
	mustContain(t, got, "resourceEventGridDomainFlatten(rd, id, &item)")
}

func TestRenderList_UntypedListMethodOptions(t *testing.T) {
	res := untypedListIR()
	res.ListByResourceGroupOp = "ListByResourceGroup"
	res.ListBySubscriptionHasOptions = true
	res.ListByResourceGroupHasOptions = true

	got := RenderListResource(res)
	if _, err := format.Source([]byte(got)); err != nil {
		t.Fatalf("rendered list is not valid Go: %v\n\n%s", err, got)
	}

	// The SDK list methods take a trailing options argument, which must be
	// supplied from the resource's SDK package.
	mustContain(t, got, "commonids.NewSubscriptionID(subscriptionID), domains.DefaultListBySubscriptionOperationOptions())")
	mustContain(t, got, "data.ResourceGroupName.ValueString()), domains.DefaultListByResourceGroupOperationOptions())")
}

func TestRenderList_UntypedNoListMethodOptions(t *testing.T) {
	res := untypedListIR() // HasOptions flags default to false

	got := RenderListResource(res)
	if strings.Contains(got, "OperationOptions()") {
		t.Errorf("did not expect an options argument when the SDK method takes none")
	}
	mustContain(t, got, "commonids.NewSubscriptionID(subscriptionID))")
}

func mustContain(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Errorf("expected rendered output to contain:\n%s", needle)
	}
}
