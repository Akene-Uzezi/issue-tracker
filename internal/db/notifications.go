package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationModel struct {
	DB *pgxpool.Pool
}

type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Ref       string    `json:"ref,omitempty"`
	RelatedID int       `json:"relatedId,omitempty"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"createdAt"`
}

// CreateForRoles inserts a notification row for every active (non-disabled) user
// whose role is in the provided list. Used to alert superadmins, admins and
// managers when a death report is submitted, regardless of which frontend they
// are currently using.
func (m *NotificationModel) CreateForRoles(ctx context.Context, roles []string, n Notification) (int64, error) {
	if len(roles) == 0 {
		return 0, nil
	}
	query := `
		INSERT INTO notifications (user_id, type, title, message, ref, related_id)
		SELECT id, $1, $2, $3, $4, $5
		FROM users
		WHERE role = ANY(string_to_array($6, ',')) AND disabled = FALSE
	`
	ref := sql.NullString{String: n.Ref, Valid: n.Ref != ""}
	relatedID := sql.NullInt64{Int64: int64(n.RelatedID), Valid: n.RelatedID != 0}
	tag, err := m.DB.Exec(ctx, query, n.Type, n.Title, n.Message, ref, relatedID, strings.Join(roles, ","))
	if err != nil {
		return 0, fmt.Errorf("database query error: %w", err)
	}
	return tag.RowsAffected(), nil
}

func (m *NotificationModel) GetForUser(ctx context.Context, userID int) ([]Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, ref, related_id, read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT 100
	`
	rows, err := m.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	notifications, err := pgx.CollectRows(rows, pgx.RowToStructByName[Notification])
	if err != nil {
		return nil, fmt.Errorf("failed to scan notifications: %w", err)
	}
	return notifications, nil
}

func (m *NotificationModel) MarkAllRead(ctx context.Context, userID int) error {
	_, err := m.DB.Exec(ctx, `UPDATE notifications SET read = TRUE WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("database query error: %w", err)
	}
	return nil
}

func (m *NotificationModel) MarkRead(ctx context.Context, id, userID int) error {
	_, err := m.DB.Exec(ctx, `UPDATE notifications SET read = TRUE WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("database query error: %w", err)
	}
	return nil
}

func (m *NotificationModel) DeleteAll(ctx context.Context, userID int) error {
	_, err := m.DB.Exec(ctx, `DELETE FROM notifications WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("database query error: %w", err)
	}
	return nil
}
