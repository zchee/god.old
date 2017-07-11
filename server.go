// Copyright 2017 The god Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package god

import (
	"github.com/zchee/god/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	s *grpc.Server
}

func NewServer() *Server {
	s := grpc.NewServer()
	srv := &Server{
		s: s,
	}
	pb.RegisterGodServer(s, srv)
	return srv
}

func (s *Server) GetCallees(ctx context.Context, loc *pb.Location) (*pb.Callees, error) {
	return nil, nil
}

func (s *Server) GetCallers(ctx context.Context, loc *pb.Location) (*pb.Callers, error) {
	return nil, nil
}

func (s *Server) GetCallStack(ctx context.Context, loc *pb.Location) (*pb.CallStack, error) {
	return nil, nil
}

func (s *Server) GetDefinition(ctx context.Context, loc *pb.Location) (*pb.Definition, error) {
	return nil, nil
}

func (s *Server) GetDescribe(ctx context.Context, loc *pb.Location) (*pb.DescribeMethods, error) {
	return nil, nil
}

func (s *Server) GetFreeVar(ctx context.Context, loc *pb.Location) (*pb.FreeVars, error) {
	return nil, nil
}

func (s *Server) GetImplements(ctx context.Context, loc *pb.Location) (*pb.Implements, error) {
	return nil, nil
}

func (s *Server) GetPeers(ctx context.Context, loc *pb.Location) (*pb.Peers, error) {
	return nil, nil
}

func (s *Server) GetPointsTo(ctx context.Context, loc *pb.Location) (*pb.PointsTo, error) {
	return nil, nil
}

func (s *Server) GetReferrers(ctx context.Context, loc *pb.Location) (*pb.ReferrersPackage, error) {
	return nil, nil
}

func (s *Server) GetWhat(ctx context.Context, loc *pb.Location) (*pb.What, error) {
	return nil, nil
}

func (s *Server) GetWhichErrs(ctx context.Context, loc *pb.Location) (*pb.WhichErrs, error) {
	return nil, nil
}
