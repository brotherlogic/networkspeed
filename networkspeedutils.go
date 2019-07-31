package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/networkspeed/proto"
)

func (s *Server) addTransfer(t *pb.Transfer) {
	s.config.Transfers = append(s.config.Transfers, t)
}

func (s *Server) runTransfer(ctx context.Context) error {
	servers, err := s.bridge.getServers(ctx)
	if err != nil {
		return err
	}

	// Process the list randomly
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for _, i := range r.Perm(len(servers)) {
		t := time.Now()
		resp, err := s.bridge.makeTransfer(ctx, servers[i])
		if err == nil {
			s.bridge.recordTransfer(ctx, &pb.Transfer{Destination: servers[i], Origin: s.Registry.Identifier, MessageSize: resp.MessageSize, TimeInNanoseconds: time.Now().Sub(t).Nanoseconds()})
			return nil
		}
	}

	return fmt.Errorf("Unable to find any suitable servers from list of %v", len(servers))
}

func buildPayload(sizeInBytes int) []byte {
	resp := make([]byte, sizeInBytes)
	rand.Read(resp)
	return resp
}

func (s *Server) buildProps() properties {
	props := properties{Servers: []string{}, Timing: make(map[string]map[string]int64)}
	counts := make(map[string]map[string]int64)
	for _, transfer := range s.config.Transfers {
		found := false
		for _, server := range props.Servers {
			if server == transfer.Origin {
				found = true
			}
		}

		if !found {
			props.Servers = append(props.Servers, transfer.Origin)
		}
		if _, ok := props.Timing[transfer.Origin]; !ok {
			props.Timing[transfer.Origin] = make(map[string]int64)
			counts[transfer.Origin] = make(map[string]int64)
		}

		if _, ok := props.Timing[transfer.Destination]; !ok {
			props.Timing[transfer.Origin][transfer.Destination] = 0
			counts[transfer.Origin][transfer.Destination] = 0
		}

		props.Timing[transfer.Origin][transfer.Destination] += transfer.TimeInNanoseconds
		counts[transfer.Origin][transfer.Destination]++
	}

	for origin, omap := range props.Timing {
		for destination, val := range omap {
			props.Timing[origin][destination] = val / counts[origin][destination]
		}
	}

	props.Servers2 = props.Servers
	return props
}
