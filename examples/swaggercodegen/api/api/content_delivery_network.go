/*
 * Service Provider (swaggercodegen)
 *
 * This service provider allows the creation of fake 'cdns' resources
 *
 * API version: 1.0.0
 * Contact: apiteam@serviceprovider.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package api

type ContentDeliveryNetwork struct {

	Id string `json:"id,omitempty"`

	Label string `json:"label"`

	Ips []string `json:"ips"`

	Hostnames []string `json:"hostnames"`

	ExampleInt int32 `json:"exampleInt,omitempty"`

	ExampleNumber float32 `json:"exampleNumber,omitempty"`

	ExampleBoolean bool `json:"example_boolean,omitempty"`
}
