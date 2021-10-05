package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	api "github.com/alecchampaign/remotecmds/proto/commands"
	"github.com/golang/protobuf/proto"
)

func makeCurrTimeRequest(req *api.CommandRequest) (*api.CurrentTimeResponse, error) {
	request, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://0.0.0.0:3000", "application/x-binary", bytes.NewReader(request))
	if err != nil {
		return nil, fmt.Errorf("Error making post request : %d", err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respObj := &api.CurrentTimeResponse{}
	err = proto.Unmarshal(respBytes, respObj)
	if err != nil {
		return nil, err
	}

	return respObj, nil
}

func main() {
	req := &api.CommandRequest{
		Command: "get time test",
	}
	res, err := makeCurrTimeRequest(req)
	if err != nil {
		log.Fatalf("error making current time request : %d", err)
	}
	fmt.Println(res.CurrTime)
}
