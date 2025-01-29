CREATE TABLE IF NOT EXISTS users (
                                     id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    phone VARCHAR(16) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name VARCHAR(64),
    last_name VARCHAR(64),
    email VARCHAR(128) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (now() AT TIME ZONE 'utc') NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT (now() AT TIME ZONE 'utc') NOT NULL
    );
