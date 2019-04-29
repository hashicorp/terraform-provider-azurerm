#!/bin/bash

SERVER="ci-oss.hashicorp.engineering"

if [ $# -ne 2 ]; then
    echo "use: tcazure PR_NUMBER [TEST_PATTERN]"
    exit 1
fi

PR=$1
TESTS=$2

BUILD=$(cat <<-EOF
<build>
    <buildType id="Azure_ProviderMicrosoftAzureRM"/>
    <properties>
        <property name="BRANCH_NAME" value="refs/pull/$PR/merge"/>
        <property name="TEST_PATTERN" value="$TESTS"/>
    </properties>
</build>
EOF
)

curl -su "katbyte:$TCPASS" --request POST  "https://$SERVER/app/rest/buildQueue" --header "Content-Type:application/xml" --data-binary "$BUILD" > /dev/null
