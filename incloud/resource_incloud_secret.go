package incloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type IncloudSecret struct {
    Id string `json:"id"`
    Data map[string]interface{} `json:"data"`
}

type IncloudSecretResponse struct {
    Id string `json:"id"`
    Data map[string]interface{} `json:"data"`
}

func resourceIncloudSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncloudSecretCreate,
		ReadContext:   resourceIncloudSecretRead,
		UpdateContext: resourceIncloudSecretUpdate,
		DeleteContext: resourceIncloudSecretDelete,
		Schema: map[string]*schema.Schema{
			"identifier": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
			},
			"data": &schema.Schema{
                Type:     schema.TypeMap,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
                Sensitive: true,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceIncloudSecretCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("identifier").(string)
	instance := IncloudSecret{
        Id: d.Get("identifier").(string),
        Data: d.Get("data").(map[string]interface{}),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/tenant/%s/secret/%s", strings.Trim(provider.HostURL, "/"), provider.tenant_id, d.Get("identifier").(string)), strings.NewReader(string(rb)))

	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.api_token))

	res, err := provider.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var incloudSecretResponse IncloudSecretResponse
	err = json.Unmarshal(body, &incloudSecretResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("identifier", incloudSecretResponse.Id)
    d.Set("data", incloudSecretResponse.Data)

	return diags
}

func resourceIncloudSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/tenant/%s/secret/%s", strings.Trim(provider.HostURL, "/"), provider.tenant_id, d.Get("identifier").(string)), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.api_token))

	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode == 404 {
		log.Printf("[WARN] incloud_secret %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var incloudSecretResponse IncloudSecretResponse
	err = json.Unmarshal(body, &incloudSecretResponse)
	//err = json.NewDecoder(resp.Body).Decode(IncloudSecretResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("identifier", incloudSecretResponse.Id)
    d.Set("data", incloudSecretResponse.Data)

	return diags
}
func resourceIncloudSecretUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceIncloudSecretCreate(ctx, d, m)
}

func resourceIncloudSecretDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/tenant/%s/secret/%s", strings.Trim(provider.HostURL, "/"), provider.tenant_id, d.Get("identifier").(string)), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.api_token))

	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode >= 300 {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	return diags
}