terraform {
  required_version = ">= 0.13"
  required_providers {
    improvmx = {
      source  = "christippett/improvmx"
      version = ">= 0.0.1"
    }
  }
}

provider "improvmx" {
  api_key = var.improvmx_api_key
}
