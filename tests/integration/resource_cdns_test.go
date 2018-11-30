package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/dikhan/terraform-provider-openapi/examples/swaggercodegen/api/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"regexp"
)

const resourcePathCDN = "/v1/cdns"
const resouceSchemaDefinitionNameCDN = "ContentDeliveryNetworkV1"

const resourceCDNName = "cdn_v1"

var openAPIResourceNameCDN = fmt.Sprintf("%s_%s", providerName, resourceCDNName)
var openAPIResourceInstanceNameCDN = "my_cdn"
var openAPIResourceStateCDN = fmt.Sprintf("%s.%s", openAPIResourceNameCDN, openAPIResourceInstanceNameCDN)

var cdn api.ContentDeliveryNetworkV1
var testCreateConfigCDN string

func init() {
	// Setting this up here as it is used by many different tests
	cdn = newContentDeliveryNetwork("someLabel", []string{"192.168.0.2"}, []string{"www.google.com"}, 10, 12.22, true, "some updated message news", "some message news with details", "http", 80, "https", 443)
	testCreateConfigCDN = populateTemplateConfigurationCDN(cdn.Label, cdn.Ips, cdn.Hostnames, cdn.ExampleInt, cdn.ExampleNumber, cdn.ExampleBoolean, cdn.ObjectProperty.Message, cdn.ObjectProperty.DetailedMessage, cdn.ArrayOfObjectsExample[0].Protocol, cdn.ArrayOfObjectsExample[0].OriginPort, cdn.ArrayOfObjectsExample[1].Protocol, cdn.ArrayOfObjectsExample[1].OriginPort)
}

func TestAccCDN_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckCDNsV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConfigCDN,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistCDN(),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "label", cdn.Label),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.#", fmt.Sprintf("%d", len(cdn.Ips))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.0", arrayToString(cdn.Ips)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.#", fmt.Sprintf("%d", len(cdn.Hostnames))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.0", arrayToString(cdn.Hostnames)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_int", fmt.Sprintf("%d", cdn.ExampleInt)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "better_example_number_field_name", floatToString(cdn.ExampleNumber)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_boolean", fmt.Sprintf("%v", cdn.ExampleBoolean)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.%", fmt.Sprintf("%d", 2)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.message", cdn.ObjectProperty.Message),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.detailed_message", cdn.ObjectProperty.DetailedMessage),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_nested_scheme_property.%", fmt.Sprintf("%d", 1)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.#", fmt.Sprintf("%d", len(cdn.ArrayOfObjectsExample))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.origin_port", fmt.Sprint(cdn.ArrayOfObjectsExample[0].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.protocol", cdn.ArrayOfObjectsExample[0].Protocol),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.origin_port", fmt.Sprint(cdn.ArrayOfObjectsExample[1].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.protocol", cdn.ArrayOfObjectsExample[1].Protocol),
				),
			},
		},
	})
}

func TestAccCDN_Create_Using_Provider_Env_Variables(t *testing.T) {
	os.Setenv("APIKEY_AUTH", "apiKeyValue")
	testCDNCreateConfigWithoutProviderAuthProperty := fmt.Sprintf(`provider "%s" {
  x_request_id = "some value..."
}
resource "%s" "my_cdn" {
  label = "%s"
  ips = ["%s"]
  hostnames = ["%s"]
}`, providerName, openAPIResourceNameCDN, cdn.Label, arrayToString(cdn.Ips), arrayToString(cdn.Hostnames))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckCDNsV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCDNCreateConfigWithoutProviderAuthProperty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistCDN(),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "label", cdn.Label),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.#", fmt.Sprintf("%d", len(cdn.Ips))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.0", arrayToString(cdn.Ips)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.#", fmt.Sprintf("%d", len(cdn.Hostnames))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.0", arrayToString(cdn.Hostnames)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_nested_scheme_property.%", fmt.Sprintf("%d", 1)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_nested_scheme_property.name", "autogenerated name"),
				),
			},
		},
	})
	os.Unsetenv("APIKEY_AUTH")
}

