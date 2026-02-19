package constants

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleViewer Role = "viewer"
)

var AllowedRoles = map[Role]bool{
	RoleAdmin:  true,
	RoleUser:   true,
	RoleViewer: true,
}
