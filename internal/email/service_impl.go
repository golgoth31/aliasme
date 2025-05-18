package email

import (
	"context"
	"time"

	"github.com/golgoth31/aliasme/internal/models"
	"github.com/golgoth31/aliasme/internal/ovh"
	aliasme "github.com/golgoth31/aliasme/pkg/proto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// EmailService handles email-related operations
type EmailService struct {
	aliasme.UnimplementedEmailServiceServer
	db           *gorm.DB
	ovhClient    *ovh.Client
	emailService *Service
}

// NewEmailService creates a new email service
func NewEmailService(db *gorm.DB, ovhClient *ovh.Client, emailService *Service) *EmailService {
	return &EmailService{
		db:           db,
		ovhClient:    ovhClient,
		emailService: emailService,
	}
}

// RegisterEmail registers a new email address for a user
func (s *EmailService) RegisterEmail(ctx context.Context, req *aliasme.RegisterEmailRequest) (*aliasme.Email, error) {
	// Generate verification token
	token, err := generateToken()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate verification token")
		return nil, err
	}

	// Generate email ID
	id, err := generateID()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate email ID")
		return nil, err
	}

	email := &models.Email{
		ID:        id,
		UserID:    req.UserId,
		Address:   req.EmailAddress,
		Verified:  false,
		Token:     token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(email).Error; err != nil {
		log.Error().Err(err).Msg("Failed to create email")
		return nil, err
	}

	// Send verification email
	if err := s.emailService.SendVerificationEmail(email.Address, token); err != nil {
		log.Error().Err(err).Msg("Failed to send verification email")
		return nil, err
	}

	return &aliasme.Email{
		Id:        email.ID,
		UserId:    email.UserID,
		Address:   email.Address,
		Verified:  email.Verified,
		CreatedAt: timestamppb.New(email.CreatedAt),
		UpdatedAt: timestamppb.New(email.UpdatedAt),
	}, nil
}

// VerifyEmail verifies an email address using the token
func (s *EmailService) VerifyEmail(ctx context.Context, req *aliasme.VerifyEmailRequest) (*aliasme.Email, error) {
	var email models.Email
	if err := s.db.First(&email, "token = ?", req.Token).Error; err != nil {
		log.Error().Err(err).Msg("Failed to find email with token")
		return nil, err
	}

	email.Verified = true
	email.Token = ""
	email.UpdatedAt = time.Now()

	if err := s.db.Save(&email).Error; err != nil {
		log.Error().Err(err).Msg("Failed to update email")
		return nil, err
	}

	return &aliasme.Email{
		Id:        email.ID,
		UserId:    email.UserID,
		Address:   email.Address,
		Verified:  email.Verified,
		CreatedAt: timestamppb.New(email.CreatedAt),
		UpdatedAt: timestamppb.New(email.UpdatedAt),
	}, nil
}

// CreateAlias creates a new email alias
func (s *EmailService) CreateAlias(ctx context.Context, req *aliasme.CreateAliasRequest) (*aliasme.Alias, error) {
	// Verify that the email belongs to the user and is verified
	var email models.Email
	if err := s.db.First(&email, "id = ? AND user_id = ? AND verified = ?", req.EmailId, req.UserId, true).Error; err != nil {
		log.Error().Err(err).Msg("Failed to find verified email")
		return nil, err
	}

	// Generate alias ID
	id, err := generateID()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate alias ID")
		return nil, err
	}

	// Create alias address
	aliasAddress := req.AliasPrefix + "@yourdomain.com"

	// Create alias in OVH
	if err := s.ovhClient.CreateEmailAlias("yourdomain.com", req.AliasPrefix, email.Address); err != nil {
		log.Error().Err(err).Msg("Failed to create alias in OVH")
		return nil, err
	}

	alias := &models.Alias{
		ID:           id,
		UserID:       req.UserId,
		EmailID:      req.EmailId,
		AliasAddress: aliasAddress,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(alias).Error; err != nil {
		log.Error().Err(err).Msg("Failed to create alias")
		return nil, err
	}

	return &aliasme.Alias{
		Id:           alias.ID,
		UserId:       alias.UserID,
		EmailId:      alias.EmailID,
		AliasAddress: alias.AliasAddress,
		CreatedAt:    timestamppb.New(alias.CreatedAt),
		UpdatedAt:    timestamppb.New(alias.UpdatedAt),
	}, nil
}

// ListAliases lists all aliases for a user
func (s *EmailService) ListAliases(ctx context.Context, req *aliasme.ListAliasesRequest) (*aliasme.ListAliasesResponse, error) {
	var aliases []models.Alias
	if err := s.db.Find(&aliases, "user_id = ?", req.UserId).Error; err != nil {
		log.Error().Err(err).Msg("Failed to list aliases")
		return nil, err
	}

	protoAliases := make([]*aliasme.Alias, len(aliases))
	for i, alias := range aliases {
		protoAliases[i] = &aliasme.Alias{
			Id:           alias.ID,
			UserId:       alias.UserID,
			EmailId:      alias.EmailID,
			AliasAddress: alias.AliasAddress,
			CreatedAt:    timestamppb.New(alias.CreatedAt),
			UpdatedAt:    timestamppb.New(alias.UpdatedAt),
		}
	}

	return &aliasme.ListAliasesResponse{
		Aliases: protoAliases,
	}, nil
}

// DeleteAlias deletes an alias
func (s *EmailService) DeleteAlias(ctx context.Context, req *aliasme.DeleteAliasRequest) (*aliasme.DeleteAliasResponse, error) {
	if err := s.db.Delete(&models.Alias{}, "id = ?", req.Id).Error; err != nil {
		log.Error().Err(err).Msg("Failed to delete alias")
		return nil, err
	}

	return &aliasme.DeleteAliasResponse{Success: true}, nil
}

// UpdateAlias updates an alias
func (s *EmailService) UpdateAlias(ctx context.Context, req *aliasme.UpdateAliasRequest) (*aliasme.Alias, error) {
	var alias models.Alias
	if err := s.db.First(&alias, "id = ?", req.Id).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get alias")
		return nil, err
	}

	alias.EmailID = req.EmailId
	alias.AliasAddress = req.AliasPrefix
	alias.UpdatedAt = time.Now()

	if err := s.db.Save(&alias).Error; err != nil {
		log.Error().Err(err).Msg("Failed to update alias")
		return nil, err
	}

	return &aliasme.Alias{
		Id:           alias.ID,
		UserId:       alias.UserID,
		EmailId:      alias.EmailID,
		AliasAddress: alias.AliasAddress,
		CreatedAt:    timestamppb.New(alias.CreatedAt),
		UpdatedAt:    timestamppb.New(alias.UpdatedAt),
	}, nil
}

// generateToken generates a random verification token
func generateToken() (string, error) {
	return uuid.New().String(), nil
}
