-- Migration: Create FreakDream schema
-- Run once against a fresh database.
-- Supports MySQL 8.0+

CREATE DATABASE IF NOT EXISTS freakdream CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE freakdream;

-- Users
CREATE TABLE IF NOT EXISTS usuario (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    nome       VARCHAR(100)  NOT NULL,
    sobre      VARCHAR(200)  NOT NULL DEFAULT '',
    email      VARCHAR(255)  NOT NULL UNIQUE,
    senha      VARCHAR(255)  NOT NULL,  -- bcrypt hash
    created_at DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Login history
CREATE TABLE IF NOT EXISTS login_historico (
    id      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    data    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ip      VARCHAR(45)     NOT NULL,
    FOREIGN KEY (user_id) REFERENCES usuario(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Photos (stored filenames)
CREATE TABLE IF NOT EXISTS fotos (
    foto_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    foto    VARCHAR(255) NOT NULL,
    texto   TEXT         NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Profile photo link (latest wins via ORDER BY foto_id DESC)
CREATE TABLE IF NOT EXISTS perfil_foto (
    foto_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (user_id, foto_id),
    FOREIGN KEY (foto_id) REFERENCES fotos(foto_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES usuario(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Posts (tipo: 1=text, 2=photo, 3=share/repost)
CREATE TABLE IF NOT EXISTS post (
    post_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    tipo    TINYINT UNSIGNED NOT NULL DEFAULT 1,
    texto   TEXT             NOT NULL,
    data    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_data (data),
    FOREIGN KEY (user_id) REFERENCES usuario(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Follow relationships
CREATE TABLE IF NOT EXISTS seguir (
    user_id     BIGINT UNSIGNED NOT NULL,  -- who is followed
    seguidor_id BIGINT UNSIGNED NOT NULL,  -- who follows
    PRIMARY KEY (user_id, seguidor_id),
    FOREIGN KEY (user_id)     REFERENCES usuario(id) ON DELETE CASCADE,
    FOREIGN KEY (seguidor_id) REFERENCES usuario(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Friends / friend requests
CREATE TABLE IF NOT EXISTS amigo (
    user_id   BIGINT UNSIGNED NOT NULL,
    amigo_id  BIGINT UNSIGNED NOT NULL,
    status    ENUM('pending','accepted') NOT NULL DEFAULT 'pending',
    criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, amigo_id),
    FOREIGN KEY (user_id)  REFERENCES usuario(id) ON DELETE CASCADE,
    FOREIGN KEY (amigo_id) REFERENCES usuario(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Notifications
CREATE TABLE IF NOT EXISTS notificacao (
    id      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    tipo    VARCHAR(50)     NOT NULL,
    texto   TEXT            NOT NULL,
    data    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    lida    BOOLEAN         NOT NULL DEFAULT FALSE,
    INDEX idx_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES usuario(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Chat messages
CREATE TABLE IF NOT EXISTS chat_message (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    sender_id   BIGINT UNSIGNED NOT NULL,
    receiver_id BIGINT UNSIGNED NOT NULL,
    texto       TEXT            NOT NULL,
    data        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_conversation (sender_id, receiver_id),
    FOREIGN KEY (sender_id)   REFERENCES usuario(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES usuario(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
