package api

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type PostBody struct {
    URL string `json:"url"`
}

type Response struct {
    Error string `json:"error,omitempty"`
    Data any `json:"data,omitempty"`
}

func NewHandler(db map[string]string) http.Handler {
    router := chi.NewMux()
    // call middlewares
    router.Use(middleware.Recoverer)
    router.Use(middleware.RequestID)
    router.Use(middleware.Logger)

    /**
     * Set routes/endpoints
     *
     * Those functions have to be in the format:
     * func(rw http.ResponseWriter, req *http.Request) => void
    **/
    router.Post("/api/shorten", handleShorten(db))
    // router.Get("/{code}", handleGet(db))

    return router
}



// Post method that sends the full URL to be shortened
func handleShorten(db map[string]string) http.HandlerFunc {
    return func (rw http.ResponseWriter, req *http.Request) {
        var body PostBody
        // if user trying to send invalid body
        if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
            sendJSON(
                rw, Response{Error: "invalid body"},
                http.StatusUnprocessableEntity,
                )
            return 
        }
        if _, err := url.Parse(body.URL); err != nil {
            sendJSON(
                rw,
                Response{Error: "invalid url parsed"},
                http.StatusBadRequest,
            )
        }
        code := genCode()
        db[code] = body.URL
        sendJSON(
            rw,
            Response{Data: code},
            http.StatusCreated,
        )
    }
}


// Get method I guess will return the shortened url
// func handleGet(db map[string]string) http.HandlerFunc {
//     return func (rw http.ResponseWriter, req *http.Request) {
//         
//     }
// }

/**
 * Auxiliary functions
 */
const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
func genCode() string {
    const size = 8
    bytes := make([]byte, size)
    for i := range size {
        bytes[i] = characters[rand.Intn(len(characters))]
    }
    return string(bytes)
}

func sendJSON(rw http.ResponseWriter, resp Response, status int) {
    rw.Header().Set("Content-Type", "application/json")

    data, err := json.Marshal(resp)
    if err != nil {
        slog.Error("failed to marshal json data", "error", err)
        sendJSON(
            rw,
            Response{Error: "something went wrong"},
            http.StatusInternalServerError,
        )
        return
    }
    rw.WriteHeader(status)
    if _, err := rw.Write(data); err != nil {
        slog.Error("failed to write json data", "error", err)
        return
    }
}
