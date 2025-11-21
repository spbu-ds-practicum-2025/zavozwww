CREATE TABLE IF NOT EXISTS user_profiles
(
    user_id    UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name  TEXT NOT NULL,
    age        INT NOT NULL,
    info       TEXT,
    city       TEXT,

    CONSTRAINT fk_user_profile
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);