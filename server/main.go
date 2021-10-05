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

func ParseRequest(res http.ResponseWriter, req *http.Request) (string, error) {
	request := &api.CommandRequest{}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read message from request : %v", err)
	}
	if err = proto.Unmarshal(data, request); err != nil {
		return "", fmt.Errorf("error while unmarshaling request : %d", err)
	}
	return request.GetCommand(), nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received")
		command, err := ParseRequest(w, r)
		if err != nil {
			log.Fatal(err)
		}

		switch command {
		case "get time test":
			result := &api.CurrentTimeResponse{
				CurrTime: timestamppb.New(time.Now()),
			}
			response, err := proto.Marshal(result)
			if err != nil {
				log.Fatalf("error while marshaling response : %d", err)
			}
			w.Write(response)
			// TODO: write case in event of error?
		}
	})

	fmt.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}
