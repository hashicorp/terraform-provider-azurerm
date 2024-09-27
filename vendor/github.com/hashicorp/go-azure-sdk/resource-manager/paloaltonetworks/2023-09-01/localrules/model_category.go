package localrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Category struct {
	Feeds     []string `json:"feeds"`
	UrlCustom []string `json:"urlCustom"`
}
