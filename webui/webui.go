package webui

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	sdb "github.com/rob05c/sauropoda/db"
)

// resourceHandler loads the given file and returns a http Handler, setting CORS to * and the MIME type to mimeType. If mimeType is the empty string, it infers the mime type (when this function is called and the file is loaded, NOT every time the returned HandlerFunc is called).
func resourceHandler(filename string, mimeType string) (http.HandlerFunc, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if mimeType == "" {
		mimeType = http.DetectContentType(b)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", mimeType)
		fmt.Fprintf(w, "%s", string(b))
	}, nil
}

// image/png
// text/html

const imageDir = "images"
const imageExt = "png"
const imageUiPath = "images"
const indexPath = "index.html"

// RegisterHandlers registers the HTTP endpoints for the web UI, with the given mux.
// The pathPrefix is a path to prefix all served paths with. For example, to serve the UI at `/ui/` pass the pathPrefix `ui`.
func RegisterHandlers(mux *http.ServeMux, pathPrefix string, species map[string]sdb.Species) error {
	for name, _ := range species {
		lname := strings.ToLower(name)
		filename := lname + "." + imageExt
		imgPath := imageDir + "/species/" + filename
		handler, err := resourceHandler(imgPath, "")
		if err != nil {
			return fmt.Errorf("%s image '%s' failed to load: %v", name, imgPath, err)
		}
		mux.HandleFunc(pathPrefix+"/images/"+lname+"png", handler)
		mux.HandleFunc(pathPrefix+"/images/"+lname+".png", handler)
	}

	indexHandler, err := resourceHandler(indexPath, "")
	if err != nil {
		return fmt.Errorf("%s failed to load: %v", indexPath, err)
	}
	if pathPrefix == "" {
		pathPrefix = "/"
	}
	mux.HandleFunc(pathPrefix, indexHandler)

	return nil
}
