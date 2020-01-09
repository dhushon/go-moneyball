/*
 * BoxScore API
 *
 * BoxScore API
 *
 * API version: 1.0.0
 * Contact: dan.hushon@dxc.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type Arena struct {

	ArenaID string `json:"arenaID,omitempty"`

	Name string `json:"name"`

	Address *Address `json:"address"`

	Url *interface{} `json:"url,omitempty"`

	Occupancy int32 `json:"occupancy,omitempty"`
}