variable "improvmx_api_key" {
  type = string
}

variable "domain" {
  type = string
}

variable "aliases" {
  type = map(list(string))
}
