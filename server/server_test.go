package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	api "github.com/alecchampaign/remotecmds/proto/commands"
	"google.golang.org/protobuf/proto"
)

func TestServer(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"marshal/unmarshal a request": testMarshalUnmarshal,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

func testMarshalUnmarshal(t *testing.T) {
	want := &api.CommandRequest{Command: "get time"}
	data, err := proto.Marshal(want)
	require.NoError(t, err)

	got := &api.CommandRequest{}
	err = proto.Unmarshal(data, got)
	require.NoError(t, err)

	require.Equal(t, want.GetCommand(), got.GetCommand())
}