func TestAccCDN_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckCDNsV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConfigCDN,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistCDN(),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "label", cdn.Label),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.#", fmt.Sprintf("%d", len(cdn.Ips))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.0", arrayToString(cdn.Ips)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.#", fmt.Sprintf("%d", len(cdn.Hostnames))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.0", arrayToString(cdn.Hostnames)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_int", fmt.Sprintf("%d", cdn.ExampleInt)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "better_example_number_field_name", floatToString(cdn.ExampleNumber)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_boolean", fmt.Sprintf("%v", cdn.ExampleBoolean)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.%", fmt.Sprintf("%d", 2)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.message", cdn.ObjectProperty.Message),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.detailed_message", cdn.ObjectProperty.DetailedMessage),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_nested_scheme_property.%", fmt.Sprintf("%d", 1)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_nested_scheme_property.name", "autogenerated name"),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.#", fmt.Sprintf("%d", len(cdn.ArrayOfObjectsExample))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.origin_port", fmt.Sprint(cdn.ArrayOfObjectsExample[0].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.protocol", cdn.ArrayOfObjectsExample[0].Protocol),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.origin_port", fmt.Sprint(cdn.ArrayOfObjectsExample[1].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.protocol", cdn.ArrayOfObjectsExample[1].Protocol),
				),
			},
			{
				Config:            testCreateConfigCDN,
				ResourceName:      openAPIResourceStateCDN,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistCDN(),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "label", cdn.Label),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.#", fmt.Sprintf("%d", len(cdn.Ips))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.0", arrayToString(cdn.Ips)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.#", fmt.Sprintf("%d", len(cdn.Hostnames))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.0", arrayToString(cdn.Hostnames)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_int", fmt.Sprintf("%d", cdn.ExampleInt)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "better_example_number_field_name", floatToString(cdn.ExampleNumber)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_boolean", fmt.Sprintf("%v", cdn.ExampleBoolean)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.%", fmt.Sprintf("%d", 2)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.message", cdn.ObjectProperty.Message),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.detailed_message", cdn.ObjectProperty.DetailedMessage),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_nested_scheme_property.%", fmt.Sprintf("%d", 1)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_nested_scheme_property.name", "autogenerated name"),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.#", fmt.Sprintf("%d", len(cdn.ArrayOfObjectsExample))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.origin_port", fmt.Sprint(cdn.ArrayOfObjectsExample[0].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.protocol", cdn.ArrayOfObjectsExample[0].Protocol),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.origin_port", fmt.Sprint(cdn.ArrayOfObjectsExample[1].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.protocol", cdn.ArrayOfObjectsExample[1].Protocol),
				),
			},
		},
	})
}

func TestAccCDN_CreateFailsDueToMissingMandatoryApiKeyAuth(t *testing.T) {
	testCDNCreateMissingAPIKeyAuthConfig := fmt.Sprintf(`provider "%s" {
  # apikey_auth = "apiKeyValue" simulating configuration that is missing the mandatory apikey_auth (commented out for the reference)
  x_request_id = "some value..."
}
resource "%s" "my_cdn" {}`, providerName, openAPIResourceNameCDN)

	expectedValidationError, _ := regexp.Compile(".*\"apikey_auth\": required field is not set.*")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckCDNsV1Destroy,
		Steps: []resource.TestStep{
			{
				Config:      testCDNCreateMissingAPIKeyAuthConfig,
				ExpectError: expectedValidationError,
			},
		},
	})
}

