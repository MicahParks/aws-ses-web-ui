package template

import (
	"net/http"
	"net/mail"
	"strings"

	hh "github.com/MicahParks/httphandle"
	hhconst "github.com/MicahParks/httphandle/constant"

	aswu "github.com/MicahParks/aws-ses-web-ui"
	"github.com/MicahParks/aws-ses-web-ui/server"
)

type IndexData struct {
	DefaultFrom string
	WrapperData *server.WrapperData
}

type Index struct {
	s server.Server
}

func (i *Index) ApplyMiddleware(h http.Handler) http.Handler {
	return h
}
func (i *Index) Authorize(w http.ResponseWriter, r *http.Request) (authorized bool, modified *http.Request, skipTemplate bool) {
	return true, r, false
}
func (i *Index) Initialize(s server.Server) error {
	i.s = s
	return nil
}
func (i *Index) Respond(r *http.Request) (meta hh.TemplateRespMeta, templateData any, wrapperData hh.WrapperData) {
	meta.ResponseCode = http.StatusOK
	wData := i.s.NewWrapperData(r)
	wData.Title = "AWS SES Web UI"
	tData := &IndexData{
		WrapperData: wData,
	}
	if d := i.s.Conf.ASWU.DefaultFrom.Get(); d != nil {
		for _, v := range i.s.Conf.ASWU.AllowedFrom {
			if strings.HasPrefix(v, "@") {
				if strings.HasSuffix(d.Address, v) {
					tData.DefaultFrom = d.String()
					break
				}
			} else {
				addr, _ := mail.ParseAddress(v)
				if addr != nil && addr.Address == d.Address {
					if d.Name == "" {
						d.Name = addr.Name
					}
					tData.DefaultFrom = d.String()
					break
				}
			}
		}
		tData.DefaultFrom = d.String()
	}
	return meta, tData, wData
}
func (i *Index) TemplateName() string {
	return "index.gohtml"
}
func (i *Index) URLPattern() string {
	return hhconst.PathIndex
}
func (i *Index) WrapperTemplateName() string {
	return aswu.TemplateWrapper
}
