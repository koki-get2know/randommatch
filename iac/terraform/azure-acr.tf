resource "azurerm_container_registry" "container_registry" {
  name                = local.container_registry_name
  location            = azurerm_resource_group.aks.location
  resource_group_name = azurerm_resource_group.aks.name
  sku                 = "Standard"
  admin_enabled       = true
}