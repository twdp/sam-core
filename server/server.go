
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/rpcx-ecosystem/rpcx-examples3"
	"github.com/smallnest/rpcx/server"
	"tianwei.pro/sam-core/config"
	_ "tianwei.pro/sam-core/config"
	"tianwei.pro/sam-core/rpc"
)

func main() {
	flag.Parse()
	s := server.NewServer()
	rpc.AddRegistryPlugin(s)

	s.RegisterName("Arith", new(example.Arith), "")

	s.Serve("tcp", config.Conf.DefaultString("addr", "localhost:29998"))
}

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

type Arith int

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	fmt.Printf("call: %d * %d = %d\n", args.A, args.B, reply.C)
	return nil
}

func (t *Arith) Add(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A + args.B
	fmt.Printf("call: %d + %d = %d\n", args.A, args.B, reply.C)
	return nil
}

func (t *Arith) Say(ctx context.Context, args *string, reply *string) error {
	*reply = "hello " + *args
	return nil
}