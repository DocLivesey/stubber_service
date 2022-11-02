package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DocLivesey/stubber_service/bash"
	serv "github.com/DocLivesey/stubber_service/gen_service"
	"github.com/go-chi/chi/v5"
)

type stubsMsg struct {
	Success bool        `json:"sucess"`
	Stubs   []serv.Stub `json:"stubs,omitempty"`
	Stub    serv.Stub   `json:"stub,omitempty"`
}

type failMsg struct {
	Success bool `json:"success"`
}

type StubberApp struct{}

func (s StubberApp) GetStubAll(w http.ResponseWriter, r *http.Request) {
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
	var stub serv.PostStubJSONRequestBody

	j := json.NewEncoder(w)
	m := failMsg{Success: false}
	if err := json.NewDecoder(r.Body).Decode(&stub); err != nil {
		j.Encode(m)
		return
	}
	if stub.Path == "" {
		j.Encode(m)
		return
	}
	bash.StubStatus(&stub)
	j.Encode(stubsMsg{Success: true, Stub: stub})

}

func (s StubberApp) PostStubStart(w http.ResponseWriter, r *http.Request) {
	var stub serv.PostStubStartJSONRequestBody

	j := json.NewEncoder(w)
	m := failMsg{Success: false}
	if err := json.NewDecoder(r.Body).Decode(&stub); err != nil {
		j.Encode(m)
		return
	}
	if stub.Path == "" {
		j.Encode(m)
		return
	}
	if err := bash.StartStub(stub); err != nil {
		j.Encode(m)
		return
	}
	time.Sleep(time.Millisecond * 100)
	bash.StubStatus(&stub)
	j.Encode(stubsMsg{Success: true, Stub: stub})
}

func (s StubberApp) PostStubStop(w http.ResponseWriter, r *http.Request) {
	var stub serv.PostStubStartJSONRequestBody

	j := json.NewEncoder(w)
	m := failMsg{Success: false}
	if err := json.NewDecoder(r.Body).Decode(&stub); err != nil {
		j.Encode(m)
		return
	}
	if stub.Path == "" {
		j.Encode(m)
		return
	}
	if err := bash.StopStub(stub); err != nil {
		j.Encode(m)
		return
	}
	time.Sleep(time.Millisecond * 100)
	bash.StubStatus(&stub)
	j.Encode(stubsMsg{Success: true, Stub: stub})
}

func main() {
	i := StubberApp{}
	s := serv.HandlerFromMuxWithBaseURL(i, chi.NewRouter(), "/api/v1")
	http.ListenAndServe(":3000", s)
}
