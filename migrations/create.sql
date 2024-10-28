
CREATE TABLE IF NOT EXISTS ads (
    "id"           SERIAL PRIMARY KEY,
    "name"         CHARACTER VARYING(150) NOT NULL,
    "description"  TEXT,
    "price"        NUMERIC(10, 2) NOT NULL,
    "is_active"    BOOLEAN NOT NULL DEFAULT false,
    "created_at"   TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
    UNIQUE("name", "is_active")
);

CREATE INDEX ads_name_index ON ads ("name");
