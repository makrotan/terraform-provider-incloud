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

type IncloudProcess struct {
    Id string `json:"id"`
    Name string `json:"name"`
    ProjectId string `json:"project_id"`
    Spec string `json:"spec"`
}

type IncloudProcessResponse struct {
    Id string `json:"id"`
    Name string `json:"name"`
    ProjectId string `json:"project_id"`
    Spec string `json:"spec"`
}

func resourceIncloudProcess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncloudProcessCreate,
		ReadContext:   resourceIncloudProcessRead,
		UpdateContext: resourceIncloudProcessUpdate,
		DeleteContext: resourceIncloudProcessDelete,
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
			"project_id": &schema.Schema{
                Type:     schema.TypeString,
				Optional: false,
				Required: true,
				Computed: false,
				ForceNew: true,
			},
			"spec": &schema.Schema{
                Type:     schema.TypeString,
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

func resourceIncloudProcessCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Get("identifier").(string)
	instance := IncloudProcess{
        Id: d.Get("identifier").(string),
        Name: d.Get("name").(string),
        ProjectId: d.Get("project_id").(string),
        Spec: d.Get("spec").(string),
	}

	rb, err := json.Marshal(instance)
	if err != nil {
		return diag.FromErr(err)
	}

	// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/ssh-key/%s", strings.Trim(provider.HostURL, "/"), pk), strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/tenant/%s/project/%s/process/%s", strings.Trim(provider.HostURL, "/"), provider.tenant_id, d.Get("project_id").(string), d.Get("identifier").(string)), strings.NewReader(string(rb)))

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

	var incloudProcessResponse IncloudProcessResponse
	err = json.Unmarshal(body, &incloudProcessResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pk)
    d.Set("identifier", incloudProcessResponse.Id)
    d.Set("name", incloudProcessResponse.Name)
    d.Set("project_id", incloudProcessResponse.ProjectId)
    d.Set("spec", incloudProcessResponse.Spec)

	return diags
}

func resourceIncloudProcessRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	pk := d.Id()
	req, err := http.NewRequest("GET",  fmt.Sprintf("%s/api/v1/tenant/%s/project/%s/process/%s", strings.Trim(provider.HostURL, "/"), provider.tenant_id, d.Get("project_id").(string), d.Get("identifier").(string)), nil)
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
		log.Printf("[WARN] incloud_process %s not present", pk)
		d.SetId("")
		return nil
	} else if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("status: %d, body: %s", res.StatusCode, body))
	}

	var incloudProcessResponse IncloudProcessResponse
	err = json.Unmarshal(body, &incloudProcessResponse)
	//err = json.NewDecoder(resp.Body).Decode(IncloudProcessResponse)
	if err != nil {
		return diag.FromErr(err)
	}
    d.Set("identifier", incloudProcessResponse.Id)
    d.Set("name", incloudProcessResponse.Name)
    d.Set("project_id", incloudProcessResponse.ProjectId)
    d.Set("spec", incloudProcessResponse.Spec)

	return diags
}

func resourceIncloudProcessUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceIncloudProcessCreate(ctx, d, m)
}

func resourceIncloudProcessDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*Client)
	client := provider.HTTPClient

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

// 	pk := d.Id()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/tenant/%s/project/%s/process/%s", strings.Trim(provider.HostURL, "/"), provider.tenant_id, d.Get("project_id").(string), d.Get("identifier").(string)), nil)
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