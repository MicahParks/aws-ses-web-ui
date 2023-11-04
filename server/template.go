package server

import (
	"net/http"

	hh "github.com/MicahParks/httphandle"

	aswu "github.com/MicahParks/aws-ses-web-ui"
)

type WrapperData struct {
	Title  string
	Path   aswu.Path
	Result hh.TemplateDataResult
}

func (w *WrapperData) SetResult(result hh.TemplateDataResult) {
	w.Result = result
}

func (s Server) ErrorTemplate(meta hh.TemplateRespMeta, r *http.Request, w http.ResponseWriter) {
	s.logger.Error("Error template not implemented.",
		"url", r.URL.String(),
	)
}

func (s Server) NotFound(w http.ResponseWriter, r *http.Request) {
	s.logger.Error("Not found handler not implemented.",
		"url", r.URL.String(),
	)
}

func (s Server) NewWrapperData(r *http.Request) (wData *WrapperData) {
	wData = &WrapperData{}
	return wData
}
