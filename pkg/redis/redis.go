package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// Config redis配置参数
type Config struct {
	Addr      string // 地址(IP:Port)
	DB        int    // 数据库
	Password  string // 密码
	KeyPrefix string // 存储key的前缀
}

// NewStore 创建基于redis存储实例
func NewStore(cfg *Config) *Store {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	})
	return &Store{
		cli:    cli,
		prefix: cfg.KeyPrefix,
	}
}

// NewStoreWithClient 使用redis客户端创建存储实例
func NewStoreWithClient(cli *redis.Client, keyPrefix string) *Store {
	return &Store{
		cli:    cli,
		prefix: keyPrefix,
	}
}

// NewStoreWithClusterClient 使用redis集群客户端创建存储实例
func NewStoreWithClusterClient(cli *redis.ClusterClient, keyPrefix string) *Store {
	return &Store{
		cli:    cli,
		prefix: keyPrefix,
	}
}

type redisClient interface {
	redis.Cmdable
	Close() error
}

// Store redis存储
type Store struct {
	cli    redisClient
	prefix string
}

func (s *Store) wrapperKey(key string) string {
	return fmt.Sprintf("%s:%s", s.prefix, key)
}

func (s *Store) Get(key string) interface{} {
	result, err := s.cli.Get(s.wrapperKey(key)).Result()
	if err != nil {
		return ""
	}
	return result
}

func (s *Store) Set(key string, v interface{}, expiration time.Duration) error {
	cmd := s.cli.Set(s.wrapperKey(key), v, expiration)
	return cmd.Err()
}

func (s *Store) IsExist(key string) bool {
	cmd := s.cli.Exists(s.wrapperKey(key))
	return cmd.Err() == nil
}

func (s *Store) Delete(key string) error {
	cmd := s.cli.Del(s.wrapperKey(key))
	return cmd.Err()
}

func (s *Store) Check(key string) (bool, error) {
	cmd := s.cli.Exists(s.wrapperKey(key))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}
func (s *Store) ZAdd(key string, members ...redis.Z) (bool, error) {
	cmd := s.cli.ZAdd(s.wrapperKey(key), members...)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (s *Store) ZIncrBy(key string, increment float64, member string) (bool, error) {
	cmd := s.cli.ZIncrBy(s.wrapperKey(key), increment, member)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (s *Store) ZRem(key string, members ...interface{}) (bool, error) {
	cmd := s.cli.ZRem(s.wrapperKey(key), members...)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}
func (s *Store) ZRemRangeByScore(key, min, max string) (bool, error) {
	cmd := s.cli.ZRemRangeByScore(s.wrapperKey(key), min, max)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (s *Store) ZRange(key string, start, stop int64) ([]string, error) {
	return s.cli.ZRange(s.wrapperKey(key), start, stop).Result()
}

func (s *Store) ZRevRange(key string, start, stop int64) ([]string, error) {
	return s.cli.ZRevRange(s.wrapperKey(key), start, stop).Result()
}

func (s *Store) ZRangeByScore(key string, opt redis.ZRangeBy) ([]string, error) {
	return s.cli.ZRangeByScore(s.wrapperKey(key), opt).Result()
}

func (s *Store) ZCount(key, min, max string) (int64, error) {
	return s.cli.ZCount(s.wrapperKey(key), min, max).Result()
}

func (s *Store) Close() error {
	return s.cli.Close()
}
