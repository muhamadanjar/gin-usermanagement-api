package utils

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Inisialisasi Redis Client (Singleton)
var (
	ctx         = context.Background()
	redisClient *redis.Client
)

// InitRedis menghubungkan ke Redis
func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Ganti sesuai konfigurasi Redis
		Password: "",               // Kosongkan jika tidak ada password
		DB:       0,                // Gunakan database default
	})

	// Cek koneksi
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Gagal konek ke Redis:", err)
	}
	log.Println("Terhubung ke Redis")
}

// GetRedisClient mengembalikan instance Redis client
func GetRedisClient() *redis.Client {
	return redisClient
}

// SetCache menyimpan data ke Redis dengan TTL
func SetCache(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisClient.Set(ctx, key, data, ttl).Err()
}

// GetCache mengambil data dari Redis
func GetCache(key string, dest interface{}) (bool, error) {
	data, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil // Tidak ada data di Redis
	} else if err != nil {
		return false, err
	}

	// Decode JSON ke struct
	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteCache menghapus cache di Redis
func DeleteCache(key string) error {
	return redisClient.Del(ctx, key).Err()
}
