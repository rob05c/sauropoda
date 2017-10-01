package ui

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rob05c/sauropoda/dino"
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

const staticDir = "static"
const imageDir = "img"
const imageExt = "png"
const imageUIPath = "img"

const MimeHTML = "text/html"
const MimeJS = "application/javascript"
const MimeCSS = "text/css"

type SrvFile struct {
	File string
	Path string
	Mime string
}

// serveFiles returns a map of files to serve, to their content types
func serveFiles() []SrvFile {
	return []SrvFile{
		SrvFile{Path: "", File: staticDir + "/" + "index.html", Mime: MimeHTML},
		SrvFile{Path: "index.html", File: staticDir + "/" + "index.html", Mime: MimeHTML},
		SrvFile{Path: "index.js", File: staticDir + "/" + "index.js", Mime: MimeJS},
		SrvFile{Path: "index.css", File: staticDir + "/" + "index.css", Mime: MimeCSS},
		SrvFile{Path: "login", File: staticDir + "/" + "login.html", Mime: MimeHTML},
		SrvFile{Path: "login.html", File: staticDir + "/" + "login.html", Mime: MimeHTML},
		SrvFile{Path: "login.js", File: staticDir + "/" + "login.js", Mime: MimeJS},
		SrvFile{Path: "login.css", File: staticDir + "/" + "login.css", Mime: MimeCSS},
		SrvFile{Path: "player", File: staticDir + "/" + "player.html", Mime: MimeHTML},
		SrvFile{Path: "player.html", File: staticDir + "/" + "player.html", Mime: MimeHTML},
		SrvFile{Path: "player.js", File: staticDir + "/" + "player.js", Mime: MimeJS},
		SrvFile{Path: "player.css", File: staticDir + "/" + "player.css", Mime: MimeCSS},
		SrvFile{Path: "dinos", File: staticDir + "/" + "dinos.html", Mime: MimeHTML},
		SrvFile{Path: "dinos.html", File: staticDir + "/" + "dinos.html", Mime: MimeHTML},
		SrvFile{Path: "dinos.js", File: staticDir + "/" + "dinos.js", Mime: MimeJS},
		SrvFile{Path: "dinos.css", File: staticDir + "/" + "dinos.css", Mime: MimeCSS},
		SrvFile{Path: "journal", File: staticDir + "/" + "journal.html", Mime: MimeHTML},
		SrvFile{Path: "journal.html", File: staticDir + "/" + "journal.html", Mime: MimeHTML},
		SrvFile{Path: "journal.js", File: staticDir + "/" + "journal.js", Mime: MimeJS},
		SrvFile{Path: "journal.css", File: staticDir + "/" + "journal.css", Mime: MimeCSS},
	}
}

// RegisterHandlers registers the HTTP endpoints for the web UI, with the given mux.
// The pathPrefix is a path to prefix all served paths with. For example, to serve the UI at `/ui/` pass the pathPrefix `ui`.
func RegisterHandlers(mux *http.ServeMux, pathPrefix string, species map[string]dino.Species) error {
	for name, _ := range species {
		lname := strings.ToLower(name)
		filename := lname + "." + imageExt
		imgPath := staticDir + "/" + imageDir + "/species/" + filename
		handler, err := resourceHandler(imgPath, "")
		if err != nil {
			return fmt.Errorf("%s image '%s' failed to load: %v", name, imgPath, err)
		}
		mux.HandleFunc(pathPrefix+"/"+imageUIPath+"/"+lname+".png", handler)
	}
	if pathPrefix == "" {
		pathPrefix = "/"
	}
	for _, srv := range serveFiles() {
		h, err := resourceHandler(srv.File, srv.Mime)
		if err != nil {
			return fmt.Errorf("%s failed to load: %v", srv.File, err)
		}
		mux.HandleFunc(pathPrefix+srv.Path, h)
	}

	return nil
}
