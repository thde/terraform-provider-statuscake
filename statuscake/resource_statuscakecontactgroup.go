package statuscake

import (
	"fmt"
	"log"
	"strconv"

	"github.com/DreamItGetIT/statuscake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//nolint:errcheck
func resourceStatusCakeContactGroup() *schema.Resource {
	return &schema.Resource{
		Create: CreateContactGroup,
		Update: UpdateContactGroup,
		Delete: DeleteContactGroup,
		Read:   ReadContactGroup,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"contact_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"desktop_alert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ping_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pushover": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"boxcar": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mobiles": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"emails": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

//nolint:errcheck
func CreateContactGroup(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	newContactGroup := &statuscake.ContactGroup{
		GroupName:    d.Get("group_name").(string),
		Emails:       castSetToSliceStrings(d.Get("emails").(*schema.Set).List()),
		Mobiles:      d.Get("mobiles").(string),
		Boxcar:       d.Get("boxcar").(string),
		Pushover:     d.Get("pushover").(string),
		DesktopAlert: d.Get("desktop_alert").(string),
		PingURL:      d.Get("ping_url").(string),
	}

	log.Printf("[DEBUG] Creating new StatusCake Contact group: %s", d.Get("group_name").(string))

	response, err := statuscake.NewContactGroups(client).Create(newContactGroup)
	if err != nil {
		return fmt.Errorf("error creating statuscake contact group: %w", err)
	}

	d.Set("mobiles", newContactGroup.Mobiles)
	d.Set("boxcar", newContactGroup.Boxcar)
	d.Set("pushover", newContactGroup.Pushover)
	d.Set("desktop_alert", newContactGroup.DesktopAlert)
	d.Set("contact_id", newContactGroup.ContactID)
	d.SetId(strconv.Itoa(response.ContactID))

	return ReadContactGroup(d, meta)
}

//nolint:errcheck
func UpdateContactGroup(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	params := &statuscake.ContactGroup{
		GroupName:    d.Get("group_name").(string),
		Emails:       castSetToSliceStrings(d.Get("emails").(*schema.Set).List()),
		Mobiles:      d.Get("mobiles").(string),
		ContactID:    d.Get("contact_id").(int),
		Boxcar:       d.Get("boxcar").(string),
		Pushover:     d.Get("pushover").(string),
		DesktopAlert: d.Get("desktop_alert").(string),
		PingURL:      d.Get("ping_url").(string),
	}
	log.Printf("[DEBUG] StatusCake ContactGroup Update for %s", d.Id())
	_, err := statuscake.NewContactGroups(client).Update(params)
	d.Set("mobiles", params.Mobiles)
	d.Set("boxcar", params.Boxcar)
	d.Set("pushover", params.Pushover)
	d.Set("desktop_alert", params.DesktopAlert)
	if err != nil {
		return fmt.Errorf("error updating statuscake contact group: %w", err)
	}
	return ReadContactGroup(d, meta)
}

func DeleteContactGroup(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)
	id, _ := strconv.Atoi(d.Id())
	log.Printf("[DEBUG] Deleting StatusCake ContactGroup: %s", d.Id())
	err := statuscake.NewContactGroups(client).Delete(id)

	return err
}

//nolint:errcheck
func ReadContactGroup(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)
	id, _ := strconv.Atoi(d.Id())
	response, err := statuscake.NewContactGroups(client).Detail(id)
	if err != nil {
		return fmt.Errorf("error getting statuscake contact group details for %s: Error: %w", d.Id(), err)
	}
	d.Set("group_name", response.GroupName)
	d.Set("emails", response.Emails)
	d.Set("contact_id", response.ContactID)
	d.Set("ping_url", response.PingURL)
	d.SetId(strconv.Itoa(response.ContactID))

	return nil
}