func TestAccCDN_CreateFailsDueToWrongAuthKeyValue(t *testing.T) {
	testCDNCreateWrongAPIKeyAuthConfig := fmt.Sprintf(`provider "%s" {
  apikey_auth = "This is not the key expected by the API to authenticate the client, it should be 'apiKeyValue'' :)"
  x_request_id = "some value..."
}
resource "%s" "my_cdn" {
  label = "%s"
  ips = ["%s"]
  hostnames = ["%s"]
}`, providerName, openAPIResourceNameCDN, cdn.Label, arrayToString(cdn.Ips), arrayToString(cdn.Hostnames))

	expectedValidationError, _ := regexp.Compile(".*{\"code\":\"401\", \"message\": \"unauthorized user\"}.*")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckCDNsV1Destroy,
		Steps: []resource.TestStep{
			{
				Config:      testCDNCreateWrongAPIKeyAuthConfig,
				ExpectError: expectedValidationError,
			},
		},
	})
}

func TestAccCDN_CreateFailsDueToRequiredPropertyMissing(t *testing.T) {
	testCDNCreateConfigMissingRequiredProperty := fmt.Sprintf(`provider "%s" {
  apikey_auth = "apiKeyValue"
  x_request_id = "some value..."
}
resource "%s" "my_cdn" {
  #label = "%s" # ==> Simulating required field is missing (commented out for the reference)
  ips = ["%s"]
  hostnames = ["%s"]
}`, providerName, openAPIResourceNameCDN, cdn.Label, arrayToString(cdn.Ips), arrayToString(cdn.Hostnames))

	expectedValidationError, _ := regexp.Compile(".*\"label\": required field is not set.*")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckCDNsV1Destroy,
		Steps: []resource.TestStep{
			{
				Config:      testCDNCreateConfigMissingRequiredProperty,
				ExpectError: expectedValidationError,
			},
		},
	})
}

func TestAccCDN_Update(t *testing.T) {
	var cdnUpdated = newContentDeliveryNetwork(cdn.Label, cdn.Ips, cdn.Hostnames, 14, 14.14, false, "some updated message news", "some message news with details", "http", 80, "https", 443)
	testCDNUpdatedConfig := populateTemplateConfigurationCDN(cdnUpdated.Label, cdnUpdated.Ips, cdnUpdated.Hostnames, cdnUpdated.ExampleInt, cdnUpdated.ExampleNumber, cdnUpdated.ExampleBoolean, cdnUpdated.ObjectProperty.Message, cdnUpdated.ObjectProperty.DetailedMessage, cdnUpdated.ArrayOfObjectsExample[0].Protocol, cdnUpdated.ArrayOfObjectsExample[0].OriginPort, cdnUpdated.ArrayOfObjectsExample[1].Protocol, cdnUpdated.ArrayOfObjectsExample[1].OriginPort)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckCDNsV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConfigCDN,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistCDN(),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "label", cdn.Label),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.#", fmt.Sprintf("%d", len(cdn.Ips))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.0", arrayToString(cdn.Ips)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.#", fmt.Sprintf("%d", len(cdn.Hostnames))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.0", arrayToString(cdn.Hostnames)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_int", fmt.Sprintf("%d", cdn.ExampleInt)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "better_example_number_field_name", floatToString(cdn.ExampleNumber)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_boolean", fmt.Sprintf("%v", cdn.ExampleBoolean)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.%", fmt.Sprintf("%d", 2)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.message", cdn.ObjectProperty.Message),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.detailed_message", cdn.ObjectProperty.DetailedMessage),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_nested_scheme_property.%", fmt.Sprintf("%d", 1)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_nested_scheme_property.name", "autogenerated name"),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.#", fmt.Sprintf("%d", len(cdn.ArrayOfObjectsExample))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.origin_port", fmt.Sprint(cdn.ArrayOfObjectsExample[0].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.protocol", cdn.ArrayOfObjectsExample[0].Protocol),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.origin_port", fmt.Sprint(cdn.ArrayOfObjectsExample[1].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.protocol", cdn.ArrayOfObjectsExample[1].Protocol),
				),
			},
			{
				Config: testCDNUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistCDN(),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "label", cdnUpdated.Label),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.#", fmt.Sprintf("%d", len(cdnUpdated.Ips))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.0", arrayToString(cdnUpdated.Ips)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.#", fmt.Sprintf("%d", len(cdnUpdated.Hostnames))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.0", arrayToString(cdnUpdated.Hostnames)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_int", fmt.Sprintf("%d", cdnUpdated.ExampleInt)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "better_example_number_field_name", floatToString(cdnUpdated.ExampleNumber)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_boolean", fmt.Sprintf("%v", cdnUpdated.ExampleBoolean)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.%", fmt.Sprintf("%d", 2)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.message", cdnUpdated.ObjectProperty.Message),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.detailed_message", cdnUpdated.ObjectProperty.DetailedMessage),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_nested_scheme_property.%", fmt.Sprintf("%d", 1)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_nested_scheme_property.name", "autogenerated name"),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.#", fmt.Sprintf("%d", len(cdnUpdated.ArrayOfObjectsExample))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.origin_port", fmt.Sprint(cdnUpdated.ArrayOfObjectsExample[0].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.protocol", cdnUpdated.ArrayOfObjectsExample[0].Protocol),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.origin_port", fmt.Sprint(cdnUpdated.ArrayOfObjectsExample[1].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.protocol", cdnUpdated.ArrayOfObjectsExample[1].Protocol),
				),
			},
		},
	})
}

