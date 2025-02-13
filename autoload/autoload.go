// Package autoload provides an auto-loading mechanism for the godotenvpgp package.
// Can be used in place of (godotenv)envfile.Load().
package autoload

import (
	"godotenvpgp/envfile"
)

var load = envfile.Load

func init() {
	load()
}
