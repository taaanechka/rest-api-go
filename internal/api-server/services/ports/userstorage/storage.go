package userstorage

import (
	"context"
)

type Config struct {
	Host       string   `yaml:"host"`
	Port       string   `yaml:"port"`
	Database   string   `yaml:"database"`
	AuthDB     string   `yaml:"auth_db"`
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
	Collection string   `yaml:"collection"`
	Indexes    []string `yaml:"indexes"`
}

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	FindAll(ctx context.Context) ([]User, error)
	FindOne(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, id string, user User) error
	Delete(ctx context.Context, id string) error
}
