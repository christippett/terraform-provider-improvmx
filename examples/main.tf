resource "improvmx_domain" "test" {
  domain = "terraform-test.christippett.dev"
}

data "improvmx_check" "test" {
  domain = improvmx_domain.test.domain
  depends_on = [
    improvmx_domain.test
  ]
}

data "improvmx_dns" "test" {
  domain = improvmx_domain.test.domain
  depends_on = [
    improvmx_domain.test
  ]
}


output "domain_check" {
  value = data.improvmx_check.test
}
output "domain_dns" {
  value = data.improvmx_dns.test
}
