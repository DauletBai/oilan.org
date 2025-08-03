-- 001_create_users_table.up.sql

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    provider VARCHAR(50) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Создаем уникальный индекс, чтобы один и тот же пользователь
-- не мог зарегистрироваться дважды через одного и того же провайдера.
CREATE UNIQUE INDEX IF NOT EXISTS users_provider_provider_id_idx ON users (provider, provider_id);