package brand

import _ "embed"

//go:embed HHHQ
var banner string

// Banner returns the interactive startup banner.
func Banner() string {
	return banner
}
