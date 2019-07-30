package main

import (
	"fmt"
	"testing"

	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"

	pbd "github.com/brotherlogic/discovery/proto"
	pb "github.com/brotherlogic/networkspeed/proto"
)

func InitTest() *Server {
	s := Init()
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient("./testing")
	s.bridge = &testBridge{servers: []string{}}
	s.Registry = &pbd.RegistryEntry{Identifier: "test"}

	return s
}

type testBridge struct {
	servers        []string
	failGetServers bool
}

func (t *testBridge) getServers(ctx context.Context) ([]string, error) {
	if t.failGetServers {
		return []string{}, fmt.Errorf("Built to fail")
	}
	return t.servers, nil
}
func (t *testBridge) makeTransfer(ctx context.Context, server string) (*pb.TransferResponse, error) {
	return &pb.TransferResponse{}, nil
}
func (t *testBridge) recordTransfer(ctx context.Context, trans *pb.Transfer) (*pb.RecordResponse, error) {
	return &pb.RecordResponse{}, nil
}

func TestStuff(t *testing.T) {
	s := InitTest()
	s.addTransfer(&pb.Transfer{})
}

func TestRunTransfer(t *testing.T) {
	s := InitTest()
	s.bridge = &testBridge{servers: []string{"testserver1"}}

	err := s.runTransfer(context.Background())
	if err != nil {
		t.Errorf("Transfer failed")
	}
}

func TestRunTransferFailServerGet(t *testing.T) {
	s := InitTest()
	s.bridge = &testBridge{failGetServers: true}

	err := s.runTransfer(context.Background())
	if err == nil {
		t.Errorf("Fail transfer did not fail")
	}
}

func TestRunTransferWithNoServers(t *testing.T) {
	s := InitTest()

	err := s.runTransfer(context.Background())
	if err == nil {
		t.Errorf("Transfer did not fail with no servers")
	}
}

func TestBuildPayload(t *testing.T) {
	payload := buildPayload(100)

	if len(payload) != 100 {
		t.Error("Payload is wrong length: %v", len(payload))
	}

	found := false
	for _, val := range payload {
		if val > 0 {
			found = true
		}
	}

	if !found {
		t.Errorf("Payload is empty")
	}
}
