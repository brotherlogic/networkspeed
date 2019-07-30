package main

import (
	"testing"

	"github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/networkspeed/proto"
)

func InitTest() *Server {
	s := Init()
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient("./testing")

	return s
}

func TestStuff(t *testing.T) {
	s := InitTest()
	s.addTransfer(&pb.Transfer{})
}
