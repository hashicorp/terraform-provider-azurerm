#!/usr/bin/env bash

files=$(find ./azurerm -type f -name "*.go")
error=false

echo "==> Checking that Custom Timeouts are used..."

for f in $files; do
  if grep "ctx := meta." "$f" > /dev/null; then
    echo $f
    error=true
  fi
done

if $error; then
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
