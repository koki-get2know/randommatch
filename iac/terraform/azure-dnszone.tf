resource "azurerm_dns_zone" "dns_zone" {
  name                = local.dns_zone
  resource_group_name = azurerm_resource_group.aks.name
}


resource "azurerm_dns_a_record" "dns_record" {
  name                = local.zone_prefix_name
  zone_name           = azurerm_dns_zone.dns_zone.name
  resource_group_name = azurerm_resource_group.aks.name
  ttl                 = 3600
  target_resource_id  = azurerm_public_ip.aks.id
}

resource "azurerm_dns_mx_record" "mx_record" {
  name                = local.zone_prefix_name
  zone_name           = azurerm_dns_zone.dns_zone.name
  resource_group_name = azurerm_resource_group.aks.name
  ttl                 = 3600

  record {
    preference = 10
    exchange   = "feedback-smtp.eu-west-1.amazonses.com"
  }
}

resource "azurerm_dns_txt_record" "txt_record" {
  name                = local.zone_prefix_name
  zone_name           = azurerm_dns_zone.dns_zone.name
  resource_group_name = azurerm_resource_group.aks.name
  ttl                 = 3600

  record {
    value = "v=spf1 include:amazonses.com ~all"
  }
}

resource "azurerm_dns_cname_record" "mailservice1" {
  name                = "3y5fuenaapsknlhelrsgmupajkrrq2a7._domainkey"
  zone_name           = azurerm_dns_zone.dns_zone.name
  resource_group_name = azurerm_resource_group.aks.name
  ttl                 = 3600
  record              = "3y5fuenaapsknlhelrsgmupajkrrq2a7.dkim.amazonses.com"
}

resource "azurerm_dns_cname_record" "mailservice2" {
  name                = "4ti5k6e2hxzqejgjeauyx57mab2rmqlh._domainkey"
  zone_name           = azurerm_dns_zone.dns_zone.name
  resource_group_name = azurerm_resource_group.aks.name
  ttl                 = 3600
  record              = "4ti5k6e2hxzqejgjeauyx57mab2rmqlh.dkim.amazonses.com"
}

resource "azurerm_dns_cname_record" "mailservice3" {
  name                = "7gihl7ksnx6xkcr67npoasrjvgrrtpnw._domainkey"
  zone_name           = azurerm_dns_zone.dns_zone.name
  resource_group_name = azurerm_resource_group.aks.name
  ttl                 = 3600
  record              = "7gihl7ksnx6xkcr67npoasrjvgrrtpnw.dkim.amazonses.com"
}

resource "azurerm_dns_cname_record" "mailservice4" {
  name                = "h4wbw7mq5s7t5wacdqk6d6klagccwuvw._domainkey"
  zone_name           = azurerm_dns_zone.dns_zone.name
  resource_group_name = azurerm_resource_group.aks.name
  ttl                 = 3600
  record              = "h4wbw7mq5s7t5wacdqk6d6klagccwuvw.dkim.amazonses.com"
}

resource "azurerm_dns_cname_record" "mailservice5" {
  name                = "pcu4qgi43ehxauzxrgfr4kgcvuz25bfb._domainkey"
  zone_name           = azurerm_dns_zone.dns_zone.name
  resource_group_name = azurerm_resource_group.aks.name
  ttl                 = 3600
  record              = "pcu4qgi43ehxauzxrgfr4kgcvuz25bfb.dkim.amazonses.com"
}

resource "azurerm_dns_cname_record" "mailservice6" {
  name                = "vnz7diohlp6do44rxxuj2nwbkgdyf7na._domainkey"
  zone_name           = azurerm_dns_zone.dns_zone.name
  resource_group_name = azurerm_resource_group.aks.name
  ttl                 = 3600
  record              = "vnz7diohlp6do44rxxuj2nwbkgdyf7na.dkim.amazonses.com"
}