// oilan/internal/infrastructure/repository/postgres/dialog_postgres.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"oilan/internal/domain"
	"oilan/internal/domain/repository"
	"time"
)

// dialogRepo implements the repository.DialogRepository interface.
type dialogRepo struct {
	db *sql.DB
}

// NewDialogRepository creates a new instance of the dialog repository.
func NewDialogRepository(db *sql.DB) repository.DialogRepository {
	return &dialogRepo{db: db}
}

// Save creates a new dialog session.
func (r *dialogRepo) Save(ctx context.Context, dialog *domain.Dialog) error {
	query := `
        INSERT INTO dialogs (user_id, title, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id;
    `
	now := time.Now()
	dialog.CreatedAt = now
	dialog.UpdatedAt = now
	
	return r.db.QueryRowContext(ctx, query, dialog.UserID, dialog.Title, dialog.CreatedAt, dialog.UpdatedAt).Scan(&dialog.ID)
}

// AddMessage adds a new message to an existing dialog and updates the dialog's timestamp.
func (r *dialogRepo) AddMessage(ctx context.Context, message *domain.Message) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	msgQuery := `
        INSERT INTO messages (dialog_id, role, content, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id;
    `
	message.CreatedAt = time.Now()
	err = tx.QueryRowContext(ctx, msgQuery, message.DialogID, message.Role, message.Content, message.CreatedAt).Scan(&message.ID)
	if err != nil {
		return err
	}

	dialogQuery := `UPDATE dialogs SET updated_at = $1 WHERE id = $2;`
	_, err = tx.ExecContext(ctx, dialogQuery, message.CreatedAt, message.DialogID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// FindByID finds a single dialog with all its messages.
func (r *dialogRepo) FindByID(ctx context.Context, id int64) (*domain.Dialog, error) {
	dialogQuery := `SELECT id, user_id, title, created_at, updated_at FROM dialogs WHERE id = $1;`
	dialog := &domain.Dialog{}
	err := r.db.QueryRowContext(ctx, dialogQuery, id).Scan(
		&dialog.ID, &dialog.UserID, &dialog.Title, &dialog.CreatedAt, &dialog.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	messagesQuery := `SELECT id, dialog_id, role, content, created_at FROM messages WHERE dialog_id = $1 ORDER BY created_at ASC;`
	rows, err := r.db.QueryContext(ctx, messagesQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(&msg.ID, &msg.DialogID, &msg.Role, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	dialog.Messages = messages
	return dialog, nil
}

// FindAllByUserID finds all dialogs for a specific user (without messages for performance).
func (r *dialogRepo) FindAllByUserID(ctx context.Context, userID int64) ([]*domain.Dialog, error) {
	query := `SELECT id, user_id, title, created_at, updated_at FROM dialogs WHERE user_id = $1 ORDER BY updated_at DESC;`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dialogs []*domain.Dialog
	for rows.Next() {
		var dialog domain.Dialog
		if err := rows.Scan(&dialog.ID, &dialog.UserID, &dialog.Title, &dialog.CreatedAt, &dialog.UpdatedAt); err != nil {
			return nil, err
		}
		dialogs = append(dialogs, &dialog)
	}

	return dialogs, nil
}

// GetAll retrieves all dialogs from the database.
func (r *dialogRepo) GetAll(ctx context.Context) ([]*domain.Dialog, error) {
	query := `SELECT id, user_id, title, created_at, updated_at FROM dialogs ORDER BY updated_at DESC;`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dialogs []*domain.Dialog
	for rows.Next() {
		var dialog domain.Dialog
		if err := rows.Scan(
			&dialog.ID, &dialog.UserID, &dialog.Title, &dialog.CreatedAt, &dialog.UpdatedAt,
		); err != nil {
			return nil, err
		}
		dialogs = append(dialogs, &dialog)
	}
	return dialogs, nil
}