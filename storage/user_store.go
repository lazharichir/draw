package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/lazharichir/draw/core"
	"github.com/lazharichir/draw/storage/dbtx"
)

type IAMStore struct {
}

func NewIAMStorePG() *IAMStore {
	return &IAMStore{}
}

func (iam *IAMStore) SaveUser(ctx context.Context, db dbtx.DBTx, user core.User) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("users")
	ib.Cols("id", "status", "username", "email", "created_at", "last_signed_in_at", "last_drawn_pixel_at")
	ib.Values(user.ID, user.Status, user.Username, user.Email, user.CreatedAt, user.LastSignedInAt, user.LastDrawnPixelAt)
	ib.SQL(`
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			username = EXCLUDED.username,
			email = EXCLUDED.email,
			last_signed_in_at = EXCLUDED.last_signed_in_at,
			last_drawn_pixel_at = EXCLUDED.last_drawn_pixel_at
	`)

	query, args := ib.Build()
	_, err := db.ExecContext(ctx, query, args...)
	return err
}

func (iam *IAMStore) GetUserBy(ctx context.Context, db dbtx.DBTx, field string, value string) (*core.User, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "status", "username", "email", "created_at", "last_signed_in_at", "last_drawn_pixel_at").From("users")

	switch strings.ToLower(field) {
	case "id":
		sb.Where(sb.Equal("id", value))
	case "username":
		sb.Where(sb.Equal("username", value))
	case "email":
		sb.Where(sb.Equal("email", value))
	}

	query, args := sb.Build()
	row := db.QueryRowContext(ctx, query, args...)

	user := core.User{}
	if err := row.Scan(&user.ID, &user.Status, &user.Username, &user.Email, &user.CreatedAt, &user.LastSignedInAt, &user.LastDrawnPixelAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (iam *IAMStore) SaveUserProfile(ctx context.Context, db dbtx.DBTx, userID string, profile core.UserProfile) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("user_profiles")
	ib.Cols("user_id", "first_name", "last_name", "gender", "dob", "bio", "website_url", "facebook_url", "twitter_url", "instagram_url", "linkedin_url", "tiktok_url", "youtube_url")
	ib.Values(
		userID,
		profile.FirstName,
		profile.LastName,
		profile.Gender,
		profile.DateOfBirth,
		profile.Bio,
		profile.WebsiteURL,
		profile.FacebookURL,
		profile.TwitterURL,
		profile.InstagramURL,
		profile.LinkedInURL,
		profile.TikTokURL,
		profile.YouTubeURL,
	)

	ib.SQL(`
		ON CONFLICT (user_id) DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			gender = EXCLUDED.gender,
			dob = EXCLUDED.dob,
			bio = EXCLUDED.bio,
			website_url = EXCLUDED.website_url,
			facebook_url = EXCLUDED.facebook_url,
			twitter_url = EXCLUDED.twitter_url,
			instagram_url = EXCLUDED.instagram_url,
			linkedin_url = EXCLUDED.linkedin_url,
			tiktok_url = EXCLUDED.tiktok_url,
			youtube_url = EXCLUDED.youtube_url
	`)

	query, args := ib.Build()
	_, err := db.ExecContext(ctx, query, args...)
	return err
}

func (iam *IAMStore) GetUserProfile(ctx context.Context, db dbtx.DBTx, userID string) (*core.UserProfile, error) {
	lookup, err := iam.LoadUserProfiles(ctx, db, userID)
	if err != nil {
		return nil, err
	}

	if lookup[userID] == nil {
		return nil, fmt.Errorf("user profile not found for user %s", userID)
	}

	return lookup[userID], nil
}

func (iam *IAMStore) LoadUserProfiles(ctx context.Context, db dbtx.DBTx, ids ...string) (map[string]*core.UserProfile, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("user_id", "first_name", "last_name", "gender", "dob", "bio", "website_url", "facebook_url", "twitter_url", "instagram_url", "linkedin_url", "tiktok_url", "youtube_url")
	sb.From("user_profiles")
	sb.Where(sb.In("user_id", toAnySlice(ids)...))

	query, args := sb.Build()
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userProfiles := map[string]*core.UserProfile{}
	for rows.Next() {
		var userid string
		var profile core.UserProfile
		if err = rows.Scan(
			&userid,
			&profile.FirstName,
			&profile.LastName,
			&profile.Gender,
			&profile.DateOfBirth,
			&profile.Bio,
			&profile.WebsiteURL,
			&profile.FacebookURL,
			&profile.TwitterURL,
			&profile.InstagramURL,
			&profile.LinkedInURL,
			&profile.TikTokURL,
			&profile.YouTubeURL,
		); err != nil {
			return nil, err
		}

		userProfiles[userid] = &profile
	}

	return userProfiles, nil
}

func (iam *IAMStore) SaveVerificationToken(ctx context.Context, db dbtx.DBTx, token core.VerificationToken) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("verification_tokens")
	ib.Cols("token", "kind", "user_id", "email", "created_at", "expires_at", "used_at")
	ib.Values(token.Token, token.Kind, token.UserID, token.Email, token.CreatedAt, token.ExpiresAt, token.UsedAt)
	ib.SQL(`
		ON CONFLICT (token) DO UPDATE SET
			email = EXCLUDED.email,
			expires_at = EXCLUDED.expires_at,
			used_at = EXCLUDED.used_at
	`)

	query, args := ib.Build()
	_, err := db.ExecContext(ctx, query, args...)
	return err
}

func (iam *IAMStore) GetVerificationToken(ctx context.Context, db dbtx.DBTx, token string) (*core.VerificationToken, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("token", "kind", "user_id", "email", "created_at", "expires_at", "used_at")
	sb.From("verification_tokens")
	sb.Where(sb.Equal("token", token))

	query, args := sb.Build()
	row := db.QueryRowContext(ctx, query, args...)

	vt := core.VerificationToken{}
	if err := row.Scan(&vt.Token, &vt.Kind, &vt.UserID, &vt.Email, &vt.CreatedAt, &vt.ExpiresAt, &vt.UsedAt); err != nil {
		if err == ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &vt, nil
}
