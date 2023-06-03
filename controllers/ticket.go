package controllers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"com.routee/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/google/btree"
	"github.com/oklog/ulid"
)

type Res struct {
	StatusCode int     `json:"status_code"`
	Message    string  `json:"message"`
	TicketData ResData `json:"ticket_data"`
}

type ResData struct {
	Id    string            `json:"uid"`
	Arena string            `json:"arena"`
	PData interfaces.Ticket `json:"person_data"`
}

func (r ResData) Less(item btree.Item) bool {
	return r.Id < item.(ResData).Id
}

var tree = btree.New(200)

func CreateTicket(w http.ResponseWriter, r *http.Request) {
	var ticket interfaces.Ticket

	headerType := r.Header.Get("Content-Type")

	if headerType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("Content Type is not application/json"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	//decoder.DisallowUnknownFields()
	err := decoder.Decode(&ticket)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fid := func() string {

		entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
		ms := ulid.Timestamp(time.Now())
		gid, _ := ulid.New(ms, entropy)
		rid := string(gid.String())

		return rid
	}

	resd := ResData{
		Id:    fid(),
		Arena: "green",
		PData: ticket,
	}

	res := Res{
		StatusCode: int(201),
		Message:    "ticket generated successfully",
		TicketData: resd,
	}

	tree.ReplaceOrInsert(resd)

	jsonRes, errm := json.Marshal(res)

	if errm != nil {
		log.Panicln("Marshal unsuccesfull")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(string(jsonRes))
	w.WriteHeader(http.StatusCreated)
}

func GetTicketData(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")

	var Response ResData

	result := tree.Get(ResData{Id: params})

	if result != nil {
		Response = result.(ResData)
		w.Header().Set("Content-Type", "application/json")
		jsonRes, err := json.Marshal(Response)
		if err != nil {
			log.Panicln("Marshal unsuccesfull")
		}

		json.NewEncoder(w).Encode(string(jsonRes))
		w.WriteHeader(http.StatusOK)

	} else {
		type custom struct {
			Msg string `json:"msg"`
		}
		res := custom{
			Msg: "id not found",
		}

		w.Header().Set("Content-Type", "application/json")
		jsonRes, err := json.Marshal(res)

		if err != nil {
			log.Panicln("Marshal unsuccesfull")
		}

		json.NewEncoder(w).Encode(string(jsonRes))
		w.WriteHeader(http.StatusNotFound)
	}

}

func CancelTicket(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")

	var Response Res

	DeleteItem := tree.Delete(ResData{Id: params})
	
	if DeleteItem != nil {
		Response = Res{
			StatusCode: 200,
			Message:    "Successfully Deleted",
		}
		w.WriteHeader(http.StatusOK)
	} else {
		Response = Res{
			StatusCode: 404,
			Message:    "id not found",
		}
		w.WriteHeader(http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	jsonRes, err := json.Marshal(Response)

	if err != nil {
		log.Println("marshal unsucces:", err)
	}

	json.NewEncoder(w).Encode(string(jsonRes))
}
