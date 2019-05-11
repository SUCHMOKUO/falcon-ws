package server

import (
	"fmt"
	"github.com/SUCHMOKUO/falcon-ws/util"
	"net/http"
)

// handle the request of ip position query.
func handleLocationReq(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	hostEncoded := query.Get("h")
	if hostEncoded == "" {
		http.NotFound(w, r)
		return
	}
	host, err := util.Decode(hostEncoded)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	location, err := getLocation(host)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	response(w, location)
}

func getLocation(host string) (string, error) {
	location, ok := getLocationFromDB(host)
	if ok {
		return location, nil
	}
	location, err := getLocationFromAPI(host)
	if err != nil {
		return "", err
	}
	go setLocationToDB(host, location)
	return location, nil
}

func response(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(200)
	fmt.Fprintln(w, msg)
}
