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

func makeCurrTimeRequest(req *api.CommandRequest) (*api.CommandResponse, error) {
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

	respObj := new(api.CommandResponse)
	if err = proto.Unmarshal(respBytes, respObj); err != nil {
		return nil, err
	}

	return respObj, nil
}

func makeStatusRequest(req *api.StatusRequest) (*api.StatusResponse, error) {
	request, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://0.0.0.0:3000/status", "application/x-binary", bytes.NewReader(request))
	if err != nil {
		return nil, fmt.Errorf("error making post request : %d", err)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respObj := new(api.StatusResponse)
	if err = proto.Unmarshal(respBytes, respObj); err != nil {
		return nil, err
	}

	return respObj, nil
}

func main() {
	{
		req := &api.CommandRequest{
			Commands: []string{
				"get time test",
				"say something",
			},
		}
		res, err := makeCurrTimeRequest(req)
		if err != nil {
			log.Fatalf("error making current time request : %d", err)
		}
		fmt.Println(res)
	}
	{
		req := &api.StatusRequest{
			Command: "say something",
		}
		res, err := makeStatusRequest(req)
		if err != nil {
			log.Fatalf("error making status request : %d", err)
		}
		fmt.Println(res)
	}
}
