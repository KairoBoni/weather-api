package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type API struct {
	// Liveness and readiness probes
	GetWeatherHandler http.HandlerFunc
}

// Routes godoc
//
//	@title						Weather API
//	@version					1.0.0
//	@description				---------------
//	@host						---------------
//	@BasePath					/api/v1
func (a *API) Routes(router *chi.Mux) {
	router.Get("/api/v1/weather", a.GetWeatherHandler)
}

type errorResponseBody struct {
	Error string `json:"error"`
}

func errorResponse(w http.ResponseWriter, r *http.Request, code int, err error) {
	render.Status(r, code)
	render.JSON(w, r, errorResponseBody{
		Error: err.Error(),
	})
}

func unknownErrorResponse(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusInternalServerError)
	render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
}
