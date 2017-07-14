// Copyright 2017 The god Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"go/build"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/zchee/god"
	"github.com/zchee/god/internal/log"

	"golang.org/x/tools/go/buildutil"
	"google.golang.org/grpc"
)

var (
	daemonize  = flag.Bool("d", false, "run god daemon instead of client")
	cpuprofile = flag.String("cpuprofile", "", "write CPU profile to file")
	scope      = flag.String("scope", "", "comma-separated list of packages the analysis should be limited to")
)

func init() {
	flag.Var((*buildutil.TagsFlag)(&build.Default.BuildTags), "tags", buildutil.TagsFlagDoc)
	flag.Parse()
}

func main() {
	if *daemonize {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

		if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
		srv := god.NewServer()
		errc := make(chan error, 1)
		go func() { errc <- srv.Start() }()
		select {
		case err := <-errc:
			log.Fatal(err)
		case <-sigc:
			return
		}
	}

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(2)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := grpc.DialContext(ctx, god.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := god.NewClient(conn)
	defer c.Close()

	if _, err := net.Dial("tcp", god.Address); err != nil {
		err := runServer()
		if err != nil {
			log.Fatal(err)
		}

		for {
			if _, err := c.Ping(); err != nil {
				continue
			}
			break
		}
	}

	opt := new(god.ClientOptions)
	if *scope != "" {
		opt.Scope = *scope
	}

	cmd := args[0]
	switch cmd {
	case "callees":
		c.Callees(ctx, args[1], opt)
	case "callers":
		c.Callers(ctx, args[1], opt)
	case "definition":
		c.Definition(ctx, args[1], opt)
	case "stop":
		c.Stop()
	default:
		log.Fatalf("unknown subcommand: %s", cmd)
	}
}

func runServer() error {
	log.Debug("runServer")
	path, err := os.Executable()
	if err != nil {
		return err
	}
	args := []string{path, "-d"}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	stdin, err := os.Open(os.DevNull)
	if err != nil {
		return err
	}
	stdout, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	stderr, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	procAttr := os.ProcAttr{
		Dir:   cwd,
		Env:   syscall.Environ(),
		Files: []*os.File{stdin, stdout, stderr},
	}
	proc, err := os.StartProcess(path, args, &procAttr)
	if err != nil {
		return err
	}

	return proc.Release()
}
