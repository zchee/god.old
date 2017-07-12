// Copyright 2017 The god Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package god

import (
	"go/build"
	"go/token"
	"net"
	"strconv"
	"sync"

	"github.com/zchee/god/internal/guru"
	"github.com/zchee/god/internal/log"
	serialpb "github.com/zchee/god/serial"

	"golang.org/x/net/context"
	"golang.org/x/tools/cmd/guru/serial"
	"google.golang.org/grpc"
)

// Address is the god gRPC server address.
const Address = ":7154" // g: 7, o: 15, d: 4

// Server represents a god server.
type Server struct {
	s      *grpc.Server
	lis    net.Listener
	mu     sync.Mutex
	done   chan struct{}
	result interface{}
}

// NewServer returns the new Server.
func NewServer() *Server {
	s := grpc.NewServer()
	srv := &Server{
		s: s,
	}
	serialpb.RegisterGodServer(s, srv)
	return srv
}

// serve serve the god gRPC server.
func (s *Server) serve() error {
	log.Debug("serve")
	lis, err := net.Listen("tcp", Address)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.lis = lis
	s.mu.Unlock()

	return s.s.Serve(s.lis)
}

// Start starts the god gRPC server.
func (s *Server) Start() error {
	log.Debug("Start")
	errc := make(chan error, 1)
	go func() {
		errc <- s.serve()
	}()

	// wating for serve result or done
	select {
	case err := <-errc:
		log.Debug("<-errc")
		if err != nil {
			return err
		}
	case <-s.done:
		s.s.Stop()
		defer s.lis.Close()
	}

	return nil
}

// Stop sends empty struct to done chan, and stops the god gRPC server.
func (s *Server) Stop() {
	log.Debug("Done")
	s.mu.Lock()
	s.done <- struct{}{}
	s.mu.Unlock()
}

var outputMu sync.Mutex

func (s *Server) Output(fset *token.FileSet, qr guru.QueryResult) {
	outputMu.Lock()
	defer outputMu.Unlock()
	s.result = qr.Result(fset)
}

func (s *Server) Ping(ctx context.Context, req *serialpb.Request) (*serialpb.Response, error) {
	return &serialpb.Response{}, nil
}

func (s *Server) GetCallees(ctx context.Context, loc *serialpb.Location) (*serialpb.Callees, error) {
	return &serialpb.Callees{}, nil
}

func (s *Server) GetCallers(ctx context.Context, loc *serialpb.Location) (*serialpb.Callers, error) {
	return &serialpb.Callers{}, nil
}

func (s *Server) GetCallStack(ctx context.Context, loc *serialpb.Location) (*serialpb.CallStack, error) {
	return &serialpb.CallStack{}, nil
}

func (s *Server) GetDefinition(ctx context.Context, loc *serialpb.Location) (*serialpb.Definition, error) {
	log.Debug("GetDefinition")
	query := &guru.Query{
		Pos:    loc.Filename + ":#" + strconv.Itoa(int(loc.Offset)),
		Build:  &build.Default,
		Output: s.Output,
	}
	if err := guru.Definition(query); err != nil {
		return nil, err
	}
	s.mu.Lock()
	def := s.result.(*serial.Definition)
	s.mu.Unlock()
	return &serialpb.Definition{
		ObjPos: def.ObjPos,
		Desc:   def.Desc,
	}, nil
}

func (s *Server) GetDescribe(ctx context.Context, loc *serialpb.Location) (*serialpb.DescribeMethods, error) {
	return &serialpb.DescribeMethods{}, nil
}

func (s *Server) GetFreeVar(ctx context.Context, loc *serialpb.Location) (*serialpb.FreeVars, error) {
	return &serialpb.FreeVars{}, nil
}

func (s *Server) GetImplements(ctx context.Context, loc *serialpb.Location) (*serialpb.Implements, error) {
	return &serialpb.Implements{}, nil
}

func (s *Server) GetPeers(ctx context.Context, loc *serialpb.Location) (*serialpb.Peers, error) {
	return &serialpb.Peers{}, nil
}

func (s *Server) GetPointsTo(ctx context.Context, loc *serialpb.Location) (*serialpb.PointsTo, error) {
	return &serialpb.PointsTo{}, nil
}

func (s *Server) GetReferrers(ctx context.Context, loc *serialpb.Location) (*serialpb.ReferrersPackage, error) {
	return &serialpb.ReferrersPackage{}, nil
}

func (s *Server) GetWhat(ctx context.Context, loc *serialpb.Location) (*serialpb.What, error) {
	return &serialpb.What{}, nil
}

func (s *Server) GetWhichErrs(ctx context.Context, loc *serialpb.Location) (*serialpb.WhichErrs, error) {
	return &serialpb.WhichErrs{}, nil
}
