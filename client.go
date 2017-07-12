// Copyright 2017 The god Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package god

import (
	"context"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/zchee/god/internal/log"
	serialpb "github.com/zchee/god/serial"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// Client represents a god client.
type Client struct {
	conn      *grpc.ClientConn
	godclient serialpb.GodClient
}

func init() {
	// disable the annoying grpclog output.
	grpclog.SetLogger(log.New(ioutil.Discard, "", 0))
}

// NewClient sets serialpb.NewClient() and return the new pointer Client.
func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		conn:      conn,
		godclient: serialpb.NewGodClient(conn),
	}
}

// Definition return the definition information of current cursor position.
func (c *Client) Definition(ctx context.Context, pos string) {
	log.Debugln("Definition")

	p := strings.Split(pos, ":#")
	offset, err := strconv.Atoi(p[1])
	if err != nil {
		log.Fatal(err)
	}
	loc := CreateLocation(p[0], int64(offset))
	def, err := c.godclient.GetDefinition(ctx, loc)
	if err != nil {
		log.Fatalf("could not get Definition: %v", err)
	}
	log.Debugf("def: %T => %+v\n", def, def)
}

func (c *Client) Ping() (*serialpb.Response, error) {
	log.Debugln("Ping")
	return c.godclient.Ping(context.Background(), &serialpb.Request{})
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
