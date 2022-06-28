resource "azurerm_log_analytics_workspace" "aks" {
    name                = local.log_analytics_workspace_name
    location            = azurerm_resource_group.aks.location
    resource_group_name = azurerm_resource_group.aks.name
    sku                 = local.log_analytics_workspace_sku
}

resource "azurerm_log_analytics_solution" "aks-containerinsights" {
    solution_name         = "ContainerInsights"
    location              = azurerm_log_analytics_workspace.aks.location
    resource_group_name   = azurerm_resource_group.aks.name
    workspace_resource_id = azurerm_log_analytics_workspace.aks.id
    workspace_name        = azurerm_log_analytics_workspace.aks.name

    plan {
        publisher = "Microsoft"
        product   = "OMSGallery/ContainerInsights"
    }
}