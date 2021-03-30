// +build go1.13

// Copyright 2017 Microsoft Corporation. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// Package armcore provides connections and utilities for Go SDK ARM client modules.
//
// All Azure Resource Manager clients require a Connection, which is simply a
// combination of the desired ARM endpoint and a pipeline for handling HTTP requests
// and responses.
//
// To access the Azure public cloud, use the NewDefaultConnection() constructor with
// the required token credential.  Module azidentity provides several methods for
// obtaining token credentials.
//
//		cred, _ := azidentity.NewDefaultAzureCredential(nil)
//		con := armcore.NewDefaultConnection(cred, nil)
//
// When accessing clouds other than the Azure public cloud, use the NewConnection()
// constructor with the required ARM endpoint and token credential.  The most common
// case is connecting to an Azure sovereign cloud or Azure Stack instance.
//
// NewDefaultConnection() and NewConnection() are configured with the same pipeline
// thus have the same pipeline configuration options.  Use the NewConnectionWithPipeline()
// constructor to create a connection that uses a custom azcore.Pipeline.  Note that
// any custom pipeline will require at minimum an authentication policy obtained from
// a token credential in order to authenticate with ARM.  See the implementation of
// NewConnection() for how to obtain a credential's authentication policy.
package armcore
