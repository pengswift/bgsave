package main

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

import (
	"github.com/fzzy/radix/extra/cluster"
	"github.com/golang/snappy/snappy"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/vmihailenco/msgpack.v2"
)

import (
	log "github.com/pengswift/gamelibs/nsq-logger"
)

import (
	pb "proto"
)

const (
	SERVICE = "[BGSAVE]"
)

const (
	SAVE_DELAY          = 100 * time.Millisecond
	COUNT_DELAY         = 1 * time.Minute
	DEFAULT_REDIS_HOST  = "127.0.0.1:7000"
	DEFAULT_MONGODB_URL = "mongodb://127.0.0.1/mydb"
	ENV_REDIS_HOST      = "REDIS_HOST"
	ENV_MONGODB_URL     = "MONGDB_URL"
	ENV_SNAPPY          = "ENABLE_SNAPPY"
	BUFSIZ              = 65536
)

type server struct {
	wait          chan string
	redis_client  *cluster.Cluster
	db            *mgo.Database
	enable_snappy bool
}

func (s *server) init() {
	// snappy
	if env := os.Getenv(ENV_SNAPPY); env != "" {
		s.enable_snappy = true
	}

	// read redis host
	redis_host := DEFAULT_REDIS_HOST
	if env := os.Getenv(ENV_REDIS_HOST); env != "" {
		redis_host = env
	}

	// start connection to redis cluster
	client, err := cluster.NewCluster(redis_host)
	if err != nil {
		log.Critical(err)
		os.Exit(-1)
	}

	//存放redis client 句柄
	s.redis_client = client

	// read mongodb host
	mongodb_url := DEFAULT_MONGODB_URL
	if env := os.Getenv(ENV_MONGODB_URL); env != "" {
		mongodb_url = env
	}

	// start connection to mongodb
	sess, err := mgo.Dial(mongodb_url)
	if err != nil {
		log.Critical(err)
		os.Exit(-1)
	}

	s.db = sess.DB("")

	// wait chan
	s.wait = make(chan string, BUFSIZ)
	go s.loader_task()
}

func (s *server) MarkDirty(ctx context.Context, in *pb.BgSave_Key) (*pb.BgSave_NullResult, error) {
	s.wait <- in.Name
	return &pb.BgSave_NullResult{}, nil
}

func (s *server) MarkDirties(ctx context.Context, in *pb.BgSave_Keys) (*pb.BgSave_NullResult, error) {
	for k := range in.Names {
		s.wait <- in.Names[k]
	}
	return &pb.BgSave_NullResult{}, nil
}

// background loader, copy chan into map, execute dump every SAVE_DELAY
func (s *server) loader_task() {
	dirty := make(map[string]bool)
	timer := time.After(SAVE_DELAY)
	timer_count := time.After(COUNT_DELAY)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)

	//记录处理的个数
	var count uint64

	for {
		select {
		case key := <-s.wait:
			dirty[key] = true
		case <-timer:
			if len(dirty) > 0 {
				//计数
				count += uint64(len(dirty))
				//写入
				s.dump(dirty)
				//清空
				dirty = make(map[string]bool)
			}
			//重新计时
			timer = time.After(SAVE_DELAY)
		case <-timer_count:
			//打印计数
			log.Info("num records saved:", count)
			timer_count = time.After(COUNT_DELAY)
		case <-sig:
			//当有关闭信号发送时
			if len(dirty) > 0 {
				s.dump(dirty)
			}
			log.Info("SIGTERM")
			os.Exit(0)
		}
	}
}

func (s *server) dump(dirty map[string]bool) {
	//遍历dirty, 取出key
	for k := range dirty {
		raw, err := s.redis_client.Cmd("GET", k).Bytes()
		if err != nil {
			log.Critical(err)
			continue
		}

		// 解压缩
		if s.enable_snappy {
			if dec, err := snappy.Decode(nil, raw); err == nil {
				raw = dec
			} else {
				log.Critical(err)
				continue
			}
		}

		var record map[string]interface{}
		//反序列化
		err = msgpack.Unmarshal(raw, &record)
		if err != nil {
			log.Critical(err)
			continue
		}

		//分隔拿到key和value
		strs := strings.Split(k, ":")
		if len(strs) != 2 {
			log.Critical("canot split key", k)
			continue
		}

		tblname, id_str := strs[0], strs[1]

		id, err := strconv.Atoi(id_str)
		if err != nil {
			log.Critical(err)
			continue
		}

		// save data to mongodb
		_, err = s.db.C(tblname).Upsert(bson.M{"Id": id}, record)
		if err != nil {
			log.Critical(err)
			continue
		}
	}
}
