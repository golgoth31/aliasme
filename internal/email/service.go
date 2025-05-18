package email

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/golgoth31/aliasme/internal/models"
	aliasme "github.com/golgoth31/aliasme/pkg/proto"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// Config holds the email service configuration
type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

// Service handles email-related operations
type Service struct {
	aliasme.UnimplementedEmailServiceServer
	db     *gorm.DB
	config Config
}

// New creates a new email service
func New(db *gorm.DB, cfg Config) *Service {
	return &Service{db: db, config: cfg}
}

// CreateAlias creates a new email alias
func (s *Service) CreateAlias(ctx context.Context, req *aliasme.CreateAliasRequest) (*aliasme.Alias, error) {
	// Generate alias ID
	id, err := generateID()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate alias ID")
		return nil, err
	}

	alias := &models.Alias{
		ID:           id,
		UserID:       req.UserId,
		EmailID:      req.GetEmailId(),
		AliasAddress: req.GetAliasPrefix(),
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

// GetAlias retrieves an alias by ID
func (s *Service) GetAlias(ctx context.Context, req *aliasme.GetAliasRequest) (*aliasme.Alias, error) {
	var alias models.Alias
	if err := s.db.First(&alias, "id = ?", req.Id).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get alias")
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

// UpdateAlias updates an alias
func (s *Service) UpdateAlias(ctx context.Context, req *aliasme.UpdateAliasRequest) (*aliasme.Alias, error) {
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

// DeleteAlias deletes an alias
func (s *Service) DeleteAlias(ctx context.Context, req *aliasme.DeleteAliasRequest) (*aliasme.DeleteAliasResponse, error) {
	if err := s.db.Delete(&models.Alias{}, "id = ?", req.Id).Error; err != nil {
		log.Error().Err(err).Msg("Failed to delete alias")
		return nil, err
	}

	return &aliasme.DeleteAliasResponse{Success: true}, nil
}

// ListAliases lists all aliases for a user
func (s *Service) ListAliases(ctx context.Context, req *aliasme.ListAliasesRequest) (*aliasme.ListAliasesResponse, error) {
	var aliases []models.Alias
	if err := s.db.Where("user_id = ?", req.UserId).Find(&aliases).Error; err != nil {
		log.Error().Err(err).Msg("Failed to list aliases")
		return nil, err
	}

	response := &aliasme.ListAliasesResponse{
		Aliases: make([]*aliasme.Alias, len(aliases)),
	}

	for i, alias := range aliases {
		response.Aliases[i] = &aliasme.Alias{
			Id:           alias.ID,
			UserId:       alias.UserID,
			EmailId:      alias.EmailID,
			AliasAddress: alias.AliasAddress,
			CreatedAt:    timestamppb.New(alias.CreatedAt),
			UpdatedAt:    timestamppb.New(alias.UpdatedAt),
		}
	}

	return response, nil
}

// generateID generates a random ID using xid
func generateID() (string, error) {
	return xid.New().String(), nil
}

// SendVerificationEmail sends a verification email to the user
func (s *Service) SendVerificationEmail(to, token string) error {
	subject := "Verify your email address"
	body := fmt.Sprintf(`
		Hello,

		Please verify your email address by clicking the following link:
		%s/verify?token=%s

		This link will expire in 24 hours.

		Best regards,
		The AliasMe Team
	`, os.Getenv("BASE_URL"), token)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n",
		s.config.FromEmail, to, subject, body)

	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword,
		s.config.SMTPHost)

	err := smtp.SendMail(
		s.config.SMTPHost+":"+s.config.SMTPPort,
		auth,
		s.config.FromEmail,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		log.Error().Err(err).Msg("Failed to send verification email")
		return err
	}

	return nil
}
