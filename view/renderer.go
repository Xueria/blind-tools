package view

import "fyne.io/fyne/v2"

// baseRenderer implements the boilerplate of fyne.WidgetRenderer (Objects and
// Destroy) that Fyne's internal BaseRenderer would otherwise provide. It lets
// custom widgets in this package implement a renderer without importing
// fyne's internal packages.
type baseRenderer struct {
	objects []fyne.CanvasObject
}

// Destroy is a no-op hook.
func (r *baseRenderer) Destroy() {}

// Objects returns the child objects that should be drawn.
func (r *baseRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// SetObjects replaces the child objects.
func (r *baseRenderer) SetObjects(objects []fyne.CanvasObject) {
	r.objects = objects
}
