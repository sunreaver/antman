package redis

import (
	"fmt"
	"time"

	r "github.com/go-redis/redis"
)

// Redis Redis
type Redis struct {
	r.UniversalClient
	cfg Config
}

// Key Key
func (r *Redis) Key(id, bundle interface{}) string {
	return fmt.Sprintf("%s%v:%v", r.cfg.Prefix, id, bundle)
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
	tmp := r.NewClient(&r.Options{
		Addr:         cfg.Hosts[0],
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Millisecond,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Millisecond,
		PoolSize:     cfg.Poolsize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
	})
	if e := tmp.Ping().Err(); e != nil {
		return nil, e
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
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Millisecond,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Millisecond,
		PoolSize:     cfg.Poolsize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
	})
	if e := tmp.Ping().Err(); e != nil {
		return nil, e
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
