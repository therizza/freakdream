package models

import "time"

type User struct {
	ID        int64     `db:"id"         json:"id"`
	Nome      string    `db:"nome"       json:"nome"`
	Sobre     string    `db:"sobre"      json:"sobre"`
	Email     string    `db:"email"      json:"email"`
	Senha     string    `db:"senha"      json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type UserPublic struct {
	ID    int64  `json:"id"`
	Nome  string `json:"nome"`
	Sobre string `json:"sobre"`
	Foto  string `json:"foto"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

type RegisterRequest struct {
	Nome  string `json:"nome"`
	Sobre string `json:"sobre"`
	Email string `json:"email"`
	Senha string `json:"senha"`
}

type PostType int

const (
	PostTypeText  PostType = 1
	PostTypePhoto PostType = 2
	PostTypeShare PostType = 3
)

type Post struct {
	PostID    int64     `db:"post_id"    json:"post_id"`
	UserID    int64     `db:"user_id"    json:"user_id"`
	Tipo      PostType  `db:"tipo"       json:"tipo"`
	Texto     string    `db:"texto"      json:"texto"`
	Data      time.Time `db:"data"       json:"data"`
	User      *UserPublic `db:"-"        json:"user,omitempty"`
}

type Photo struct {
	FotoID int64  `db:"foto_id" json:"foto_id"`
	Foto   string `db:"foto"    json:"foto"`
	Texto  string `db:"texto"   json:"texto"`
}

type ProfilePhoto struct {
	FotoID int64 `db:"foto_id" json:"foto_id"`
	UserID int64 `db:"user_id" json:"user_id"`
}

type Follow struct {
	UserID      int64 `db:"user_id"      json:"user_id"`
	SeguidorID  int64 `db:"seguidor_id"  json:"seguidor_id"`
}

type Friend struct {
	UserID  int64     `db:"user_id"   json:"user_id"`
	AmigoID int64     `db:"amigo_id"  json:"amigo_id"`
	Status  string    `db:"status"    json:"status"` // pending, accepted
	CriadoEm time.Time `db:"criado_em" json:"criado_em"`
}

type Notification struct {
	ID     int64     `db:"id"      json:"id"`
	UserID int64     `db:"user_id" json:"user_id"`
	Tipo   string    `db:"tipo"    json:"tipo"`
	Texto  string    `db:"texto"   json:"texto"`
	Data   time.Time `db:"data"    json:"data"`
	Lida   bool      `db:"lida"    json:"lida"`
}

type ChatMessage struct {
	ID         int64     `db:"id"          json:"id"`
	SenderID   int64     `db:"sender_id"   json:"sender_id"`
	ReceiverID int64     `db:"receiver_id" json:"receiver_id"`
	Texto      string    `db:"texto"       json:"texto"`
	Data       time.Time `db:"data"        json:"data"`
}

type WSMessage struct {
	Type       string `json:"type"`
	ReceiverID int64  `json:"receiver_id,omitempty"`
	Texto      string `json:"texto,omitempty"`
	SenderID   int64  `json:"sender_id,omitempty"`
	Data       string `json:"data,omitempty"`
}
