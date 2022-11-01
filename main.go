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
	Stub    serv.Stub   `json:"stub,omitempty"`
}

type failMsg struct {
	Success bool `json:"success"`
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
	var stub serv.PostStubJSONRequestBody
	// b, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	m, _ := json.Marshal(stubsMsg{Success: false})
	// 	w.Write(m)
	// }
	// err = json.Unmarshal(b, &stub)
	// if err != nil {
	// 	m, _ := json.Marshal(stubsMsg{Success: false})
	// 	w.Write(m)
	// }
	// if stub.Path == "" {
	// 	m, _ := json.Marshal(stubsMsg{Success: false})
	// 	w.Write(m)
	// }
	// m, _ := json.Marshal(stubsMsg{Success: true, Stub: stub})
	// w.Write(m)

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

}

func (s StubberApp) PostStubStop(w http.ResponseWriter, r *http.Request) {

}

func main() {
	i := StubberApp{}
	s := serv.HandlerFromMuxWithBaseURL(i, chi.NewRouter(), "/api/v1")
	http.ListenAndServe(":3000", s)
}
