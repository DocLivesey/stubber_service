package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const cfg string = "./cfg.json"

var vp viper.Viper
var hosts []string
var templ *template.Template

type StubsMsg struct {
	Success bool   `json:"success"`
	Stubs   []Stub `json:"stubs,omitempty"`
}

type Stub struct {
	Jar   string `json:"jar,omitempty"`
	Path  string `json:"path"`
	State bool   `json:"state,omitempty"`
	Pid   string `json:"pid,omitempty"`
	Cpu   string `json:"cpu,omitempty"`
	Mem   string `json:"mem,omitempty"`
	Port  string `json:"port,omitempty"`
}

type contextKey string

func init() {
	vp := viper.New()
	vp.SetConfigFile(cfg)
	if err := vp.ReadInConfig(); err != nil {
		logrus.WithFields(logrus.Fields{
			"file": cfg,
		}).Fatal("error parsing cfg")
	}
	hosts = vp.GetStringSlice("slaves")
	loglevel, err := logrus.ParseLevel(vp.GetString("loglevel"))
	if err != nil {
		loglevel = logrus.ErrorLevel
	}
	logrus.SetLevel(loglevel)

}

func viewStubs(w http.ResponseWriter, r *http.Request) {
	data := make(map[string][]Stub, 4)
	var m StubsMsg
	var s []Stub
	for _, host := range hosts {
		resp, err := http.Get(host + "/api/v1/stub/all")
		if err != nil || resp.StatusCode != 200 {
			logrus.WithFields(logrus.Fields{
				"path":   resp.Request.URL,
				"Status": resp.Status,
			}).Error("fail to get all stubs")
			break
		}
		if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
			logrus.WithField("host", resp.Request.URL).Error("fail to parse response")
		}
		s = append(s, m.Stubs...)
		data[host] = s
	}
	templ.Execute(w, data)
}

func start(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if len(query) != 2 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.WithFields(logrus.Fields{
			"URL":            r.URL,
			"params_in_code": query,
		}).Error("Bad request")
		return
	}
	m := struct {
		Path string `json:"path"`
	}{
		Path: query["path"][0],
	}
	s := new(bytes.Buffer)
	json.NewEncoder(s).Encode(m)
	resp, err := http.Post(query["host"][0]+"/api/v1/stub/start", "application/json", s)
	if err != nil || resp.StatusCode != 200 {
		w.WriteHeader(http.StatusBadRequest)
		logrus.WithFields(logrus.Fields{
			"send_to":        resp.Request.URL,
			"Body":           resp.Request.Body,
			"params_in_code": query,
		}).Error("Bad request to slave")
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	// http.RedirectHandler("/", http.StatusSeeOther)
}

func stop(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if len(query) != 3 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.WithFields(logrus.Fields{
			"URL":            r.URL,
			"params_in_code": query,
		}).Error("Bad request")
		return
	}
	m := struct {
		Path string `json:"path"`
		Pid  string `json:"pid"`
	}{
		Path: query["path"][0],
		Pid:  query["pid"][0],
	}
	s := new(bytes.Buffer)
	json.NewEncoder(s).Encode(m)
	fmt.Println(s)
	resp, err := http.Post(query["host"][0]+"/api/v1/stub/stop", "application/json", s)
	if err != nil || resp.StatusCode != 200 {
		w.WriteHeader(http.StatusBadRequest)
		logrus.WithFields(logrus.Fields{
			"send_to":        resp.Request.URL,
			"Body":           resp.Request.Body,
			"params_in_code": query,
		}).Error("Bad request to slave")
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

	// http.RedirectHandler("/", http.StatusSeeOther)
}

func main() {

	templ = template.Must(template.ParseFiles("templates/index.html"))

	r := chi.NewRouter()
	r.Use(addId, logger, middleware.CleanPath, middleware.Recoverer, middleware.Timeout(10*time.Second))

	fs := http.FileServer(http.Dir("./static"))
	// r.Get("/", viewStubs)
	r.Group(func(r chi.Router) {
		r.Get("/", viewStubs)
		r.Handle("/static/*", http.StripPrefix("/static/", fs))
		r.Get("/start", start)
		r.Get("/stop", stop)
	})

	http.ListenAndServe(":3030", r)

}

func addId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id contextKey = "id"

		if r.Context().Value(id) == nil {
			uid := xid.New()
			ctx := context.WithValue(r.Context(), id, uid.String())
			cookie := http.Cookie{
				Name:     "id",
				Value:    uid.String(),
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, &cookie)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func logger(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		var id contextKey = "id"

		logrus.WithFields(logrus.Fields{
			"url":           r.URL,
			"id":            r.Context().Value(id),
			"body":          r.Body,
			"response_time": time.Since(start) * time.Millisecond,
		}).Info("request_details")

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
