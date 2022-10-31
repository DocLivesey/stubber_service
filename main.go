package main

import (
	"encoding/json"
	"net/http"

	"github.com/DocLivesey/stubber_service/bash"
	serv "github.com/DocLivesey/stubber_service/gen_service"
	"github.com/go-chi/chi/v5"
)

type stubsMsg struct {
	Success bool        `json:"sucess"`
	Stubs   []serv.Stub `json:"stubs,omitempty"`
}

type StubberApp struct{}

func (s StubberApp) GetStubAll(w http.ResponseWriter, r *http.Request) {
	// j := json.NewEncoder(w)
	ss, err := bash.Populate()
	if err != nil {
		m, _ := json.Marshal(stubsMsg{Success: false})
		w.Write(m)
	}
	w.Header().Set("Content-Type", "application/json")
	m, _ := json.Marshal(stubsMsg{Success: true, Stubs: ss})
	w.Write(m)
}

func (s StubberApp) PostStub(w http.ResponseWriter, r *http.Request) {

}

func (s StubberApp) PostStubStart(w http.ResponseWriter, r *http.Request) {

}

func (s StubberApp) PostStubStop(w http.ResponseWriter, r *http.Request) {

}

func main() {
	i := StubberApp{}
	s := serv.HandlerFromMuxWithBaseURL(i, chi.NewRouter(), "/api/v1")
	http.ListenAndServe(":3000", s)
}
