package statuscake

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/DreamItGetIT/statuscake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccStatusCakeSslBasic(t *testing.T) {
	var ssl statuscake.Ssl

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSslCheckDestroy(&ssl),
		Steps: []resource.TestStep{
			{
				Config: interpolateTerraformTemplateSsl(testAccSslConfigBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccSslCheckExists("statuscake_ssl.example", &ssl),
					testAccSslCheckAttributes("statuscake_ssl.example", &ssl),
				),
			},
		},
	})
}

func TestAccStatusCakeSslWithUpdate(t *testing.T) {
	var ssl statuscake.Ssl

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSslCheckDestroy(&ssl),
		Steps: []resource.TestStep{
			{
				Config: interpolateTerraformTemplateSsl(testAccSslConfigBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccSslCheckExists("statuscake_ssl.example", &ssl),
					testAccSslCheckAttributes("statuscake_ssl.example", &ssl),
				),
			},

			{
				Config: testAccSslConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccSslCheckExists("statuscake_ssl.example", &ssl),
					testAccSslCheckAttributes("statuscake_ssl.example", &ssl),
					resource.TestCheckResourceAttr("statuscake_ssl.example", "checkrate", "86400"),
					resource.TestCheckResourceAttr("statuscake_ssl.example", "domain", "https://www.example.com"),
					resource.TestCheckResourceAttr("statuscake_ssl.example", "contact_groups_c", ""),
					resource.TestCheckResourceAttr("statuscake_ssl.example", "alert_at", "18,81,2019"),
					resource.TestCheckResourceAttr("statuscake_ssl.example", "alert_reminder", "false"),
					resource.TestCheckResourceAttr("statuscake_ssl.example", "alert_expiry", "false"),
					resource.TestCheckResourceAttr("statuscake_ssl.example", "alert_broken", "true"),
					resource.TestCheckResourceAttr("statuscake_ssl.example", "alert_mixed", "false"),
				),
			},
		},
	})
}

func testAccSslCheckExists(rn string, ssl *statuscake.Ssl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ssl ID not set")
		}

		client := testAccProvider.Meta().(*statuscake.Client)
		sslID := rs.Primary.ID

		gotSsl, err := statuscake.NewSsls(client).Detail(sslID)
		if err != nil {
			return fmt.Errorf("error getting ssl: %w", err)
		}
		gotSsl.LastUpdatedUtc = "0000-00-00 00:00:00" // quick fix to avoid issue with it because the state is updated before the value change but it is changed when gotSsl is created
		*ssl = *gotSsl

		return nil
	}
}

func testAccSslCheckAttributes(rn string, ssl *statuscake.Ssl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs := s.RootModule().Resources[rn].Primary.Attributes

		check := func(key, stateValue, sslValue string) error {
			if sslValue != stateValue {
				return fmt.Errorf("different values for %s in state (%s) and in statuscake (%s)",
					key, stateValue, sslValue)
			}
			return nil
		}

		for key, value := range attrs {
			var err error

			switch key {
			case "domain":
				err = check(key, value, ssl.Domain)
			case "contact_groups_c":
				err = check(key, value, ssl.ContactGroupsC)
			case "checkrate":
				err = check(key, value, strconv.Itoa(ssl.Checkrate))
			case "alert_at":
				err = check(key, value, ssl.AlertAt)
			case "alert_reminder":
				err = check(key, value, strconv.FormatBool(ssl.AlertReminder))
			case "alert_expiry":
				err = check(key, value, strconv.FormatBool(ssl.AlertExpiry))
			case "alert_broken":
				err = check(key, value, strconv.FormatBool(ssl.AlertBroken))
			case "alert_mixed":
				err = check(key, value, strconv.FormatBool(ssl.AlertMixed))
			case "last_updated_utc":
				err = check(key, value, ssl.LastUpdatedUtc)
			case "paused":
				err = check(key, value, strconv.FormatBool(ssl.Paused))
			case "issuer_cn":
				err = check(key, value, ssl.IssuerCn)
			case "contact_groups":
				for _, tv := range ssl.ContactGroups {
					err = check(key, value, tv)
					if err != nil {
						return err
					}
				}
			case "cert_score":
				err = check(key, value, ssl.CertScore)
			case "cert_status":
				err = check(key, value, ssl.CertStatus)
			case "cipher":
				err = check(key, value, ssl.Cipher)
			case "valid_from_utc":
				err = check(key, value, ssl.ValidFromUtc)
			case "valid_until_utc":
				err = check(key, value, ssl.ValidUntilUtc)
			case "last_reminder":
				err = check(key, value, strconv.Itoa(ssl.LastReminder))
			case "flags":
				for _, tv := range ssl.Flags {
					err = check(key, value, strconv.FormatBool(tv))
					if err != nil {
						return err
					}
				}

			case "mixed_content":
				for _, tv := range ssl.MixedContent {
					for _, tv2 := range tv {
						err = check(key, value, tv2)
						if err != nil {
							return err
						}
					}
				}
			}
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccSslCheckDestroy(ssl *statuscake.Ssl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*statuscake.Client)
		_, err := statuscake.NewSsls(client).Detail(ssl.ID)
		if err == nil {
			return fmt.Errorf("ssl still exists")
		}

		return nil
	}
}

func interpolateTerraformTemplateSsl(template string) string {
	sslContactGroupID := "43402"

	if v := os.Getenv("STATUSCAKE_SSL_CONTACT_GROUP_ID"); v != "" {
		sslContactGroupID = v
	}
	if sslContactGroupID == "-1" {
		sslContactGroupID = ""
	}

	return fmt.Sprintf(template, sslContactGroupID)
}

const testAccSslConfigBasic = `
resource "statuscake_ssl" "example" {
	domain = "https://www.example.com"
	contact_groups_c = "%s"
        checkrate = 3600
        alert_at = "18,71,2019"
        alert_reminder = true
	alert_expiry = true
        alert_broken = false
        alert_mixed = true
}
`

const testAccSslConfigUpdate = `
resource "statuscake_ssl" "example" {
	domain = "https://www.example.com"
        contact_groups_c = ""
        checkrate = 86400
        alert_at = "18,81,2019"
        alert_reminder = false
	alert_expiry = false
        alert_broken = true
        alert_mixed = false
}
`
