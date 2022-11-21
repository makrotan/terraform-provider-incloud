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

type IncloudApp struct {
    Id string `json:"id"`
    Name string `json:"name"`
    GitUrl string `json:"git_url"`
    Branch string `json:"branch"`
    Status string `json:"status"`
}

type IncloudAppResponse struct {
    Id string `json:"id"`
    Name string `json:"name"`
    GitUrl string `json:"git_url"`
    Branch string `json:"branch"`
    Status string `json:"status"`
}

func resourceIncloudApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncloudAppCreate,
		ReadContext:   resourceIncloudAppRead,
		UpdateContext: resourceIncloudAppUpdate,
		DeleteContext: resourceIncloudAppDelete,
		Schema: map[string]*schema.Schema{
			"identifier": &schema.Schema{
                Type:     schema.TypeString,
                Required: true, Computed: false, Optional: false, ForceNew: true,
			},
			"name": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"git_url": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"branch": &schema.Schema{
                Type:     schema.TypeString,
                Default: "main",
				Optional: true,
				Required: false,
				Computed: false,
				ForceNew: false,
			},
			"status": &schema.Schema{
                Type:     schema.TypeString,
                Required: false, Computed: true, Optional: false, ForceNew: false,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceIncloudAppCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("identifier").(string)
	instance := IncloudApp{
        Id: d.Get("identifier").(string),
        Name: d.Get("name").(string),
        GitUrl: d.Get("git_url").(string),
        Branch: d.Get("branch").(string),
        Status: d.Get("status").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/tenant/%s/app/%s", strings.Trim(provider.HostURL, "/"), provider.tenant_id, d.Get("identifier").(string)), strings.NewReader(string(rb)))

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

	var incloudAppResponse IncloudAppResponse
	err = json.Unmarshal(body, &incloudAppResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("identifier", incloudAppResponse.Id)
    d.Set("name", incloudAppResponse.Name)
    d.Set("git_url", incloudAppResponse.GitUrl)
    d.Set("branch", incloudAppResponse.Branch)
    d.Set("status", incloudAppResponse.Status)

	return diags
}

func resourceIncloudAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/tenant/%s/app/%s", strings.Trim(provider.HostURL, "/"), provider.tenant_id, d.Get("identifier").(string)), nil)
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
		log.Printf("[WARN] incloud_app %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var incloudAppResponse IncloudAppResponse
	err = json.Unmarshal(body, &incloudAppResponse)
	//err = json.NewDecoder(resp.Body).Decode(IncloudAppResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("identifier", incloudAppResponse.Id)
    d.Set("name", incloudAppResponse.Name)
    d.Set("git_url", incloudAppResponse.GitUrl)
    d.Set("branch", incloudAppResponse.Branch)
    d.Set("status", incloudAppResponse.Status)

	return diags
}
func resourceIncloudAppUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceIncloudAppCreate(ctx, d, m)
}

func resourceIncloudAppDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/tenant/%s/app/%s", strings.Trim(provider.HostURL, "/"), provider.tenant_id, d.Get("identifier").(string)), nil)
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