package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	api "github.com/alecchampaign/remotecmds/proto/commands"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ParseRequest(res http.ResponseWriter, req *http.Request) ([]string, error) {
	request := &api.CommandRequest{}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read message from request : %v", err)
	}
	if err = proto.Unmarshal(data, request); err != nil {
		return nil, fmt.Errorf("error while unmarshaling request : %d", err)
	}
	return request.GetCommands(), nil
}

type jobs struct {
	mu     sync.Mutex
	status map[string]bool
	result *api.CommandResponse
}

func (j *jobs) Complete(key string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.status[key] = true
}

func newJobs() *jobs {
	var j jobs
	j.status = make(map[string]bool)
	j.result = new(api.CommandResponse)
	return &j
}

func main() {
	jobs := newJobs()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received")
		commands, err := ParseRequest(w, r)
		if err != nil {
			log.Fatal(err)
		}

		var response []byte
		var wg sync.WaitGroup

		for _, command := range commands {
			wg.Add(1)
			switch command {
			case "get time test":
				go func() {
					defer func() {
						jobs.status[command] = true
						wg.Done()
					}()
					jobs.result.CurrTime = timestamppb.New(time.Now())
				}()
			case "say something":
				go func() {
					defer func() {
						jobs.status[command] = true
						wg.Done()
					}()
					jobs.result.Speak = "Hello world!"
				}()
			default:
				w.WriteHeader(http.StatusNotFound)
				jobs.result.Error = "Invalid command"
			}
		}

		wg.Wait()
		response, err = proto.Marshal(jobs.result)
		if err != nil {
			log.Fatalf("error while marshaling response : %d", err)
		}

		w.Write(response)
	})

	fmt.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}
