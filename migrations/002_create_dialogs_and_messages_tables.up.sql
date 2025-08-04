-- 002_create_dialogs_and_messages_tables.up.sql

-- Создаем таблицу для хранения сессий диалогов
CREATE TABLE IF NOT EXISTS dialogs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL DEFAULT 'New Chat',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Создаем таблицу для хранения отдельных сообщений
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    dialog_id BIGINT NOT NULL REFERENCES dialogs(id) ON DELETE CASCADE,
    role VARCHAR(10) NOT NULL, -- 'user' or 'ai'
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Создаем индексы для ускорения поиска
CREATE INDEX IF NOT EXISTS dialogs_user_id_idx ON dialogs (user_id);
CREATE INDEX IF NOT EXISTS messages_dialog_id_idx ON messages (dialog_id);