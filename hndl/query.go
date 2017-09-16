package hndl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rob05c/sauropoda/dinogen"
)

// TODO change to return proper error codes
func handleQuery(d RouteData, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fmt.Printf("%v %v X-Real-IP %v X-Forwarded-For %v requested %v\n", time.Now(), r.RemoteAddr, r.Header.Get("X-Real-IP"), r.Header.Get("X-Forwarded-For"), r.URL)
	urlParts := strings.Split(r.URL.Path, "/")
	if len(urlParts) < 4 {
		fmt.Fprintf(w, "Error: Not enough parts")
		return
	}
	latStr := urlParts[2]
	lonStr := urlParts[3]
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		fmt.Fprintf(w, "Error: latitude not a number")
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		fmt.Fprintf(w, "Error: longitude not a number")
		return
	}

	if lat > 90.0 || lat < -90.0 {
		fmt.Fprintf(w, "Error: latitude not between 90 and -90")
		return
	}

	if lon > 180.0 || lat < -180.0 {
		fmt.Fprintf(w, "Error: longitude not between 180 and -180")
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
