package userservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/taaanechka/rest-api-go/internal/api-server/services/ports/userstorage"
	"github.com/taaanechka/rest-api-go/internal/apperror"
	"github.com/taaanechka/rest-api-go/pkg/logging"
)

type Service struct {
	storage userstorage.Storage
	lg      *logging.Logger
}

func NewService(lg *logging.Logger, storage userstorage.Storage) *Service {
	return &Service{
		storage: storage,
		lg: lg,
	}
}

func (s *Service) Create(ctx context.Context, uBL userstorage.User) (string, error) {
	oid, err := s.storage.Create(ctx, uBL)
	if err != nil {
		return oid, fmt.Errorf("failed to create user: %w", err)
	}
	return oid, nil
}

func (s *Service) GetList(ctx context.Context) ([]userstorage.User, error) {
	us, err := s.storage.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return us, nil
}

func (s *Service) GetByUUID(ctx context.Context, id string) (userstorage.User, error) {
	u, err := s.storage.FindOne(ctx, id)
	if err != nil {
		return userstorage.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return u, nil
}

func (s *Service) Update(ctx context.Context, id string, user userstorage.User) (error, bool) {
	found := true

	_, err := s.storage.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			found = false
		} else {
			return err, false
		}
	}

	if found {
		if err = s.storage.Update(ctx, id, user); err != nil {
			return fmt.Errorf("failed to update user: %w", err), false
		}
	} else {
		_, err = s.storage.Create(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err), found
		}
	}
	return nil, found
}

func (s *Service) Patch(ctx context.Context, id string, user userstorage.User) error {
	u, err := s.storage.FindOne(ctx, id)
	if err != nil {
		return err
	}

	var uPatch userstorage.User
	uPatch.Email = u.Email
	if user.Email != "" {
		uPatch.Email = user.Email
	}

	uPatch.Username = u.Username
	if user.Username != "" {
		uPatch.Username = user.Username
	}

	uPatch.PasswordHash = u.PasswordHash
	if user.PasswordHash != "" {
		uPatch.PasswordHash = user.PasswordHash
	}

	if err := s.storage.Update(ctx, id, uPatch); err != nil {
		return fmt.Errorf("failed to partiupdate user: %w", err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	_, err := s.storage.FindOne(ctx, id)
	if err != nil {
		return err
	}

	if err := s.storage.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to partially update user: %w", err)
	}
	return nil
}
