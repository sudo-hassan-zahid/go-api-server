package constants

import "time"

const (
	ENV_LOCAL = "local"
	ENV_DEV   = "dev"
	ENV_PROD  = "prod"
)

var (
	JWTExpiration   = time.Minute * 15
	RefreshDuration = time.Hour * 24 * 7
)
