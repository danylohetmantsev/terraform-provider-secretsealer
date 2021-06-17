package secretsealer

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"strings"

	secretsealer "terraform-provider-secretsealer/utils/kubeseal"

	"terraform-provider-secretsealer/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTemplateRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SealedSecret name.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Kubernetes namespace name.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Description: "A map of labels.",
				Optional:    true,
			},
			"data": {
				Type:        schema.TypeMap,
				Description: "A map of the secret data.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Opaque",
				Description: "Time in seconds to wait for any individual kubernetes operation.",
			},
			"certificate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Time in seconds to wait for any individual kubernetes operation.",
			},
			"sealed_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Time in seconds to wait for any individual kubernetes operation.",
			},
			"unsealed_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Time in seconds to wait for any individual kubernetes operation.",
			},
		},
	}
}

func dataTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	t_name := d.Get("name").(string)
	t_namespace := d.Get("namespace").(string)
	t_type := d.Get("type").(string)
	t_labels := d.Get("labels").(map[string]interface{})
	t_data := d.Get("data").(map[string]interface{})

	var diags diag.Diagnostics

	plaintext_certificate := strings.NewReader(d.Get("certificate").(string))
	bytes_certificate, err := ioutil.ReadAll(plaintext_certificate)
	if err != nil {
		return diag.FromErr(err)
	}

	var cert *x509.Certificate

	block, _ := pem.Decode([]byte(bytes_certificate))
	cert, _ = x509.ParseCertificate(block.Bytes)
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)

	secretManifest, err := utils.GenerateSecretManifest(t_name, t_namespace, t_data, t_labels, t_type)
	if err != nil {
		return diag.FromErr(err)
	}

	sealedSecretManifest, err := secretsealer.Seal(secretManifest, rsaPublicKey, 0, false)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("sealed_secret", sealedSecretManifest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(t_namespace + "/" + t_name)

	return diags
}
