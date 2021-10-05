package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received")
		commands, err := ParseRequest(w, r)
		if err != nil {
			log.Fatal(err)
		}

		result := &api.CommandResponse{}
		var response []byte
		for _, command := range commands {
			// For concurrency, maybe create a channel here and lanch this as a goroutine?
			switch command {
			case "get time test":
				result.CurrTime = timestamppb.New(time.Now())
			case "say something":
				result.Speak = "Hello world!"
			default:
				w.WriteHeader(http.StatusNotFound)
				result.Error = "Invalid command"
			}
		}

		response, err = proto.Marshal(result)
		if err != nil {
			log.Fatalf("error while marshaling response : %d", err)
		}

		w.Write(response)
	})

	fmt.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}
