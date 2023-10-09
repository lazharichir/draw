package core

import (
	"fmt"
	"time"

	"github.com/lazharichir/draw/utils"
)

type User struct {
	ID               string
	Status           string // probation, active, suspended
	Username         string
	Email            string
	CreatedAt        time.Time
	LastSignedInAt   time.Time
	LastDrawnPixelAt time.Time
}

type UserProfile struct {
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth time.Time
	Bio         string
	//
	WebsiteURL   string
	FacebookURL  string
	TwitterURL   string
	InstagramURL string
	LinkedInURL  string
	TikTokURL    string
	YouTubeURL   string
}

func NewVerificationToken(userID string, kind string) (VerificationToken, error) {
	vt := VerificationToken{
		Kind:      kind,
		UserID:    userID,
		CreatedAt: time.Now(),
		UsedAt:    nil,
	}

	switch kind {
	case "signup":
		vt.Token = utils.NewVerificationTokenSignup()
		vt.ExpiresAt = vt.CreatedAt.Add(30 * time.Minute)
	case "signin":
		vt.Token = utils.NewVerificationTokenSignin()
		vt.ExpiresAt = vt.CreatedAt.Add(30 * time.Minute)
	case "change_email":
		vt.Token = utils.NewVerificationTokenChangeEmail()
		vt.ExpiresAt = vt.CreatedAt.Add(30 * time.Minute)
	case "reset_password":
		vt.Token = utils.NewVerificationTokenResetPassword()
		vt.ExpiresAt = vt.CreatedAt.Add(30 * time.Minute)
	default:
		return vt, fmt.Errorf("invalid verification token kind '%s'", kind)
	}

	return vt, nil
}

type VerificationToken struct {
	Token     string
	Kind      string
	UserID    string
	Email     *string
	CreatedAt time.Time
	ExpiresAt time.Time
	UsedAt    *time.Time
}

func (vt VerificationToken) IsActive() bool {
	return !vt.hasBeenUsed() && !vt.hasExpired()
}

func (vt *VerificationToken) Used() {
	now := time.Now()
	vt.UsedAt = &now
}

func (vt *VerificationToken) SetEmail(email string) {
	vt.Email = &email
}

func (vt VerificationToken) hasBeenUsed() bool {
	return vt.UsedAt != nil
}

func (vt VerificationToken) hasExpired() bool {
	return vt.ExpiresAt.Before(time.Now())
}
