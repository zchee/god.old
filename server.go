// Copyright 2017 The god Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package god

import (
	"go/build"
	"go/token"
	"net"
	"strings"
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
	grpcs  *grpc.Server
	mu     sync.RWMutex
	done   chan struct{}
	result interface{}
}

// NewServer returns the new Server.
func NewServer() *Server {
	s := grpc.NewServer()
	srv := &Server{
		grpcs: s,
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

	return s.grpcs.Serve(lis)
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
		s.grpcs.Stop()
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

func (s *Server) Output(fset *token.FileSet, qr guru.QueryResult) {
	s.mu.Lock()
	s.result = qr.Result(fset)
	s.mu.Unlock()
}

func (s *Server) query(loc *serialpb.Location) *guru.Query {
	q := &guru.Query{
		Pos:    loc.Pos,
		Build:  &build.Default,
		Output: s.Output,
	}
	if loc.Options != nil {
		// avoid corner case of split("")
		if loc.Options.Scope != "" {
			scopes := strings.Split(loc.Options.Scope, ",")
			q.Scope = scopes
		}
	}
	return q
}

func (s *Server) Ping(ctx context.Context, req *serialpb.Request) (*serialpb.Response, error) {
	return &serialpb.Response{}, nil
}

func (s *Server) GetCallees(ctx context.Context, loc *serialpb.Location) (*serialpb.Callees, error) {
	query := s.query(loc)
	if err := guru.Callees(query); err != nil {
		return nil, err
	}

	s.mu.RLock()
	res := s.result.(*serial.Callees)
	s.mu.RUnlock()

	callees := make([]*serialpb.Callee, len(res.Callees))
	for i, callee := range res.Callees {
		callees[i] = &serialpb.Callee{
			Name: callee.Name,
			Pos:  callee.Pos,
		}
	}

	return &serialpb.Callees{
		Pos:     res.Pos,
		Desc:    res.Desc,
		Callees: callees,
	}, nil
}

func (s *Server) GetCallers(ctx context.Context, loc *serialpb.Location) (*serialpb.Callers, error) {
	query := s.query(loc)
	if err := guru.Callers(query); err != nil {
		return nil, err
	}

	s.mu.RLock()
	res := s.result.([]serial.Caller)
	s.mu.RUnlock()

	callers := &serialpb.Callers{
		Callers: make([]*serialpb.Caller, len(res)),
	}
	for i, caller := range res {
		callers.Callers[i] = &serialpb.Caller{
			Pos:    caller.Pos,
			Desc:   caller.Desc,
			Caller: caller.Caller,
		}
	}

	return callers, nil
}

func (s *Server) GetCallStack(ctx context.Context, loc *serialpb.Location) (*serialpb.CallStack, error) {
	query := s.query(loc)
	if err := guru.Callstack(query); err != nil {
		return nil, err
	}

	s.mu.RLock()
	res := s.result.(*serial.CallStack)
	s.mu.RUnlock()

	callers := make([]serialpb.Caller, len(res.Callers))
	for i, caller := range res.Callers {
		callers[i] = serialpb.Caller{
			Pos:    caller.Pos,
			Desc:   caller.Desc,
			Caller: caller.Caller,
		}
	}

	return &serialpb.CallStack{
		Pos:     res.Pos,
		Target:  res.Target,
		Callers: callers,
	}, nil
}

func (s *Server) GetDefinition(ctx context.Context, loc *serialpb.Location) (*serialpb.Definition, error) {
	query := s.query(loc)
	if err := guru.Definition(query); err != nil {
		return nil, err
	}

	s.mu.RLock()
	def := s.result.(*serial.Definition)
	s.mu.RUnlock()

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
