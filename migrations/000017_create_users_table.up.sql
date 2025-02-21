CREATE TABLE IF NOT EXISTS "users_profiles" (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    phone VARCHAR(16) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name VARCHAR(64),
    last_name VARCHAR(64),
    email VARCHAR(128) UNIQUE,
    refresh_token VARCHAR(64),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (now() AT TIME ZONE 'utc') NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT (now() AT TIME ZONE 'utc') NOT NULL
    );
CREATE TABLE failed_logins (
                               phone VARCHAR(20) PRIMARY KEY,
                               attempts INT DEFAULT 0,
                               blocked_until TIMESTAMP
);

CREATE TABLE users_rights (
                                id SERIAL PRIMARY KEY,            -- ID пользователя, автоматически увеличивается
                                name VARCHAR(255) NOT NULL,       -- Имя пользователя
                                context VARCHAR(255),             -- Контекст пользователя (может быть пустым)
                                rights JSONB                      -- Права пользователя в формате JSONB
);

CREATE TABLE IF NOT EXISTS "users" (
                                     id UUID PRIMARY KEY,
                                     phone VARCHAR(255) UNIQUE NOT NULL,
                                     password_hash VARCHAR(255) NOT NULL,
                                     rights JSONB NOT NULL,
                                     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);



