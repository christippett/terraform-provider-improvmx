resource "improvmx_domain" "example" {
  domain = "piedpiper.com"

  alias {
    alias   = "hello"
    forward = "me@example.com"
  }
}
