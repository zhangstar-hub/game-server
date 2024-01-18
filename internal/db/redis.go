package db

import (
	"context"
	"my_app/internal/config"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

// 定义一个redis连接池
var RedisClient *Redis

func InitRedis() {
	RedisClient = NewRedis()
}

func NewRedis() *Redis {
	conf := config.GetC()
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Env.Redis.Address,
		DB:           conf.Env.Redis.DB,
		PoolSize:     10,               // 连接池中连接的最大数量
		MinIdleConns: 5,                // 最小空闲连接数
		IdleTimeout:  30 * time.Second, // 空闲连接超时时间
	})
	return &Redis{client: client}
}

func (r *Redis) Set(key, value string, expiration time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Get(key string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, key).Result()
}

// ListPush 将元素推入列表
func (r *Redis) ListPush(key, value string) error {
	ctx := context.Background()
	return r.client.LPush(ctx, key, value).Err()
}

// ListPop 从列表弹出元素
func (r *Redis) ListPop(key string) (string, error) {
	ctx := context.Background()
	return r.client.LPop(ctx, key).Result()
}

// ListRange 获取列表中指定范围的元素
func (r *Redis) ListRange(key string, start, stop int64) ([]string, error) {
	ctx := context.Background()
	return r.client.LRange(ctx, key, start, stop).Result()
}

// ListLength 获取列表长度
func (r *Redis) ListLength(key string) (int64, error) {
	ctx := context.Background()
	return r.client.LLen(ctx, key).Result()
}

// SetAdd 添加一个或多个元素到集合
func (r *Redis) SetAdd(key string, members ...interface{}) error {
	ctx := context.Background()
	return r.client.SAdd(ctx, key, members...).Err()
}

// SetMembers 获取集合中的所有成员
func (r *Redis) SetMembers(key string) ([]string, error) {
	ctx := context.Background()
	return r.client.SMembers(ctx, key).Result()
}

// SetIsMember 判断元素是否在集合中
func (r *Redis) SetIsMember(key string, member interface{}) (bool, error) {
	ctx := context.Background()
	return r.client.SIsMember(ctx, key, member).Result()
}

// SetRemove 从集合中移除一个或多个元素
func (r *Redis) SetRemove(key string, members ...interface{}) error {
	ctx := context.Background()
	return r.client.SRem(ctx, key, members...).Err()
}

// ZSetAdd 向有序集合添加一个或多个成员，同时指定成员的分数
func (r *Redis) ZSetAdd(key string, members ...*redis.Z) error {
	ctx := context.Background()
	return r.client.ZAdd(ctx, key, members...).Err()
}

// ZSetRangeByScore 获取指定分数范围内的有序集合成员
func (r *Redis) ZSetRangeByScore(key string, opt *redis.ZRangeBy) ([]string, error) {
	ctx := context.Background()
	return r.client.ZRangeByScore(ctx, key, opt).Result()
}

// ZSetRank 获取有序集合指定成员的排名
func (r *Redis) ZSetRank(key, member string) (int64, error) {
	ctx := context.Background()
	return r.client.ZRank(ctx, key, member).Result()
}

// ZSetRemove 从有序集合中移除一个或多个成员
func (r *Redis) ZSetRemove(key string, members ...interface{}) error {
	ctx := context.Background()
	return r.client.ZRem(ctx, key, members...).Err()
}

// HashSet 设置哈希表字段的字符串值
func (r *Redis) HashSet(key, field, value string) error {
	ctx := context.Background()
	return r.client.HSet(ctx, key, field, value).Err()
}

// HashGet 获取哈希表字段的值
func (r *Redis) HashGet(key, field string) (string, error) {
	ctx := context.Background()
	return r.client.HGet(ctx, key, field).Result()
}

// HashGetAll 获取哈希表中所有字段和值
func (r *Redis) HashGetAll(key string) (map[string]string, error) {
	ctx := context.Background()
	return r.client.HGetAll(ctx, key).Result()
}

// HIncrBy 对哈希表指定字段进行增量操作
func (r *Redis) HIncrBy(key, field string, incr int64) (int64, error) {
	ctx := context.Background()
	return r.client.HIncrBy(ctx, key, field, incr).Result()
}

// HashDelete 删除哈希表中一个或多个字段
func (r *Redis) HashDelete(key string, fields ...string) error {
	ctx := context.Background()
	return r.client.HDel(ctx, key, fields...).Err()
}

// SetWithExpiration 设置带有过期时间的键值对
func (r *Redis) SetWithExpiration(key, value string, expiration time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, expiration).Err()
}

// GetExpiration 获取键的剩余过期时间
func (r *Redis) GetExpiration(key string) (time.Duration, error) {
	ctx := context.Background()
	return r.client.TTL(ctx, key).Result()
}

// Expire 设置键的过期时间
func (r *Redis) Expire(key string, expiration time.Duration) error {
	ctx := context.Background()
	return r.client.Expire(ctx, key, expiration).Err()
}

// Delete 删除一个或多个键
func (r *Redis) Delete(keys ...string) error {
	ctx := context.Background()
	return r.client.Del(ctx, keys...).Err()
}

// Incr 对键进行递增操作
func (r *Redis) Incr(key string) (int64, error) {
	ctx := context.Background()
	return r.client.Incr(ctx, key).Result()
}

// Exists 检查指定键是否存在于数据库中
func (r *Redis) Exists(key string) (int64, error) {
	ctx := context.Background()
	return r.client.Exists(ctx, key).Result()
}

// 关闭redis连接
func (r *Redis) Close() error {
	return r.client.Close()
}
