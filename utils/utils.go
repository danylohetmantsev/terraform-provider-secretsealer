package utils

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"os"
	"os/exec"
	"text/template"
)

var (
	secretManifestTemplate = `
apiVersion: v1
kind: Secret
type: {{ .Type }}
data:
{{ range $key, $value := .Data }}
  {{ $key }}: {{ $value -}}
{{- end }}
metadata:
  creationTimestamp: null
  name: {{ .Name }}
  namespace: {{ .Namespace }}
  {{- if .Labels }} 
  labels:
  {{ range $label, $label_value := .Labels }}
    {{ $label }}: {{ $label_value }}
  {{- end }}
  {{- end }}
`
)

type SecretManifest struct {
	Name      string
	Namespace string
	Data      map[string]interface{}
	Labels    map[string]interface{}
	Type      string
}

func Which(cmd string) string {
	p, _ := exec.LookPath(cmd)
	return p
}

func Log(message string) {
	log.Printf("[sealed_secrets_provider] ================= %s\n", message)
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// doesn't exist
		return false
	}

	return true
}

func GenerateSecretManifest(t_name string, t_namespace string, t_data map[string]interface{}, t_labels map[string]interface{}, t_type string) (io.Reader, error) {
	secretManifestYAML := new(bytes.Buffer)

	for k, v := range t_data {
		t_data[k] = base64.StdEncoding.EncodeToString([]byte(v.(string)))
	}

	secretManifest := SecretManifest{
		Name:      t_name,
		Namespace: t_namespace,
		Data:      t_data,
		Labels:    t_labels,
		Type:      t_type,
	}

	t := template.Must(template.New("secretManifestTemplate").Parse(secretManifestTemplate))
	err := t.Execute(secretManifestYAML, secretManifest)
	if err != nil {
		return nil, err
	}

	return secretManifestYAML, nil
}
