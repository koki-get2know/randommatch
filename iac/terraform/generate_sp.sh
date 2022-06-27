#!/bin/bash
cat << EOF > azuresp.json
{
"clientId": "$1",
"clientSecret": "$2",
"subscriptionId": "$3",
"tenantId": "$4"
}
EOF
