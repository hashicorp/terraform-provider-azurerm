#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


echo "==> Checking that Custom Timeouts are used..."

if grep -r -l --include='*.go' 'ctx := meta.' ./internal; then
  echo ""
  echo "------------------------------------------------"
  echo ""
  echo "The files listed above must use a Wrapped StopContext to enable Custom Timeouts."
  echo "You can do this by changing:"
  echo ""
  echo "> ctx := meta.(*clients.Client).StopContext"
  echo ""
  echo "to"
  echo ""
  echo "> ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)"
  echo "> defer cancel()"
  echo ""
  echo "where 'ForCreate', 'ForCreateUpdate', 'ForDelete', 'ForRead' and 'ForUpdate' are available"
  echo ""
  echo "and then configuring Timeouts on the resource:"
  echo ""
  echo "> return &schema.Resource{"
  echo ">   ..."
  echo ">   Timeouts: &schema.ResourceTimeout{"
  echo ">     Create: schema.DefaultTimeout(30 * time.Minute),"
  echo ">     Read:   schema.DefaultTimeout(5 * time.Minute),"
  echo ">     Update: schema.DefaultTimeout(30 * time.Minute),"
  echo ">     Delete: schema.DefaultTimeout(30 * time.Minute),"
  echo ">   },"
  echo ">   ..."
  echo "> }"
  exit 1
fi

exit 0
