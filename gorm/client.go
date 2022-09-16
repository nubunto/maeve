package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nubunto/maeve"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type KeyValue struct {
	BaseModel

	Key string
	Value string
}

type Client struct {
	db *gorm.DB
}

func New(dsn string) (*Client, error) {
	cfg := &gorm.Config{
		QueryFields: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "maeve_",
			SingularTable: false,
		},
	}
	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, fmt.Errorf("error opening maeve Gorm client: %w", err)
	}

	db.AutoMigrate(&KeyValue{})

	return &Client{db: db.Debug()}, nil
}

func (c *Client) Fetch(ctx context.Context, path maeve.StringPath) (maeve.KeyValueList, error) {
	tx := c.db.WithContext(ctx)

	var keys []KeyValue
	err := tx.Where("starts_with(key, ?)", string(maeve.TrimDynamic(path))).Find(&keys).Error
	if err != nil {
		return nil, err
	}

	kv := make(maeve.KeyValueList, 0)
	for _, k := range keys {
		kv = append(kv, maeve.KeyValue{
			Path: k.Key,
			Value: k.Value,
		})
	}

	return kv, nil
}

func (c *Client) Put(ctx context.Context, kv maeve.KeyValueList) error {
	tx := c.db.WithContext(ctx)

	var keys []KeyValue
	for _, k := range kv {
		keys = append(keys, KeyValue{
			Key: k.Path,
			Value: k.Value,
		})
	}

	return tx.Create(keys).Error
}

func (c *Client) Delete(ctx context.Context, path maeve.StringPath) error {
	tx := c.db.WithContext(ctx)

	var keys []KeyValue

	tx.Where("starts_with(key, ?)", maeve.TrimDynamic(path)).Find(&keys)
	return tx.Delete(&keys).Error
}
