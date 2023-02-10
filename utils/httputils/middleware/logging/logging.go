package logging

import (
	"applicationDesignTest/utils/httputils"
	"applicationDesignTest/utils/httputils/middleware"
	"applicationDesignTest/utils/log"
	"context"
	"net/http"
	"time"
)

func NewHandlerMiddleware(logger log.Logger) middleware.Handler {
	return func(next httputils.Handler) httputils.Handler {
		return func(ctx context.Context, req httputils.Request) (resp httputils.Response, err error) {
			logger.Infof(`"msg":"start handling request","http_method":"%s","http_url_path":"%s"`, req.Method, req.Path)

			startedAt := time.Now()
			resp, err = next(ctx, req)
			duration := time.Since(startedAt)

			if apiErr, ok := httputils.APIErrorFromError(err); ok {
				if apiErr.StatusCode >= http.StatusInternalServerError {
					logger.Errorf(
						`"msg":"%s","http_method":"%s","http_url_path":"%s","http_status_code":"%d","http_status_text":"%s","duration":"%s"`,
						req.Method, req.Path, apiErr.StatusCode, http.StatusText(apiErr.StatusCode), duration,
					)
				} else if apiErr.StatusCode >= http.StatusBadRequest {
					logger.Errorf(
						`"msg":"%s","http_method":"%s","http_url_path":"%s","http_status_code":"%d","http_status_text":"%s","duration":"%s"`,
						apiErr.Message, req.Method, req.Path, apiErr.StatusCode, http.StatusText(apiErr.StatusCode), duration,
					)
				}
			} else if err != nil {
				logger.Errorf(
					`"msg":"internal handler error","http_method":"%s","http_url_path":"%s","duration":"%s"`,
					req.Method, req.Path, duration,
				)
			} else {
				logger.Infof(
					`"msg":"request handled","http_method":"%s","http_url_path":"%s","http_status_code":"%d","http_status_text":"%s","duration":"%s"`,
					req.Method, req.Path, apiErr.StatusCode, http.StatusText(apiErr.StatusCode), duration,
				)
			}

			return
		}
	}
}
