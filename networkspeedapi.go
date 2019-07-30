package main

import "golang.org/x/net/context"

import pb "github.com/brotherlogic/networkspeed/proto"

//MakeTransfer responds to a given transfer request
func (s *Server) MakeTransfer(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	return &pb.TransferResponse{MessageSize: int64(len(req.Payload))}, nil
}

//RecordTransfer records a transfer
func (s *Server) RecordTransfer(ctx context.Context, req *pb.RecordRequest) (*pb.RecordResponse, error) {
	s.addTransfer(req.Transfer)
	return &pb.RecordResponse{}, nil
}
