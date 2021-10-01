package main

import (
	"fmt"
	"log"

	"github.com/alecchampaign/remotecmds/proto/commands"
	"google.golang.org/protobuf/proto"
)

func main() {
	req := &commands.CommandRequest{Command: "get time"}
	data, err := proto.Marshal(req)
	if err != nil {
		log.Fatalf("error while marshaling object : %v", err)
	}

	res := &commands.CommandRequest{}
	if err = proto.Unmarshal(data, res); err != nil {
		log.Fatalf("error while unmarshaling object: %v", err)
	}

	fmt.Printf("Data from unmarshaled object: %v", res.GetCommand())
}
