package hndl

import (
	"database/sql"
	"net/http"

	sdb "github.com/rob05c/sauropoda/db"
	"github.com/rob05c/sauropoda/quadtree"
	"github.com/rob05c/sauropoda/webui"
)

func handlers(rd RouteData) map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/query/":     wrapHandler(rd, handleQuery),
		"/now":        handleNow,
		"/login":      wrapHandler(rd, handleLogin),
		"/createuser": wrapHandler(rd, handleCreateUser),
		"/ping":       wrapHandler(rd, handlePing),
		// "/catch":      wrapHandler(rd, handleCatch),
	}
}

func RegisterHandlers(rd RouteData) error {
	// TODO create api.RegisterHandlers
	uiPathPrefix := ""
	if err := webui.RegisterHandlers(http.DefaultServeMux, uiPathPrefix, rd.Species); err != nil {
		return err
	}
	handlers := handlers(rd)
	for path, handler := range handlers {
		http.HandleFunc(path, handler)
	}
	return nil
}

type RouteData struct {
	DB       *sql.DB
	Species  map[string]sdb.Species
	QT       quadtree.Quadtree
	TokenKey []byte
}

type DataHandlerFunc func(rd RouteData, w http.ResponseWriter, r *http.Request)

func wrapHandler(d RouteData, f DataHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(d, w, r)
	}
}
