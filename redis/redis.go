package redis

import (
	"fmt"
	"time"

	r "github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// Redis Redis
type Redis struct {
	r.UniversalClient
	cfg Config
}

// Key Key
func (r *Redis) Key(id interface{}) string {
	return fmt.Sprintf("%s%v", r.cfg.Prefix, id)
}

// Client Client
var Client *Redis

// InitClient InitClient
func InitClient(cfg Config) error {
	tmp, e := MakeClient(cfg)
	if e != nil {
		return e
	}
	Client = tmp
	return nil
}

// InitCluster InitCluster
func InitCluster(cfg Config) error {
	tmp, e := MakeCluster(cfg)
	if e != nil {
		return e
	}
	Client = tmp
	return nil
}

// MakeClient MakeClient
func MakeClient(cfg Config) (*Redis, error) {
	hosts := map[string]string{}
	for _, v := range cfg.Hosts {
		hosts[v] = v
	}
	if len(hosts) == 0 {
		return nil, errors.New("no avliable hosts")
	}
	tmp := r.NewRing(&r.RingOptions{
		Addrs:        hosts,
		Password:     cfg.Password,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Millisecond,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Millisecond,
		PoolSize:     cfg.Poolsize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
		DB:           cfg.DB,
	})
	if e := tmp.Ping().Err(); e != nil {
		return nil, errors.Wrap(e, "ping")
	}
	return &Redis{
		UniversalClient: tmp,
		cfg:             cfg,
	}, nil
}

// MakeCluster MakeCluster
func MakeCluster(cfg Config) (*Redis, error) {
	tmp := r.NewClusterClient(&r.ClusterOptions{
		Addrs:        cfg.Hosts,
		Password:     cfg.Password,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Millisecond,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Millisecond,
		PoolSize:     cfg.Poolsize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
	})
	if e := tmp.Ping().Err(); e != nil {
		return nil, errors.Wrap(e, "ping")
	}
	return &Redis{
		UniversalClient: tmp,
		cfg:             cfg,
	}, nil
}

// Stop 停止redis连接
func Stop() error {
	return Client.Close()
}