func TestAccCDN_CreateWithZeroValues(t *testing.T) {
	var cdn = newContentDeliveryNetwork("label", []string{}, []string{}, 0, 0, false, "some message news", "some message news with details", "http", 80, "https", 443)
	testCDNZeroValuesConfig := populateTemplateConfigurationCDN(cdn.Label, cdn.Ips, cdn.Hostnames, cdn.ExampleInt, cdn.ExampleNumber, cdn.ExampleBoolean, cdn.ObjectProperty.Message, cdn.ObjectProperty.DetailedMessage, cdn.ArrayOfObjectsExample[0].Protocol, cdn.ArrayOfObjectsExample[0].OriginPort, cdn.ArrayOfObjectsExample[1].Protocol, cdn.ArrayOfObjectsExample[1].OriginPort)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckCDNsV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCDNZeroValuesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistCDN(),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "label", cdn.Label),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.#", "1"),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.0", arrayToString(cdn.Ips)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.#", "1"),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.0", arrayToString(cdn.Hostnames)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_int", fmt.Sprintf("%d", cdn.ExampleInt)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "better_example_number_field_name", "0"),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_boolean", fmt.Sprintf("%v", cdn.ExampleBoolean)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.#", fmt.Sprintf("%d", len(cdn.ArrayOfObjectsExample))),
				),
			},
		},
	})
}

func TestAccCDN_UpdateImmutableProperty(t *testing.T) {
	testCDNUpdatedImmutableConfig := populateTemplateConfigurationCDN("label updated", cdn.Ips, cdn.Hostnames, cdn.ExampleInt, cdn.ExampleNumber, cdn.ExampleBoolean, cdn.ObjectProperty.Message, cdn.ObjectProperty.DetailedMessage, cdn.ArrayOfObjectsExample[0].Protocol, cdn.ArrayOfObjectsExample[0].OriginPort, cdn.ArrayOfObjectsExample[1].Protocol, cdn.ArrayOfObjectsExample[1].OriginPort)
	expectedValidationError, _ := regexp.Compile(".*property label is immutable and therefore can not be updated. Update operation was aborted; no updates were performed.*")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckCDNsV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConfigCDN,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistCDN(),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "label", cdn.Label),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.#", fmt.Sprintf("%d", len(cdn.Ips))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.0", arrayToString(cdn.Ips)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.#", fmt.Sprintf("%d", len(cdn.Hostnames))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.0", arrayToString(cdn.Hostnames)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_int", fmt.Sprintf("%d", cdn.ExampleInt)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "better_example_number_field_name", floatToString(cdn.ExampleNumber)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_boolean", fmt.Sprintf("%v", cdn.ExampleBoolean)),
				),
			},
			{
				Config:      testCDNUpdatedImmutableConfig,
				ExpectError: expectedValidationError,
			},
		},
	})
}

