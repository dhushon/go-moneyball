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

type Period struct {

	Periods int32 `json:"periods"`

	CurrentPeriod int32 `json:"currentPeriod"`

	GameStatus *GameStatus `json:"gameStatus"`
}
