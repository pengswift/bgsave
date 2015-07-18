package main

import (
	"os"
	"testing"
)

import (
	"github.com/fzzy/radix/extra/cluster"
	"github.com/golang/snappy/snappy"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/vmihailenco/msgpack.v2"
)

import (
	pb "proto"
)

const (
	address            = "localhost:50004"
	test_key           = "testing:3721"
	DEFAULT_REDIS_HOST = "127.0.0.1:7000"
	ENV_SNAPPY         = "ENABLE_SNAPPY"
)

type TestStruct struct {
	Id    int32
	Name  string
	Sex   int
	Int64 int64
	F64   float64
	F32   float32
	Data  []byte
}

func TestBgSave(t *testing.T) {
	//连接mongodb
	client, err := cluster.NewCluster(DEFAULT_REDIS_HOST)
	if err != nil {
		t.Fatal(err)
	}

	//序列化数据
	bin, _ := msgpack.Marshal(&TestStruct{3721, "hello", 18, 999, 1.1, 2.2, []byte("world")})

	//压缩
	if env := os.Getenv(ENV_SNAPPY); env != "" {
		if enc, err := snappy.Encode(nil, bin); err == nil {
			bin = enc
		} else {
			t.Fatal(err)
		}
	}

	//存入
	reply := client.Cmd("set", test_key, bin)
	if reply.Err != nil {
		t.Fatal(err)
	}

	//调用bg save
	conn, err := grpc.Dial(address)
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBgSaveServiceClient(conn)

	//传入需要落地的key
	_, err = c.MarkDirty(context.Background(), &pb.BgSave_Key{Name: test_key})
	if err != nil {
		t.Fatalf("could not query: %v", err)
	}

}