func TestAccCDN_UpdateForceNewProperty(t *testing.T) {
	var cdnUpdatedForceNew = newContentDeliveryNetwork(cdn.Label, []string{"192.168.1.5"}, cdn.Hostnames, cdn.ExampleInt, cdn.ExampleNumber, cdn.ExampleBoolean, "some message news", "some message news with details", "http", 8080, "https", 8443)
	testCDNUpdatedForceNewConfig := populateTemplateConfigurationCDN(cdnUpdatedForceNew.Label, cdnUpdatedForceNew.Ips, cdnUpdatedForceNew.Hostnames, cdnUpdatedForceNew.ExampleInt, cdnUpdatedForceNew.ExampleNumber, cdnUpdatedForceNew.ExampleBoolean, cdnUpdatedForceNew.ObjectProperty.Message, cdnUpdatedForceNew.ObjectProperty.DetailedMessage, cdnUpdatedForceNew.ArrayOfObjectsExample[0].Protocol, cdnUpdatedForceNew.ArrayOfObjectsExample[0].OriginPort, cdnUpdatedForceNew.ArrayOfObjectsExample[1].Protocol, cdnUpdatedForceNew.ArrayOfObjectsExample[1].OriginPort)
	var originalID string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckCDNsV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConfigCDN,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						for _, res := range s.RootModule().Resources {
							if res.Type != openAPIResourceNameCDN {
								continue
							}
							originalID = res.Primary.ID
						}
						return nil
					},
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "label", cdn.Label),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.#", fmt.Sprintf("%d", len(cdn.Ips))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.0", arrayToString(cdn.Ips)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.#", fmt.Sprintf("%d", len(cdn.Hostnames))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.0", arrayToString(cdn.Hostnames)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_int", fmt.Sprintf("%d", cdn.ExampleInt)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "better_example_number_field_name", floatToString(cdn.ExampleNumber)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_boolean", fmt.Sprintf("%v", cdn.ExampleBoolean)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.%", fmt.Sprintf("%d", 2)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.message", cdn.ObjectProperty.Message),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_property.detailed_message", cdn.ObjectProperty.DetailedMessage),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_nested_scheme_property.%", fmt.Sprintf("%d", 1)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "object_nested_scheme_property.name", "autogenerated name"),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.#", fmt.Sprintf("%d", len(cdn.ArrayOfObjectsExample))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.origin_port", fmt.Sprint(cdn.ArrayOfObjectsExample[0].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.protocol", cdn.ArrayOfObjectsExample[0].Protocol),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.origin_port", fmt.Sprint(cdn.ArrayOfObjectsExample[1].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.protocol", cdn.ArrayOfObjectsExample[1].Protocol),
				),
			},
			{
				Config: testCDNUpdatedForceNewConfig,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						for _, res := range s.RootModule().Resources {
							if res.Type != openAPIResourceNameCDN {
								continue
							}
							// check that the ID generated in the first config apply has changed to a different one as the force new resource was required by the change applied
							forceNewID := res.Primary.ID
							if originalID == forceNewID {
								return fmt.Errorf("force new operation did not work, resource still has the same ID %s", originalID)
							}
						}
						resourceExistsFunc := testAccCheckResourceExistCDN()
						return resourceExistsFunc(s)
					},
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "label", cdnUpdatedForceNew.Label),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.#", fmt.Sprintf("%d", len(cdnUpdatedForceNew.Ips))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "ips.0", arrayToString(cdnUpdatedForceNew.Ips)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.#", fmt.Sprintf("%d", len(cdnUpdatedForceNew.Hostnames))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "hostnames.0", arrayToString(cdnUpdatedForceNew.Hostnames)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_int", fmt.Sprintf("%d", cdnUpdatedForceNew.ExampleInt)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "better_example_number_field_name", floatToString(cdnUpdatedForceNew.ExampleNumber)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "example_boolean", fmt.Sprintf("%v", cdnUpdatedForceNew.ExampleBoolean)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.#", fmt.Sprintf("%d", len(cdnUpdatedForceNew.ArrayOfObjectsExample))),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.origin_port", fmt.Sprint(cdnUpdatedForceNew.ArrayOfObjectsExample[0].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.0.protocol", cdnUpdatedForceNew.ArrayOfObjectsExample[0].Protocol),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.origin_port", fmt.Sprint(cdnUpdatedForceNew.ArrayOfObjectsExample[1].OriginPort)),
					resource.TestCheckResourceAttr(
						openAPIResourceStateCDN, "array_of_objects_example.1.protocol", cdnUpdatedForceNew.ArrayOfObjectsExample[1].Protocol),
				),
			},
		},
	})
}

