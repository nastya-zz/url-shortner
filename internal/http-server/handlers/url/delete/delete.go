package delete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	hs "url-shoter/internal/http-server"
	resp "url-shoter/internal/lib/api/response"
	"url-shoter/internal/lib/logger/sl"
	"url-shoter/internal/storage"
)

type Request struct {
	Alias string `json:"alias,omitempty" `
}

type Response struct {
	resp.Response
}

type UrlDeleter interface {
	DeleteUrl(alias string) error
}

func New(log *slog.Logger, urlDeleter UrlDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.delete.url.New"

		log = log.With(
			slog.String("op", op),
			slog.String(hs.RequestId, middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		err := urlDeleter.DeleteUrl(alias)
		if err != nil {
			if errors.Is(err, storage.ErrAliasNotFound) {
				log.Info("alias not found", slog.String("alias", alias))
				render.JSON(w, r, resp.Error("alias not found"))
				return
			}

			log.Error("failed to delete alias", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete alias"))

			return
		}

		log.Info("delete url by alias", slog.String("alias", alias))
		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
