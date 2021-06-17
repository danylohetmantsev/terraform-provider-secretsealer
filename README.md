terraform-provider-secretsealer
================================

A very barebones provider that exposes basic `secretsealer` functionality as a terraform data source.


### Usage

```HCL
terraform {
  required_providers {
    kubeseal = {
      source = "danylohetmantsev/secretsealer"
      version = "0.0.1"
    }
  }
}

provider "secretsealer" {
}

data "secretsealer_secret" "my_secret" {
  name = "my-secret"
  namespace = "my-namespace
  type = "Opaque"

  data = {
    key = "value"
  }
  certificate_path = "./cert.crt"
}
```


### Argument Reference

The following arguments are supported:
- `name` - Name of the secret, must be unique.
- `namespace` - Namespace defines the space within which name of the secret must be unique.
- `type` -  The secret type. ex: `Opaque`
- `data` - Key/value pairs to populate the secret
- `certificate_path` - Name of the SealedSecrets controller in the cluster
- `sealed_secret` - Encrypted SealedSecret
- `depends_on` - For specifying hidden dependencies.

*NOTE: All the arguments above are required*
