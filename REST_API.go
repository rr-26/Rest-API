package main
import 	"net/http"
import "encoding/json"
import "sync"
import "io/ioutil"
import "fmt"
import "strings"

type MEETINGS struct{
	Id string `json:"Id"`
	Title string `json:"Title"`
	Participants string `json:"Participants"`
	Start_Time int `json:"Start_Time"`
	End_time int `json:"End_time"`
	Creation_time_stamp int `json:"Creation_time_stamp"`
}

type meetingHandlers struct{
	sync.Mutex
	store map[string]MEETINGS
}

func (M *meetingHandlers) handle(w http.ResponseWriter, r *http.Request){
	switch r.Method {
		case "GET":
			M.get(w,r)
			return
		case "POST":
			M.post(w,r)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
	}
}

func (M *meetingHandlers) get(w http.ResponseWriter, r *http.Request){
	meetings := make([]MEETINGS, len(M.store))

	M.Lock()
	i := 0
	for _, meeting := range M.store {
		meetings[i] = meeting
		i++
	}
	M.Unlock()

	jsonBytes, err := json.Marshal(meetings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Add("content-type","application/json")
	w.Write(jsonBytes)
}

func (M *meetingHandlers) paticularHandle(w http.ResponseWriter, r *http.Request){

	aadmilog := strings.Split(r.URL.String(), "/")
	if len(aadmilog) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println(aadmilog[2])

	M.Lock()
	meeting, present := M.store[aadmilog[2]]
	M.Unlock()
	if !present {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(meeting)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Add("content-type","application/json")
	w.Write(jsonBytes)
}

// func (M *meetingHandlers) paticularTime(w http.ResponseWriter, r *http.Request){

// 	aadmilog := strings.Split(r.URL.String(), "/")
// 	if len(aadmilog) != 4 {
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	fmt.Println(aadmilog[3])

// 	M.Lock()
// 	meeting, present := M.store[aadmilog[3]]
// 	M.Unlock()
// 	if !present {
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	jsonBytes, err := json.Marshal(meeting)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
// 	w.Header().Add("content-type","application/json")
// 	w.Write(jsonBytes)
// }

func (M *meetingHandlers) post(w http.ResponseWriter, r *http.Request){
	postBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	content_type := r.Header.Get("content-type")
	if content_type != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("status 415, unsupported content type, want 'application/json', got '%s'", content_type)))
		return
	}


	var newmeet MEETINGS
	err = json.Unmarshal(postBytes, &newmeet)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	M.Lock()
	M.store[newmeet.Id] = newmeet
	defer M.Unlock()
}

func newmeetingHandlers() *meetingHandlers { // this is a dummy input for non-null DB
	return &meetingHandlers{
		store: map[string]MEETINGS{
			"id1": MEETINGS{
				Id: 					"id1",
				Title: 					"First meeting",
				Participants: 			"Mukesh",
				Start_Time: 			6000,
				End_time: 				7000,
				Creation_time_stamp: 	1000,
			},
		},
	}
}

// type participant struct{
// 	Name string `json:"Name"`
// 	Email mail	`json:"Email"`
// 	RSVP string `json:"RSVP"`
// }

// func newParticipant(){
// }


func main(){
	// participant := newParticipant()
	meetingHandlers := newmeetingHandlers()
	http.HandleFunc("/meetings",meetingHandlers.handle)
	http.HandleFunc("/meetings/",meetingHandlers.paticularHandle)
	// http.HandleFunc("/meetings/",meetingHandlers.paticularTime)
	err := http.ListenAndServe(":8080",nil)
	if err != nil {
		panic(err)
	}
}