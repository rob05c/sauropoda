package hndl

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/rob05c/sauropoda/dino"
	"github.com/rob05c/sauropoda/quadtree"
	"github.com/rob05c/sauropoda/webui"
)

const apiPrefix = "/api/"

func handlers(rd RouteData) map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		apiPrefix + "query":      withrd(rd, hndlQuery),
		apiPrefix + "now":        handleNow,
		apiPrefix + "login":      withrd(rd, handleLogin),
		apiPrefix + "createuser": withrd(rd, handleCreateUser),
		apiPrefix + "ping":       withrd(rd, handlePing),
		apiPrefix + "catch":      withrd(rd, hndlCatch),
		apiPrefix + "dinos":      withrd(rd, hndlDinos),
		apiPrefix + "player":     withrd(rd, hndlPlayer),
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
	Species  map[string]dino.Species
	QT       quadtree.Quadtree
	TokenKey []byte
}

type DataHandlerFunc func(rd RouteData, w http.ResponseWriter, r *http.Request)

func withrd(d RouteData, f DataHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(d, w, r)
	}
}

// GetLatLon gets the latitude and longitude from the given request. It assumes the query string keys 'lat' and 'lon'.
func GetLatLon(r *http.Request) (float64, float64, error) {
	query := r.URL.Query()
	lats, ok := query["lat"]
	if !ok {
		return 0, 0, errors.New("no 'lat' query string value")
	}
	if len(lats) != 1 {
		return 0, 0, errors.New("multiple 'lat' query string value")
	}
	latStr := lats[0]

	lons, ok := query["lon"]
	if !ok {
		return 0, 0, errors.New("no 'lon' query string value")
	}
	if len(lons) != 1 {
		return 0, 0, errors.New("multiple 'lon' query string value")
	}
	lonStr := lons[0]

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return 0, 0, errors.New("latitude not a number")
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return 0, 0, errors.New("longitude not a number")
	}
	if lat > 90.0 || lat < -90.0 {
		return 0, 0, errors.New("latitude not between 90 and -90")
	}
	if lon > 180.0 || lat < -180.0 {
		return 0, 0, errors.New("longitude not between 180 and -180")
	}
	return lat, lon, nil
}
