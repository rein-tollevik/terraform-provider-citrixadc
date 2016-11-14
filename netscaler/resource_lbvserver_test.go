/*
Copyright 2016 Citrix Systems, Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package netscaler

import (
	"fmt"
	"github.com/chiradeep/go-nitro/netscaler"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccLbvserver_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLbvserverDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccLbvserver_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbvserverExist("netscaler_lbvserver.foo", nil),

					resource.TestCheckResourceAttr(
						"netscaler_lbvserver.foo", "ipv46", "10.202.11.11"),
					resource.TestCheckResourceAttr(
						"netscaler_lbvserver.foo", "lbmethod", "ROUNDROBIN"),
					resource.TestCheckResourceAttr(
						"netscaler_lbvserver.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"netscaler_lbvserver.foo", "persistencetype", "COOKIEINSERT"),
					resource.TestCheckResourceAttr(
						"netscaler_lbvserver.foo", "port", "443"),
					resource.TestCheckResourceAttr(
						"netscaler_lbvserver.foo", "servicetype", "SSL"),
				),
			},
		},
	})
}

func testAccCheckLbvserverExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No lb vserver name is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client
		data, err := nsClient.FindResource(netscaler.Lbvserver.Type(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("LB vserver %s not found", n)
		}

		return nil
	}
}

func testAccCheckLbvserverDestroy(s *terraform.State) error {
	nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netscaler_lbvserver" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := nsClient.FindResource(netscaler.Lbvserver.Type(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("LB vserver %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

const testAccLbvserver_basic = `



resource "netscaler_sslcertkey" "foossl" {
  certkey = "sample_ssl_cert_foo"
  cert = "/var/certs/server.crt"
  key = "/var/certs/server.key"
  notificationperiod = 40
  expirymonitor = "ENABLED"
}

resource "netscaler_lbvserver" "foo" {
  
  ipv46 = "10.202.11.11"
  lbmethod = "ROUNDROBIN"
  name = "terraform-lb"
  persistencetype = "COOKIEINSERT"
  port = 443
  servicetype = "SSL"
  sslcertkey = "${netscaler_sslcertkey.foossl.certkey}"
}
`
