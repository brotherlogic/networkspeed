package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbg "github.com/brotherlogic/goserver/proto"
	"github.com/brotherlogic/goserver/utils"
	pb "github.com/brotherlogic/networkspeed/proto"
)

const (
	// KEY - where we store sale info
	KEY = "/github.com/brotherlogic/networkspeedr/config"

	// PayloadInBytes the size of the payload in bytes
	PayloadInBytes = 1024 * 1024 // 1Mb
)

type bridge interface {
	getServers(ctx context.Context) ([]string, error)
	makeTransfer(ctx context.Context, server string) (*pb.TransferResponse, error)
	recordTransfer(ctx context.Context, trans *pb.Transfer) (*pb.RecordResponse, error)
}

type prodBridge struct {
	dial       func(server string) (*grpc.ClientConn, error)
	dialServer func(job, server string) (*grpc.ClientConn, error)
}

func (p *prodBridge) getServers(ctx context.Context) ([]string, error) {
	response := []string{}
	entries, err := utils.ResolveAll("networkspeed")
	if err != nil {
		return response, err
	}

	for _, entry := range entries {
		response = append(response, entry.Identifier)
	}

	return response, nil
}

func (p *prodBridge) makeTransfer(ctx context.Context, server string) (*pb.TransferResponse, error) {
	t := time.Now()
	conn, err := p.dialServer("networkspeed", server)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewTransferServiceClient(conn)
	payload := buildPayload(PayloadInBytes)
	procTime := time.Now().Sub(t).Nanoseconds()
	resp, err := client.MakeTransfer(ctx, &pb.TransferRequest{ByteSize: PayloadInBytes, Payload: payload})
	if err == nil {
		resp.ProcessingTime = procTime
	}
	return resp, err
}

func (p *prodBridge) recordTransfer(ctx context.Context, trans *pb.Transfer) (*pb.RecordResponse, error) {
	conn, err := p.dial("networkspeed")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewTransferServiceClient(conn)
	return client.RecordTransfer(ctx, &pb.RecordRequest{Transfer: trans})
}

//Server main server type
type Server struct {
	*goserver.GoServer
	config *pb.Config
	bridge bridge
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
		config:   &pb.Config{},
	}
	s.bridge = &prodBridge{dial: s.DialMaster, dialServer: s.DialServer}
	return s
}

func (s *Server) save(ctx context.Context) {
	s.KSclient.Save(ctx, KEY, s.config)
}

func (s *Server) load(ctx context.Context) error {
	config := &pb.Config{}
	data, _, err := s.KSclient.Read(ctx, KEY, config)

	if err != nil {
		return err
	}

	s.config = data.(*pb.Config)
	return nil
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterTransferServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.save(ctx)
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	if master {
		err := s.load(ctx)
		return err
	}

	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "num_transfers", Value: int64(len(s.config.Transfers))},
	}
}

type properties struct {
}

func (s *Server) deliver(w http.ResponseWriter, r *http.Request) {
	data, err := Asset("templates/main.html")
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error: %v", err))
		return
	}
	err = s.render(string(data), properties{}, w)
	if err != nil {
		s.Log(fmt.Sprintf("Error writing: %v", err))
	}
}

func (s *Server) render(f string, props properties, w io.Writer) error {
	templ := template.New("main")
	templ, err := templ.Parse(f)
	if err != nil {
		return err
	}
	templ.Execute(w, props)
	return nil
}

func (s *Server) serveUp(port int32) {
	http.HandleFunc("/", s.deliver)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	var init = flag.Bool("init", false, "Do setup")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server
	server.RegisterServer("networkspeed", false)

	if *init {
		server.config.LastCheck = time.Now().Unix()
		ctx, cancel := utils.BuildContext("networkspeed", "networkspeed")
		defer cancel()
		server.save(ctx)
		return
	}

	go server.serveUp(server.Registry.Port - 1)

	server.RegisterRepeatingTaskNonMaster(server.runTransfer, "run_transfer", time.Minute)

	fmt.Printf("%v", server.Serve())
}
