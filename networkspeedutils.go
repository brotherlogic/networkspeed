package main

import pb "github.com/brotherlogic/networkspeed/proto"

func (s *Server) addTransfer(t *pb.Transfer) {
	s.config.Transfers = append(s.config.Transfers, t)
}
