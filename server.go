// Copyright 2017 The god Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package god

import (
	"github.com/zchee/god/internal/log"
	serialpb "github.com/zchee/god/serial"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Server represents a god server.
type Server struct {
	s *grpc.Server
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

func (s *Server) GetCallees(ctx context.Context, loc *serialpb.Location) (*serialpb.Callees, error) {
	return &serialpb.Callees{}, nil
}

func (s *Server) GetCallers(ctx context.Context, loc *serialpb.Location) (*serialpb.Callers, error) {
	return &serialpb.Callers{}, nil
}

func (s *Server) GetCallStack(ctx context.Context, loc *serialpb.Location) (*serialpb.CallStack, error) {
	log.Debug("GetCallStack")
	return &serialpb.CallStack{}, nil
}

func (s *Server) GetDefinition(ctx context.Context, loc *serialpb.Location) (*serialpb.Definition, error) {
	log.Debug("GetDefinition")
	return &serialpb.Definition{}, nil
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
