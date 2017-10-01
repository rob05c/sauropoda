package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rob05c/sauropoda/dinogen"
)

// TODO change to return proper error codes
func hndlQuery(d RouteData, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Printf("%v %v X-Real-IP %v X-Forwarded-For %v requested %v\n", time.Now(), r.RemoteAddr, r.Header.Get("X-Real-IP"), r.Header.Get("X-Forwarded-For"), r.URL)

	handleErr := func(code int, msg string) {
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		fmt.Print(time.Now().Format(time.RFC3339) + " Error: " + r.RequestURI + " hndlQuery: " + msg + "\n")
	}

	lat, lon, err := GetLatLon(r)
	if err != nil {
		handleErr(http.StatusBadRequest, "getting latlon: "+err.Error())
		return
	}

	dinosaurs := dinogen.Query(d.QT, d.Species, lat, lon) // []quadtree.PositionedDinosaur

	dinosaursJson, err := json.Marshal(dinosaurs)
	if err != nil {
		fmt.Printf("Error marshalling dinosaurs: %v", err)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}

	fmt.Fprintf(w, "%s", string(dinosaursJson))
}
