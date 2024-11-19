package queue

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/danielpnjt/speed-engine/internal/config"
)

type LogData struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Context struct {
	ServiceName    string `json:"_app_name"`
	ServiceVersion string `json:"_app_version"`
	ServicePort    int    `json:"_app_port"`
	Tag            string `json:"_app_tag"`

	ReqMethod string `json:"_app_method"`
	ReqURI    string `json:"_app_uri"`

	AdditionalData map[string]interface{} `json:"_app_data,omitempty"`
}

func (w *worker) EnqueueStatusTopUp(ctx context.Context, external string) (err error) {
	ctxLogger := Context{
		ServiceName:    fmt.Sprintf("%s-worker", config.GetString("app.name")),
		ServiceVersion: config.GetString("app.version"),
		ServicePort:    config.GetInt("app.port"),
		Tag:            config.GetString("app.name"),
		ReqMethod:      "enqueue-speed_engine-topup",
		ReqURI:         "enqueue-speed_engine-topup",
	}

	ctx = context.WithValue(ctx, ctx, ctxLogger)

	slog.InfoContext(ctx, "worker start running job")

	res, err := w.transactionService.TopUp(ctx, external)
	if err != nil {
		slog.ErrorContext(ctx, "failed to execute get status ojk reporting pusdafil", err)
		err = nil
		return
	}
	slog.InfoContext(ctx, "success execute get status payment", res)
	return
}
