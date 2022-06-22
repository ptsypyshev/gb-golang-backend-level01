/*
 * SimpleTwitter
 *
 * Simple REST API for service like a Twitter
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type User struct {
	ID int64 `json:"id,omitempty"`

	Email string `json:"email,omitempty"`

	Password string `json:"password,omitempty"`

	Nickname string `json:"nickname,omitempty"`

	Firstname string `json:"firstname,omitempty"`

	Lastname string `json:"lastname,omitempty"`
}
