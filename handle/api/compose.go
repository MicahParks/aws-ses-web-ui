package api

import (
	"context"
	"log/slog"
	"net/http"

	hhconst "github.com/MicahParks/httphandle/constant"
	"github.com/MicahParks/httphandle/middleware/ctxkey"

	aswu "github.com/MicahParks/aws-ses-web-ui"
	"github.com/MicahParks/aws-ses-web-ui/model"
	"github.com/MicahParks/aws-ses-web-ui/server"
)

type Compose struct {
	s server.Server
}

func (c *Compose) ApplyMiddleware(h http.Handler) http.Handler {
	return h
}
func (c *Compose) Authorize(w http.ResponseWriter, r *http.Request) (authorized bool, modified *http.Request) {
	return true, r
}
func (c *Compose) ContentType() (request, response string) {
	return hhconst.ContentTypeJSON, hhconst.ContentTypeJSON
}
func (c *Compose) HTTPMethod() string {
	return http.MethodPost
}
func (c *Compose) Initialize(s server.Server) error {
	c.s = s
	return nil
}
func (c *Compose) Respond(r *http.Request) (code int, body []byte, err error) {
	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()
	b, ctx, code, body, err := jsonBody[model.Compose](r)
	if err != nil {
		return code, body, nil
	}

	l := ctx.Value(ctxkey.Logger).(*slog.Logger)

	if c.s.Conf.ASWU.UsePostgres {
		tx1, err := c.s.BeginTx(ctx)
		if err != nil {
			l.ErrorContext(ctx, hhconst.MsgFailTransactionBegin,
				hhconst.LogErr, err,
				"txName", "tx1",
			)
			return errorResponse(ctx, http.StatusInternalServerError, hhconst.MsgFailTransactionCommit)
		}
		//goland:noinspection GoUnhandledErrorResult
		defer tx1.Rollback(ctx)

		ctx = context.WithValue(ctx, ctxkey.Tx, tx1)
		b, err = c.s.Postgres.WriteCompose(ctx, b)
		if err != nil {
			l.ErrorContext(ctx, hhconst.MsgFailTransactionCommit,
				hhconst.LogErr, err,
			)
			return errorResponse(ctx, http.StatusInternalServerError, hhconst.MsgFailTransactionCommit)
		}

		err = tx1.Commit(ctx)
		if err != nil {
			l.ErrorContext(ctx, "Failed to commit transaction.",
				hhconst.LogErr, err,
				"txName", "tx1",
			)
			return errorResponse(ctx, http.StatusInternalServerError, hhconst.MsgFailTransactionCommit)
		}
	}

	err = c.s.SES.Send(ctx, b)
	if err != nil {
		l.ErrorContext(ctx, "Failed to send email.",
			hhconst.LogErr, err,
		)
		return errorResponse(ctx, http.StatusInternalServerError, "Failed to send email.")
	}

	if c.s.Conf.ASWU.UsePostgres {
		tx2, err := c.s.BeginTx(ctx)
		if err != nil {
			l.ErrorContext(ctx, hhconst.MsgFailTransactionBegin,
				hhconst.LogErr, err,
				"txName", "tx2",
			)
			return errorResponse(ctx, http.StatusInternalServerError, hhconst.MsgFailTransactionBegin)
		}
		//goland:noinspection GoUnhandledErrorResult
		defer tx2.Rollback(ctx)

		ctx = context.WithValue(ctx, ctxkey.Tx, tx2)
		err = c.s.Postgres.UpdateComposeSESAccepted(ctx, b.UUID, true)
		if err != nil {
			l.ErrorContext(ctx, "Failed to update compose to database.",
				hhconst.LogErr, err,
			)
			return errorResponse(ctx, http.StatusInternalServerError, "Failed to update compose to database.")
		}
		b.SESAccepted = true

		err = tx2.Commit(ctx)
		if err != nil {
			l.ErrorContext(ctx, hhconst.MsgFailTransactionCommit,
				hhconst.LogErr, err,
				"txName", "tx2",
			)
			return errorResponse(ctx, http.StatusInternalServerError, hhconst.MsgFailTransactionCommit)
		}
	}

	return http.StatusAccepted, nil, nil
}
func (c *Compose) URLPattern() string {
	return aswu.PathAPICompose
}
