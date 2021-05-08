resource "improvmx_domain" "test" {
  domain = "christippett.dev"
}

data "improvmx_check" "test" {
  domain = improvmx_domain.test.domain
  depends_on = [
    improvmx_domain.test
  ]
}

output "domain_records" {
  value = data.improvmx_check.test
}
