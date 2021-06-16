terraform {
  required_providers {
    secretsealer ={
      source  = "danylohetmantsev/secretsealer"
    }
  }
}

provider secretsealer {}

data secretsealer_secret this {
  name      = "secret-example"
  namespace = "default"

  labels = {
    reflector-class = "tls"
  }

  certificate_path = "./cert.crt"

  data = {
    username = base64encode("User"),
    password = base64encode("SECRET_SO_SECRET")
  }

  type = "Opaque"
}

resource local_file this {
  filename = "./sealed-secret.yaml"
  content  = data.secretsealer_secret.this.sealed_secret
}
