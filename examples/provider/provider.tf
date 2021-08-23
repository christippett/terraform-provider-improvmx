terraform {
  required_providers {
    improvmx = {
      source  = "christippett/improvmx"
      version = "0.0.1"
    }
  }
}

provider "improvmx" {
  api_key = "<API_KEY>"
}
