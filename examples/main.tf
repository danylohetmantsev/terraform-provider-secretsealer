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

  certificate = file("./cert.crt")

  data = {
    username = "sagasgas",
    password = "sfasgasg"
  }
}

resource local_file this {
  filename = "./sealed-secret.yaml"
  content  = data.secretsealer_secret.this.sealed_secret
}
