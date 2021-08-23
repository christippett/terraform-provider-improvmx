resource "improvmx_domain" "example" {
  domain = var.domain

  dynamic "alias" {
    for_each = var.aliases

    content {
      alias   = alias.key
      forward = join(",", alias.value)
    }
  }
}
