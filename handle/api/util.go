package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	hhconst "github.com/MicahParks/httphandle/constant"
	"github.com/MicahParks/httphandle/middleware/ctxkey"
	"github.com/MicahParks/httphandle/model"
	jt "github.com/MicahParks/jsontype"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type EmptyResponseWithMeta struct {
	Meta model.RequestMetadata `json:"requestMetadata"`
}

func (e *EmptyResponseWithMeta) AttachMeta(meta model.RequestMetadata) {
	e.Meta = meta
}

type response interface {
	AttachMeta(meta model.RequestMetadata)
}

func commitTx(ctx context.Context) (code int, body []byte, err error) {
	tx := ctx.Value(ctxkey.Tx).(pgx.Tx)
	err = tx.Commit(ctx)
	if err != nil {
		l := ctx.Value(ctxkey.Logger).(*slog.Logger)
		l.ErrorContext(ctx, "Failed to commit transaction.",
			hhconst.LogErr, err,
		)
		return errorResponse(ctx, http.StatusInternalServerError, hhconst.RespInternalServerError)
	}
	return jsonResponseOK(ctx, &EmptyResponseWithMeta{})
}

func errorBody(ctx context.Context, code int, message string) ([]byte, error) {
	data, err := json.Marshal(model.NewError(ctx, code, message))
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal error response: %w", err)
	}
	return data, nil
}

func errorResponse(ctx context.Context, code int, message string) (int, []byte, error) {
	data, err := errorBody(ctx, code, message)
	if err != nil {
		return 0, nil, err
	}
	return code, data, nil
}

func jsonBody[ReqData jt.Config[ReqData]](r *http.Request) (reqData ReqData, ctx context.Context, code int, body []byte, err error) {
	ctx = r.Context()
	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		code, body, _ = errorResponse(ctx, http.StatusBadRequest, "Failed to read request body.")
		return reqData, ctx, code, body, err
	}

	err = json.Unmarshal(b, &reqData)
	if err != nil {
		code, body, _ = errorResponse(ctx, http.StatusUnsupportedMediaType, "Failed to JSON parse request body.")
		return reqData, ctx, code, body, err
	}

	reqData, err = reqData.DefaultsAndValidate()
	if err != nil {
		code, body, _ = errorResponse(ctx, http.StatusUnprocessableEntity, "Failed to validate request body.")
		return reqData, ctx, code, body, err
	}

	return reqData, ctx, http.StatusOK, nil, nil
}

func jsonResponseOK(ctx context.Context, r response) (int, []byte, error) {
	meta := model.RequestMetadata{
		UUID: ctx.Value(ctxkey.ReqUUID).(uuid.UUID),
	}
	r.AttachMeta(meta)
	data, err := json.Marshal(r)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to JSON marshal response: %w", err)
	}
	return http.StatusOK, data, nil
}
