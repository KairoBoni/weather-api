package v1

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
)

type GetWeatherResponse struct{}

type GetWeatherUseCase interface {
	Get(ctx context.Context) error
}

// GenerateReportHandler godoc
func GenerateReportHandler(uc GetWeatherUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const spanName = "v1.GenerateReportHandler"

		ctx := r.Context()

		err := uc.Get(ctx)
		if err != nil {
			unknownErrorResponse(w, r)
			return
		}
		render.JSON(w, r, GetWeatherResponse{})
	}
}
