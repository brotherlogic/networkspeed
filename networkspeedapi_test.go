package main

import (
	"context"
	"testing"

	pb "github.com/brotherlogic/networkspeed/proto"
)

func TestMakeTransfer(t *testing.T) {
	s := InitTest()
	_, err := s.MakeTransfer(context.Background(), &pb.TransferRequest{})
	if err != nil {
		t.Errorf("Transfer request failed: %v", err)
	}
}

func TestRecordTransfer(t *testing.T) {
	s := InitTest()
	_, err := s.RecordTransfer(context.Background(), &pb.RecordRequest{Transfer: &pb.Transfer{}})
	if err != nil {
		t.Errorf("Transfer record failed: %v", err)
	}

	if len(s.config.Transfers) != 1 {
		t.Errorf("Wrong number of transfers: %v", s.config.Transfers)
	}
}
