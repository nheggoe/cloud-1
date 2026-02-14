package v1

import (
	"cloud1/endpoint/countryinfo"
)

const (
	ApiVersion = "v1"
	Prefix     = countryinfo.Prefix + "/" + ApiVersion
)
