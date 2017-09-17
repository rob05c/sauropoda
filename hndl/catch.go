package hndl

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rob05c/sauropoda/dinogen"
	"github.com/rob05c/sauropoda/login"
	"github.com/rob05c/sauropoda/sdb"
)

func hndlCatch(d RouteData, w http.ResponseWriter, r *http.Request) {
	fmt.Print(time.Now().Format(time.RFC3339) + " INFO: " + r.RequestURI + " hndlCatch\n")
	defer r.Body.Close()
	// TODO put in "refreshCookie" helper, for all handlers
	// TODO log errors

	cookie, err := r.Cookie(login.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		fmt.Print(time.Now().Format(time.RFC3339) + " Error: " + r.RequestURI + " hndlCatch: getting cookie: " + err.Error() + "\n")
		return
	}

	user, err := login.TokenValid(cookie.Value, d.DB, d.TokenKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		fmt.Print(time.Now().Format(time.RFC3339) + " Error: " + r.RequestURI + " hndlCatch: token invalid: " + err.Error() + "\n")
		return
	}

	idStrs, ok := r.URL.Query()["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		fmt.Print(time.Now().Format(time.RFC3339) + " Error: " + r.RequestURI + " hndlCatch: no ID\n")
		return
	}
	if len(idStrs) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		fmt.Print(time.Now().Format(time.RFC3339) + " Error: " + r.RequestURI + " hndlCatch: multiple IDs\n")
		return
	}
	idStr := idStrs[0]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		fmt.Print(time.Now().Format(time.RFC3339) + " Error: " + r.RequestURI + " hndlCatch: ID not an int: " + err.Error() + "\n")
		return
	}

	dino, ok := d.QT.GetByID(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		fmt.Print(time.Now().Format(time.RFC3339) + " Error: " + r.RequestURI + " hndlCatch: ID has no dino\n")
		return
	}

	ownedDino := dinogen.PositionedToOwned(*dino)
	if err := sdb.InsertOwnedDino(d.DB, user, ownedDino); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		fmt.Print(time.Now().Format(time.RFC3339) + " Error: " + r.RequestURI + " hndlCatch: failed to insert dino: " + err.Error() + "\n")
		return
	}

	w.WriteHeader(http.StatusOK)
}
