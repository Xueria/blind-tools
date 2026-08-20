package res

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed assets/favicon.ico
var resourceIconData []byte

var Icon = &fyne.StaticResource{
	StaticName:    "assets/favicon.ico",
	StaticContent: resourceIconData,
}
