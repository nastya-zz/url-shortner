package redirect

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

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.get.url.New"

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

		resURL, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlNotFound) {
				log.Info("url not found", "alias", alias)

				render.JSON(w, r, resp.Error("not found"))

				return
			}

			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get url"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		// redirect to found url
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
