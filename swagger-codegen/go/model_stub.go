/*
 * Stubber API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type Stub struct {

	Jar string `json:"jar,omitempty"`

	Path string `json:"path"`

	State bool `json:"state,omitempty"`

	Pid string `json:"pid,omitempty"`

	Port string `json:"port,omitempty"`

	Cpu string `json:"cpu,omitempty"`

	Mem string `json:"mem,omitempty"`
}
