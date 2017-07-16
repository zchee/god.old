// Copyright 2017 The god Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package god

import (
	"context"
	"io/ioutil"

	"github.com/zchee/god/internal/log"
	serialpb "github.com/zchee/god/serial"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// Client represents a god client.
type Client struct {
	conn  *grpc.ClientConn
	grpcc serialpb.GodClient
}

type ClientOptions struct {
	Scope string
}

func init() {
	// disable the annoying grpclog output.
	grpclog.SetLogger(log.New(ioutil.Discard, "", 0))
}

// NewClient sets serialpb.NewClient() and return the new pointer Client.
func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		conn:  conn,
		grpcc: serialpb.NewGodClient(conn),
	}
}

// Callees return the callees information of current cursor position.
func (c *Client) Callees(ctx context.Context, pos string, opt *ClientOptions) {
	loc := &serialpb.Location{Pos: pos}
	if opt != nil {
		loc.Options = &serialpb.Options{
			Scope: opt.Scope,
		}
	}
	callees, err := c.grpcc.GetCallees(ctx, loc)
	if err != nil {
		log.Fatalf("could not get Callees: %v", err)
	}
	log.Debugf("callees: %T => %+v\n", callees, callees)
}

// Callers return the callers information of current cursor position.
func (c *Client) Callers(ctx context.Context, pos string, opt *ClientOptions) {
	loc := &serialpb.Location{Pos: pos}
	if opt != nil {
		loc.Options = &serialpb.Options{
			Scope: opt.Scope,
		}
	}
	callers, err := c.grpcc.GetCallers(ctx, loc)
	if err != nil {
		log.Fatalf("could not get Callees: %v", err)
	}
	log.Debugf("callees: %T => %+v\n", callers, callers)
}

// Callstack return the callers information of current cursor position.
func (c *Client) Callstack(ctx context.Context, pos string, opt *ClientOptions) {
	loc := &serialpb.Location{Pos: pos}
	if opt != nil {
		loc.Options = &serialpb.Options{
			Scope: opt.Scope,
		}
	}
	callstack, err := c.grpcc.GetCallStack(ctx, loc)
	if err != nil {
		log.Fatalf("could not get CallStack: %v", err)
	}
	log.Debugf("callstack: %T => %+v\n", callstack, callstack)
}

// Definition return the definition information of current cursor position.
func (c *Client) Definition(ctx context.Context, pos string, opt *ClientOptions) {
	loc := &serialpb.Location{Pos: pos}
	if opt != nil {
		loc.Options = &serialpb.Options{
			Scope: opt.Scope,
		}
	}
	def, err := c.grpcc.GetDefinition(ctx, loc)
	if err != nil {
		log.Fatalf("could not get Definition: %v", err)
	}
	log.Debugf("def: %T => %+v\n", def, def)
}

// Describe return the describe information of current cursor position.
func (c *Client) Describe(ctx context.Context, pos string, opt *ClientOptions) {
	loc := &serialpb.Location{Pos: pos}
	if opt != nil {
		loc.Options = &serialpb.Options{
			Scope: opt.Scope,
		}
	}
	desc, err := c.grpcc.GetDescribe(ctx, loc)
	if err != nil {
		log.Fatalf("could not get Describe: %v", err)
	}
	log.Debugf("desc: %T => %+v\n", desc, desc)
}

// FreeVars return the freevars information of current cursor position.
func (c *Client) FreeVars(ctx context.Context, pos string, opt *ClientOptions) {
	loc := &serialpb.Location{Pos: pos}
	if opt != nil {
		loc.Options = &serialpb.Options{
			Scope: opt.Scope,
		}
	}
	frs, err := c.grpcc.GetFreeVars(ctx, loc)
	if err != nil {
		log.Fatalf("could not get FreeVar: %v", err)
	}
	log.Debugf("frs: %T => %+v\n", frs, frs)
}

func (c *Client) Ping() (*serialpb.Response, error) {
	log.Debugln("Ping")
	return c.grpcc.Ping(context.Background(), &serialpb.Request{})
}

// Stop sends stop signal to god gRPC server.
func (c *Client) Stop() {
	log.Debugln("Stop")
}

// Close closes the grpc ClientConn.
func (c *Client) Close() error {
	log.Debugln("Close")
	return c.conn.Close()
}
