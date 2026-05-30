package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/therizza/freakdream/backend/internal/models"
)

var ErrNotFound = errors.New("not found")

type UserRepository struct{ db *sqlx.DB }

func NewUserRepository(db *sqlx.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	err := r.db.GetContext(ctx, &u, "SELECT id, nome, sobre, email, senha, created_at FROM usuario WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.GetContext(ctx, &u, "SELECT id, nome, sobre, email, senha, created_at FROM usuario WHERE email = ?", email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *UserRepository) Create(ctx context.Context, nome, sobre, email, senhaHash string) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO usuario (nome, sobre, email, senha, created_at) VALUES (?, ?, ?, ?, NOW())",
		nome, sobre, email, senhaHash)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *UserRepository) Search(ctx context.Context, query string) ([]models.User, error) {
	q := "%" + strings.ToLower(query) + "%"
	var users []models.User
	err := r.db.SelectContext(ctx, &users,
		"SELECT id, nome, sobre, email, senha, created_at FROM usuario WHERE LOWER(nome) LIKE ? OR LOWER(sobre) LIKE ?", q, q)
	return users, err
}

func (r *UserRepository) RecordLogin(ctx context.Context, userID int64, ip string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO login_historico (user_id, data, ip) VALUES (?, NOW(), ?)", userID, ip)
	return err
}

// PostRepository ---------------------------------------------------------------

type PostRepository struct{ db *sqlx.DB }

func NewPostRepository(db *sqlx.DB) *PostRepository { return &PostRepository{db: db} }

func (r *PostRepository) Create(ctx context.Context, userID int64, tipo models.PostType, texto string) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO post (user_id, tipo, texto, data) VALUES (?, ?, ?, NOW())", userID, tipo, texto)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *PostRepository) GetByID(ctx context.Context, postID int64) (*models.Post, error) {
	var p models.Post
	err := r.db.GetContext(ctx, &p, "SELECT post_id, user_id, tipo, texto, data FROM post WHERE post_id = ?", postID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &p, err
}