func newContentDeliveryNetwork(label string, ips, hostnames []string, exampleInt int32, exampleNumber float32, exampleBool bool, objectPropertyMessage, objectDetailedMessage, listObjectProtocol string, listObjectOriginPort int32, listObject2Protocol string, listObject2OriginPort int32) api.ContentDeliveryNetworkV1 {
	return api.ContentDeliveryNetworkV1{
		Label:                 label,
		Ips:                   ips,
		Hostnames:             hostnames,
		ExampleInt:            exampleInt,
		ExampleNumber:         exampleNumber,
		ExampleBoolean:        exampleBool,
		ObjectProperty:        &api.ObjectProperty{Message: objectPropertyMessage, DetailedMessage: objectDetailedMessage},
		ArrayOfObjectsExample: []api.ContentDeliveryNetworkV1ArrayOfObjectsExample{{Protocol: listObjectProtocol, OriginPort: listObjectOriginPort}, {Protocol: listObject2Protocol, OriginPort: listObject2OriginPort}},
	}
}

func populateTemplateConfigurationCDN(label string, ips, hostnames []string, exampleInt int32, exampleNumber float32, exampleBool bool, objectPropertyMessage, objectDetailedMessage, listObjectProtocol string, listObjectOriginPort int32, listObject2Protocol string, listObject2OriginPort int32) string {
	return fmt.Sprintf(`provider "%s" {
  apikey_auth = "apiKeyValue"
  x_request_id = "some value..."
}

resource "%s" "%s" {
  label = "%s"
  ips = ["%s"]
  hostnames = ["%s"]

  example_int = %d
  better_example_number_field_name = %s
  example_boolean = %v

  object_property = {
    message = "%s"
    detailed_message = "%s"
  }

  array_of_objects_example = [
    {
      protocol = "%s"
      origin_port = %d
    },
    {
      protocol = "%s"
      origin_port = %d
    }
  ]
}`, providerName, openAPIResourceNameCDN, openAPIResourceInstanceNameCDN, label, arrayToString(ips), arrayToString(hostnames), exampleInt, floatToString(exampleNumber), exampleBool, objectPropertyMessage, objectDetailedMessage, listObjectProtocol, listObjectOriginPort, listObject2Protocol, listObject2OriginPort)
}

// Acceptance test resource-destruction for openapi_cdn_v1:
//
// Check all CDNs specified in the configuration have been destroyed.
func testCheckCDNsV1Destroy(state *terraform.State) error {
	return testCheckDestroy(state, openAPIResourceNameCDN, resourceCDNName, resourcePathCDN, resouceSchemaDefinitionNameCDN)
}

func testAccCheckResourceExistCDN() resource.TestCheckFunc {
	return testAccCheckResourceExist(openAPIResourceNameCDN, resourceCDNName, resourcePathCDN, resouceSchemaDefinitionNameCDN)
}
