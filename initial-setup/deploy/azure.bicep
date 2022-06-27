targetScope = 'subscription'

param location string = deployment().location
param resourcesPrefix string = ''
param spPrincipalId string

var namePrefix = 'randommatch'
var resourcesPrefixCalculated = empty(resourcesPrefix) ? '${namePrefix}' : resourcesPrefix
var resourceGroupName = '${resourcesPrefixCalculated}staterg'

module stateResourceGroup './azureResourceGroup.bicep' = {
  name: '${resourcesPrefixCalculated}-resourceGroupDeployment'
  params: {
    resourceGroupName: resourceGroupName
    location: location
  }
}

module stateStorageAccount './azureStorageAccount.bicep' = {
  name: '${resourcesPrefixCalculated}-storageAccountDeployment'
  params: {
    storageAccountName: '${resourcesPrefixCalculated}statest'
    spPrincipalId: spPrincipalId
  }
  scope: resourceGroup(resourceGroupName)
  dependsOn: [
    stateResourceGroup
  ]
}
