// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// Function returns the Function for a given name.
func (s *Server) Function(ctx context.Context, name string) (function.Function, *function.FuncError) {
	functionFuncs, diags := s.FunctionFuncs(ctx)

	funcErr := function.FuncErrorFromDiags(ctx, diags)

	functionFunc, ok := functionFuncs[name]

	if !ok {
		funcErr = function.ConcatFuncErrors(funcErr, function.NewFuncError(fmt.Sprintf("Function Not Found: No function named %q was found in the provider.", name)))

		return nil, funcErr
	}

	return functionFunc(), nil
}

// FunctionDefinition returns the Function Definition for the given name and
// caches the result for later Function operations.
func (s *Server) FunctionDefinition(ctx context.Context, name string) (function.Definition, *function.FuncError) {
	s.functionDefinitionsMutex.RLock()
	functionDefinition, ok := s.functionDefinitions[name]
	s.functionDefinitionsMutex.RUnlock()

	if ok {
		return functionDefinition, nil
	}

	functionImpl, funcErr := s.Function(ctx, name)

	if funcErr != nil {
		return function.Definition{}, funcErr
	}

	definitionReq := function.DefinitionRequest{}
	definitionResp := function.DefinitionResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined Function Definition method", map[string]interface{}{logging.KeyFunctionName: name})
	functionImpl.Definition(ctx, definitionReq, &definitionResp)
	logging.FrameworkTrace(ctx, "Called provider defined Function Definition method", map[string]interface{}{logging.KeyFunctionName: name})

	funcErr = function.ConcatFuncErrors(funcErr, function.FuncErrorFromDiags(ctx, definitionResp.Diagnostics))

	if funcErr != nil {
		return definitionResp.Definition, funcErr
	}

	s.functionDefinitionsMutex.Lock()

	if s.functionDefinitions == nil {
		s.functionDefinitions = make(map[string]function.Definition)
	}

	s.functionDefinitions[name] = definitionResp.Definition

	s.functionDefinitionsMutex.Unlock()

	return definitionResp.Definition, funcErr
}

// FunctionDefinitions returns a map of Function Definitions for the
// GetProviderSchema RPC without caching since not all definitions are
// guaranteed to be necessary for later provider operations. The definition
// implementations are also validated.
func (s *Server) FunctionDefinitions(ctx context.Context) (map[string]function.Definition, diag.Diagnostics) {
	functionDefinitions := make(map[string]function.Definition)

	functionFuncs, diags := s.FunctionFuncs(ctx)

	for name, functionFunc := range functionFuncs {
		functionImpl := functionFunc()

		definitionReq := function.DefinitionRequest{}
		definitionResp := function.DefinitionResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined Function Definition", map[string]interface{}{logging.KeyFunctionName: name})
		functionImpl.Definition(ctx, definitionReq, &definitionResp)
		logging.FrameworkTrace(ctx, "Called provider defined Function Definition", map[string]interface{}{logging.KeyFunctionName: name})

		diags.Append(definitionResp.Diagnostics...)

		if definitionResp.Diagnostics.HasError() {
			continue
		}

		validateReq := function.DefinitionValidateRequest{
			FuncName: name,
		}

		validateResp := function.DefinitionValidateResponse{}

		definitionResp.Definition.ValidateImplementation(ctx, validateReq, &validateResp)

		diags.Append(validateResp.Diagnostics...)

		if validateResp.Diagnostics.HasError() {
			continue
		}

		functionDefinitions[name] = definitionResp.Definition
	}

	return functionDefinitions, diags
}

// FunctionFuncs returns a map of Function functions. The results are cached
// on first use.
func (s *Server) FunctionFuncs(ctx context.Context) (map[string]func() function.Function, diag.Diagnostics) {
	logging.FrameworkTrace(ctx, "Checking FunctionTypes lock")
	s.functionFuncsMutex.Lock()
	defer s.functionFuncsMutex.Unlock()

	if s.functionFuncs != nil {
		return s.functionFuncs, s.functionFuncsDiags
	}

	s.functionFuncs = make(map[string]func() function.Function)

	provider, ok := s.Provider.(provider.ProviderWithFunctions)

	if !ok {
		// Only function-specific RPCs should return diagnostics about the
		// provider not implementing functions or missing functions.
		return s.functionFuncs, s.functionFuncsDiags
	}

	logging.FrameworkTrace(ctx, "Calling provider defined Provider Functions")
	functionFuncs := provider.Functions(ctx)
	logging.FrameworkTrace(ctx, "Called provider defined Provider Functions")

	for _, functionFunc := range functionFuncs {
		functionImpl := functionFunc()

		metadataReq := function.MetadataRequest{}
		metadataResp := function.MetadataResponse{}

		functionImpl.Metadata(ctx, metadataReq, &metadataResp)

		if metadataResp.Name == "" {
			s.functionFuncsDiags.AddError(
				"Function Name Missing",
				fmt.Sprintf("The %T Function returned an empty string from the Metadata method. ", functionImpl)+
					"This is always an issue with the provider and should be reported to the provider developers.",
			)
			continue
		}

		logging.FrameworkTrace(ctx, "Found function", map[string]interface{}{logging.KeyFunctionName: metadataResp.Name})

		if _, ok := s.functionFuncs[metadataResp.Name]; ok {
			s.functionFuncsDiags.AddError(
				"Duplicate Function Name Defined",
				fmt.Sprintf("The %s function name was returned for multiple functions. ", metadataResp.Name)+
					"Function names must be unique. "+
					"This is always an issue with the provider and should be reported to the provider developers.",
			)
			continue
		}

		s.functionFuncs[metadataResp.Name] = functionFunc
	}

	return s.functionFuncs, s.functionFuncsDiags
}

// FunctionMetadatas returns a slice of FunctionMetadata for the GetMetadata
// RPC.
func (s *Server) FunctionMetadatas(ctx context.Context) ([]FunctionMetadata, diag.Diagnostics) {
	functionFuncs, diags := s.FunctionFuncs(ctx)

	functionMetadatas := make([]FunctionMetadata, 0, len(functionFuncs))

	for name := range functionFuncs {
		functionMetadatas = append(functionMetadatas, FunctionMetadata{
			Name: name,
		})
	}

	return functionMetadatas, diags
}