func (r *PostRepository) Delete(ctx context.Context, postID, userID int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM post WHERE post_id = ? AND user_id = ?", postID, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostRepository) Feed(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.SelectContext(ctx, &posts, `
		SELECT DISTINCT p.post_id, p.user_id, p.tipo, p.texto, p.data
		FROM post p
		JOIN seguir s ON s.seguidor_id = p.user_id
		WHERE s.user_id = ? OR p.user_id = ?
		ORDER BY p.post_id DESC
		LIMIT ? OFFSET ?`, userID, userID, limit, offset)
	return posts, err
}

func (r *PostRepository) GetByUser(ctx context.Context, userID int64, limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.SelectContext(ctx, &posts,
		"SELECT post_id, user_id, tipo, texto, data FROM post WHERE user_id = ? ORDER BY post_id DESC LIMIT ? OFFSET ?",
		userID, limit, offset)
	return posts, err
}

// PhotoRepository ---------------------------------------------------------------

type PhotoRepository struct{ db *sqlx.DB }

func NewPhotoRepository(db *sqlx.DB) *PhotoRepository { return &PhotoRepository{db: db} }

func (r *PhotoRepository) Create(ctx context.Context, foto, texto string) (int64, error) {
	res, err := r.db.ExecContext(ctx, "INSERT INTO fotos (foto, texto) VALUES (?, ?)", foto, texto)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *PhotoRepository) GetByID(ctx context.Context, id int64) (*models.Photo, error) {
	var p models.Photo
	err := r.db.GetContext(ctx, &p, "SELECT foto_id, foto, texto FROM fotos WHERE foto_id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &p, err
}

func (r *PhotoRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM fotos WHERE foto_id = ?", id)
	return err
}

func (r *PhotoRepository) GetProfilePhoto(ctx context.Context, userID int64) (*models.Photo, error) {
	var pp models.ProfilePhoto
	err := r.db.GetContext(ctx, &pp,
		"SELECT foto_id, user_id FROM perfil_foto WHERE user_id = ? ORDER BY foto_id DESC LIMIT 1", userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, pp.FotoID)
}

func (r *PhotoRepository) SetProfilePhoto(ctx context.Context, userID, fotoID int64) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO perfil_foto (user_id, foto_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE foto_id = ?",
		userID, fotoID, fotoID)
	return err
}

func (r *PhotoRepository) DeleteProfilePhoto(ctx context.Context, userID, fotoID int64) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM perfil_foto WHERE user_id = ? AND foto_id = ?", userID, fotoID)
	return err
}

// FollowRepository ---------------------------------------------------------------

type FollowRepository struct{ db *sqlx.DB }

func NewFollowRepository(db *sqlx.DB) *FollowRepository { return &FollowRepository{db: db} }

func (r *FollowRepository) Follow(ctx context.Context, userID, seguidorID int64) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT IGNORE INTO seguir (user_id, seguidor_id) VALUES (?, ?)", userID, seguidorID)
	return err
}

func (r *FollowRepository) Unfollow(ctx context.Context, userID, seguidorID int64) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM seguir WHERE user_id = ? AND seguidor_id = ?", userID, seguidorID)
	return err
}

func (r *FollowRepository) Count(ctx context.Context, userID int64) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM seguir WHERE user_id = ?", userID)
	return count, err
}

func (r *FollowRepository) IsFollowing(ctx context.Context, userID, seguidorID int64) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM seguir WHERE user_id = ? AND seguidor_id = ?", userID, seguidorID)
	return count > 0, err
}

// FriendRepository ---------------------------------------------------------------

type FriendRepository struct{ db *sqlx.DB }

func NewFriendRepository(db *sqlx.DB) *FriendRepository { return &FriendRepository{db: db} }

func (r *FriendRepository) SendRequest(ctx context.Context, userID, amigoID int64) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT IGNORE INTO amigo (user_id, amigo_id, status, criado_em) VALUES (?, ?, 'pending', NOW())",
		userID, amigoID)
	return err
}

func (r *FriendRepository) Accept(ctx context.Context, userID, amigoID int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(ctx,
		"UPDATE amigo SET status = 'accepted' WHERE user_id = ? AND amigo_id = ?", amigoID, userID)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		"INSERT IGNORE INTO amigo (user_id, amigo_id, status, criado_em) VALUES (?, ?, 'accepted', NOW())",
		userID, amigoID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *FriendRepository) Remove(ctx context.Context, userID, amigoID int64) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM amigo WHERE (user_id = ? AND amigo_id = ?) OR (user_id = ? AND amigo_id = ?)",
		userID, amigoID, amigoID, userID)
	return err
}

func (r *FriendRepository) IsFriend(ctx context.Context, userID, amigoID int64) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM amigo WHERE user_id = ? AND amigo_id = ? AND status = 'accepted'",
		userID, amigoID)
	return count > 0, err
}

func (r *FriendRepository) ListRequests(ctx context.Context, userID int64) ([]models.Friend, error) {
	var friends []models.Friend
	err := r.db.SelectContext(ctx, &friends,
		"SELECT user_id, amigo_id, status, criado_em FROM amigo WHERE amigo_id = ? AND status = 'pending'", userID)
	return friends, err
}

// NotificationRepository ---------------------------------------------------------------

type NotificationRepository struct{ db *sqlx.DB }

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Send(ctx context.Context, userID int64, tipo, texto string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO notificacao (user_id, tipo, texto, data, lida) VALUES (?, ?, ?, NOW(), false)",
		userID, tipo, texto)
	return err
}

func (r *NotificationRepository) List(ctx context.Context, userID int64) ([]models.Notification, error) {
	var notifs []models.Notification
	err := r.db.SelectContext(ctx, &notifs,
		"SELECT id, user_id, tipo, texto, data, lida FROM notificacao WHERE user_id = ? ORDER BY id DESC",
		userID)
	return notifs, err
}

func (r *NotificationRepository) Remove(ctx context.Context, id, userID int64) error {
	res, err := r.db.ExecContext(ctx,
		"DELETE FROM notificacao WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *NotificationRepository) MarkRead(ctx context.Context, id, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE notificacao SET lida = true WHERE id = ? AND user_id = ?", id, userID)
	return err
}

// ChatRepository ---------------------------------------------------------------

type ChatRepository struct{ db *sqlx.DB }

func NewChatRepository(db *sqlx.DB) *ChatRepository { return &ChatRepository{db: db} }

func (r *ChatRepository) Save(ctx context.Context, senderID, receiverID int64, texto string) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO chat_message (sender_id, receiver_id, texto, data) VALUES (?, ?, ?, NOW())",
		senderID, receiverID, texto)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *ChatRepository) GetConversation(ctx context.Context, userA, userB int64, limit int) ([]models.ChatMessage, error) {
	var msgs []models.ChatMessage
	err := r.db.SelectContext(ctx, &msgs, fmt.Sprintf(`
		SELECT id, sender_id, receiver_id, texto, data
		FROM chat_message
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
		ORDER BY id DESC LIMIT %d`, limit),
		userA, userB, userB, userA)
	return msgs, err
}
