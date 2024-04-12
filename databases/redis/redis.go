package redis

import (
	"context"
	"github.com/BarnabyCharles/frame/config"
	"github.com/ghodss/yaml"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

func WithRedis(address string, hand func(cli *redis.Client) (interface{}, error)) (interface{}, error) {

	err := config.ViperInit(address)
	if err != nil {
		return "", err
	}
	name := viper.GetString("Nacos.ServerName")
	group := viper.GetString("Nacos.group")
	val, err := config.GetNacosConfig(name, group)
	if err != nil {
		return "", err
	}
	var r config.AppConfig
	yaml.Unmarshal([]byte(val), &r)
	cli := redis.NewClient(&redis.Options{
		Addr: r.Redis.Host + ":" + r.Redis.Port,
		DB:   int(r.Redis.DB),
	})
	// 设置AOF同步方式为always，每次写入都同步
	_, err = cli.ConfigSet(context.Background(), "appendfsync", "always").Result()
	if err != nil {
		panic(err)
	}

	// 等待AOF文件同步完成
	_, err = cli.Do(context.Background(), "CLIENT", "REPLY", "ON").Result()
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	str, err := hand(cli)
	if err != nil {
		return "", err
	}
	return str, nil
}

func SetByVal(address, key string, val interface{}, duration time.Duration) error {
	_, err := WithRedis(address, func(cli *redis.Client) (interface{}, error) {
		err := cli.Set(context.Background(), key, val, duration).Err()
		if err != nil {
			return "", err
		}
		return "", nil
	})
	if err != nil {
		return err
	}
	return nil
}

func GetKey(address, key string) (string, error) {
	val, err := WithRedis(address, func(cli *redis.Client) (interface{}, error) {
		val, err := cli.Get(context.Background(), key).Result()
		if err != nil {
			return "", err
		}
		return val, nil
	})
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

func Exists(address, key string) bool {
	ok, _ := WithRedis(address, func(cli *redis.Client) (interface{}, error) {
		ok, err := cli.Exists(context.Background(), key).Result()
		if err != nil {
			panic(err)
		}
		return strconv.FormatInt(ok, 10), nil
	})
	if ok.(string) == "1" {
		return true

	}
	return false
}

func KeyDelete(address, key string) error {
	_, err := WithRedis(address, func(cli *redis.Client) (interface{}, error) {
		err := cli.Del(context.Background(), key).Err()
		return "", err
	})

	if err != nil {
		return err
	}
	return nil
}

func Lock(address, key string, val interface{}, duration time.Duration, isReadOnly bool) (bool, error) {
	res := false
	WithRedis(address, func(cli *redis.Client) (interface{}, error) {
		if !isReadOnly {
			for {
				result, err := cli.SetNX(context.Background(), key, val, duration).Result()
				if err != nil {
					return "", err
				}
				res = true
				if result {
					return "", nil
				}
			}
		}
		re, err := cli.SetNX(context.Background(), key, val, duration).Result()
		res = re
		if err != nil {
			return "", err
		}
		return "", nil
	})
	return res, nil
}

func UnLock(address, key string) error {
	_, err := WithRedis(address, func(cli *redis.Client) (interface{}, error) {
		return "", cli.Del(context.Background(), key).Err()
	})
	if err != nil {
		return err
	}
	return nil
}

func SetHashValue(address, key, field string, value interface{}) error {
	_, err := WithRedis(address, func(cli *redis.Client) (interface{}, error) {
		err := cli.HSet(context.Background(), key, field, value).Err()
		return "", err
	})
	return err
}

func GetHashValue(address, key, field string) (interface{}, error) {
	withRedis, err := WithRedis(address, func(cli *redis.Client) (interface{}, error) {
		return cli.HGet(context.Background(), key, field).Result()
	})
	if err != nil {
		return err, nil
	}

	return withRedis, err
}

func GetKeysInfos(address, key string) ([]string, error) {
	withRedis, err := WithRedis(address, func(cli *redis.Client) (interface{}, error) {
		return cli.HKeys(context.Background(), key).Result()
	})
	if err != nil {
		return nil, err
	}
	return withRedis.([]string), nil
}

func DecrBy(address, key string, val int64) (int64, error) {
	data, err := WithRedis(address, func(cli *redis.Client) (interface{}, error) {
		return cli.DecrBy(context.Background(), key, val).Result()
	})

	if err != nil {
		return 0, err
	}

	return data.(int64), nil

}

func IncrBy(address, key string, val int64) (int64, error) {
	data, err := WithRedis(address, func(cli *redis.Client) (interface{}, error) {
		return cli.IncrBy(context.Background(), key, val).Result()
	})

	if err != nil {
		return 0, err
	}

	return data.(int64), nil
}
