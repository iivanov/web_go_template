package webserver

import (
	"encoding/json/v2"
	"log/slog"
	"net/http"

	"go.uber.org/fx"

	apperrors "project_template/internal/shared/errors"
)

// Route is an http.Handler that knows the mux pattern
// under which it will be registered.
type Route interface {
	http.Handler
	Pattern() string
}

// AsRoute annotates the given constructor to state that
// it provides a route to the "routes" group.
func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

// AppHandler is a handler function that can return an error.
// Errors are handled centrally by the AppHandlerFunc wrapper.
type AppHandler func(w http.ResponseWriter, r *http.Request) error

// AppRoute is a Route that uses AppHandler instead of http.Handler.
type AppRoute interface {
	Pattern() string
	Handle(w http.ResponseWriter, r *http.Request) error
}

// appRouteAdapter wraps an AppRoute to implement Route.
type appRouteAdapter struct {
	route  AppRoute
	logger *slog.Logger
}

func (a *appRouteAdapter) Pattern() string {
	return a.route.Pattern()
}

func (a *appRouteAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := a.route.Handle(w, r); err != nil {
		handleError(w, a.logger, err)
	}
}

func handleError(w http.ResponseWriter, logger *slog.Logger, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		writeJSONError(w, appErr.Code, appErr.Message, appErr.Details)
		return
	}
	logger.Error("unhandled error", "error", err)
	writeJSONError(w, http.StatusInternalServerError, "internal server error", nil)
}

type errorResponse struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func writeJSONError(w http.ResponseWriter, code int, message string, details any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.MarshalWrite(w, errorResponse{
		Message: message,
		Details: details,
	})
}

// AsAppRoute annotates the given constructor to state that
// it provides an AppRoute to the "routes" group.
// It wraps the AppRoute in an adapter that handles errors.
func AsAppRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(AppRoute)),
		fx.ResultTags(`group:"approutes"`),
	)
}
