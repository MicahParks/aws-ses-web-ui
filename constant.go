package aswu

const (
	TemplateWrapper = "wrapper.gohtml"
)

const (
	PathAPICompose = "/api/compose"
)

type Path struct{}

func (p Path) APICompose() string {
	return PathAPICompose
}
