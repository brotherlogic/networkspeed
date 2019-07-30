package main

import "golang.org/x/net/context"

import pb "github.com/brotherlogic/networkspeed/proto"

//MakeTransfer responds to a given transfer request
func (s *Server) MakeTransfer(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	return &pb.TransferResponse{}, nil
}

//RecordTransfer records a transfer
func (s *Server) RecordTransfer(ctx context.Context, req *pb.RecordRequest) (*pb.RecordResponse, error) {
	s.config.Transfers = append(s.config.Transfers, req.Transfer)
	return &pb.RecordResponse{}, nil
}
