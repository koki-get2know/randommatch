resource "azurerm_resource_group" "aks" {
  name     = local.resource_group_name
  location = local.location
}