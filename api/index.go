package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/serizawa-jp/oed-proxy/client"
	errpkg "github.com/serizawa-jp/oed-proxy/error"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	logger := httplog.NewLogger("app", httplog.Options{
		JSON: true,
	})

	mux := chi.NewRouter()
	mux.Use(
		middleware.CleanPath,
		middleware.StripSlashes,
		middleware.RequestID,
		httplog.RequestLogger(logger),
		middleware.Recoverer,
		render.SetContentType(render.ContentTypeJSON),
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: false,
			MaxAge:           300,
		}),
	)

	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {

			if _, ok := r.Context().Value(render.StatusCtxKey).(int); !ok {
				w.WriteHeader(400)
			}

			oplog := httplog.LogEntry(r.Context())
			oplog.Info().Err(err).Msg("error")

			render.DefaultResponder(w, r, render.M{"status": "error"})
			return
		}

		render.DefaultResponder(w, r, v)
	}
	mux.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("oed-proxy"))
		w.WriteHeader(http.StatusOK)
	})

	mux.Post("/*", func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())

		req := &Request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			oplog.Info().Err(err).Msg("failed to decode body request")
			render.Render(w, r, errpkg.ErrInvalidRequest(err))
			return
		}
		defer r.Body.Close()

		c := client.NewOxfordClient(http.DefaultClient)
		resp, err := c.Word(req.AppID, req.AppKey, req.Lang, req.Word)
		if err != nil {
			oplog.Info().Err(err).Msg("failed to search word")
			render.Render(w, r, errpkg.ErrInvalidRequest(err))
			return
		}

		response := &Response{resp}
		render.Render(w, r, response)
	})
	mux.ServeHTTP(w, r)
}

type Request struct {
	Word   string `json:"word"`
	Lang   string `json:"lang"`
	AppID  string `json:"app_id"`
	AppKey string `json:"app_key"`
}

type Response struct {
	*client.EntriesResponse
}

func (_ *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
