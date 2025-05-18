// Copyright 2024 AliasMe
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package user

import (
	"context"
	"errors"
	"time"

	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/golgoth31/aliasme/internal/models"
	aliasme "github.com/golgoth31/aliasme/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service handles user-related operations
type Service struct {
	aliasme.UnimplementedUserServiceServer
	db *gorm.DB
}

// New creates a new user service
func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, req *aliasme.CreateUserRequest) (*aliasme.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, err
	}

	// Generate user ID
	id, err := generateID()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate user ID")
		return nil, err
	}

	user := &models.User{
		ID:        id,
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(user).Error; err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, err
	}

	return &aliasme.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(ctx context.Context, req *aliasme.GetUserRequest) (*aliasme.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", req.Id).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		return nil, err
	}

	return &aliasme.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// UpdateUser updates a user
func (s *Service) UpdateUser(ctx context.Context, req *aliasme.UpdateUserRequest) (*aliasme.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", req.Id).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		return nil, err
	}

	user.Username = req.Username
	user.Email = req.Email
	user.UpdatedAt = time.Now()

	if err := s.db.Save(&user).Error; err != nil {
		log.Error().Err(err).Msg("Failed to update user")
		return nil, err
	}

	return &aliasme.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(ctx context.Context, req *aliasme.DeleteUserRequest) (*aliasme.DeleteUserResponse, error) {
	if err := s.db.Delete(&models.User{}, "id = ?", req.Id).Error; err != nil {
		log.Error().Err(err).Msg("Failed to delete user")
		return nil, err
	}

	return &aliasme.DeleteUserResponse{Success: true}, nil
}

// generateID generates a random ID using xid
func generateID() (string, error) {
	return xid.New().String(), nil
}

// GetUserByEmail retrieves a user ID by email
func (s *Service) GetUserByEmail(ctx context.Context, req *aliasme.GetUserByEmailRequest) (*aliasme.GetUserByEmailResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &aliasme.GetUserByEmailResponse{
		UserId: user.ID,
	}, nil
}
