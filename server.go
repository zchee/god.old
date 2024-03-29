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

func (s *Server) GetDescribe(ctx context.Context, loc *serialpb.Location) (*serialpb.Describe, error) {
	query := s.query(loc)
	if err := guru.Describe(query); err != nil {
		return nil, err
	}

	s.mu.RLock()
	res := s.result.(*serial.Describe)
	s.mu.RUnlock()

	desc := &serialpb.Describe{
		Desc:   res.Desc,
		Pos:    res.Pos,
		Detail: res.Detail,
	}

	if res.Package != nil {
		members := make([]*serialpb.DescribeMember, len(res.Package.Members))
		for i, member := range res.Package.Members {
			members[i] = &serialpb.DescribeMember{
				Name:  member.Name,
				Type:  member.Type,
				Value: member.Value,
				Pos:   member.Pos,
				Kind:  member.Kind,
			}
		}
		desc.Package = &serialpb.DescribePackage{
			Path:    res.Package.Path,
			Members: members,
		}
	}
	if res.Type != nil {
		typ := &serialpb.DescribeType{
			Type:    res.Type.Type,
			NamePos: res.Type.NamePos,
			NameDef: res.Type.NameDef,
		}
		methods := make([]serialpb.DescribeMethod, len(res.Type.Methods))
		for i, method := range res.Type.Methods {
			methods[i] = serialpb.DescribeMethod{
				Name: method.Name,
				Pos:  method.Pos,
			}
		}
		typ.Methods = methods
		desc.Type = typ
	}
	if res.Value != nil {
		value := &serialpb.DescribeValue{
			Type:   res.Value.Type,
			Value:  res.Value.Value,
			ObjPos: res.Value.ObjPos,
		}
		desc.Value = value
	}

	return desc, nil
}

func (s *Server) GetFreeVars(ctx context.Context, loc *serialpb.Location) (*serialpb.FreeVars, error) {
	query := s.query(loc)
	if err := guru.Freevars(query); err != nil {
		return nil, err
	}

	s.mu.RLock()
	res := s.result.([]serial.FreeVar)
	s.mu.RUnlock()

	frs := &serialpb.FreeVars{
		FreeVar: make([]serialpb.FreeVar, len(res)),
	}
	for i, ref := range res {
		frs.FreeVar[i] = serialpb.FreeVar{
			Pos:  ref.Pos,
			Kind: ref.Kind,
			Ref:  ref.Ref,
			Type: ref.Type,
		}
	}

	return frs, nil
}

func (s *Server) GetImplements(ctx context.Context, loc *serialpb.Location) (*serialpb.Implements, error) {
	query := s.query(loc)
	if err := guru.Implements(query); err != nil {
		return nil, err
	}

	s.mu.RLock()
	res := s.result.(*serial.Implements)
	s.mu.RUnlock()

	impl := &serialpb.Implements{
		T: serialpb.ImplementsType{
			Name: res.T.Name,
			Pos:  res.T.Pos,
			Kind: res.T.Kind,
		},
		AssignableTo:            make([]serialpb.ImplementsType, len(res.AssignableTo)),
		AssignableFrom:          make([]serialpb.ImplementsType, len(res.AssignableFrom)),
		AssignableFromPtr:       make([]serialpb.ImplementsType, len(res.AssignableFromPtr)),
		AssignableToMethod:      make([]serialpb.DescribeMethod, len(res.AssignableToMethod)),
		AssignableFromMethod:    make([]serialpb.DescribeMethod, len(res.AssignableFromMethod)),
		AssignableFromPtrMethod: make([]serialpb.DescribeMethod, len(res.AssignableFromPtrMethod)),
	}
	if res.Method != nil {
		impl.Method = &serialpb.DescribeMethod{
			Name: res.Method.Name,
			Pos:  res.Method.Pos,
		}
	}

	for i, implType := range res.AssignableTo {
		impl.AssignableTo[i] = serialpb.ImplementsType{
			Name: implType.Name,
			Pos:  implType.Pos,
			Kind: implType.Kind,
		}
	}
	for i, implType := range res.AssignableFrom {
		impl.AssignableFrom[i] = serialpb.ImplementsType{
			Name: implType.Name,
			Pos:  implType.Pos,
			Kind: implType.Kind,
		}
	}
	for i, implType := range res.AssignableFromPtr {
		impl.AssignableFromPtr[i] = serialpb.ImplementsType{
			Name: implType.Name,
			Pos:  implType.Pos,
			Kind: implType.Kind,
		}
	}
	for i, descMethod := range res.AssignableToMethod {
		impl.AssignableToMethod[i] = serialpb.DescribeMethod{
			Name: descMethod.Name,
			Pos:  descMethod.Pos,
		}
	}
	for i, descMethod := range res.AssignableFromMethod {
		impl.AssignableFromMethod[i] = serialpb.DescribeMethod{
			Name: descMethod.Name,
			Pos:  descMethod.Pos,
		}
	}
	for i, descMethod := range res.AssignableFromPtrMethod {
		impl.AssignableFromPtrMethod[i] = serialpb.DescribeMethod{
			Name: descMethod.Name,
			Pos:  descMethod.Pos,
		}
	}

	return impl, nil
}

func (s *Server) GetPeers(ctx context.Context, loc *serialpb.Location) (*serialpb.Peers, error) {
	query := s.query(loc)
	if err := guru.Implements(query); err != nil {
		return nil, err
	}

	s.mu.RLock()
	res := s.result.(*serial.Peers)
	s.mu.RUnlock()

	peers := &serialpb.Peers{
		Pos:      res.Pos,
		Type:     res.Type,
		Allocs:   res.Allocs,
		Sends:    res.Sends,
		Receives: res.Receives,
		Closes:   res.Closes,
	}

	return peers, nil
}

func (s *Server) GetPointsTo(ctx context.Context, loc *serialpb.Location) (*serialpb.PointsTos, error) {
	query := s.query(loc)
	if err := guru.Implements(query); err != nil {
		return nil, err
	}

	s.mu.RLock()
	res := s.result.([]serial.PointsTo)
	s.mu.RUnlock()

	pts := &serialpb.PointsTos{
		PointsTos: make([]serialpb.PointsTo, len(res)),
	}
	for i, ptr := range res {
		pts.PointsTos[i] = serialpb.PointsTo{
			Type:    ptr.Type,
			NamePos: ptr.NamePos,
			Labels:  make([]serialpb.PointsToLabel, len(ptr.Labels)),
		}
		for j, label := range ptr.Labels {
			pts.PointsTos[i].Labels[j] = serialpb.PointsToLabel{
				Pos:  label.Pos,
				Desc: label.Desc,
			}
		}
	}
	return pts, nil
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
