PRAGMA foreign_key = ON;

CREATE TABLE user (
    id           INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')           PRIMARY KEY AUTOINCREMENT,
    user_id      TEXT NOT NULL CHECK(TYPEOF(user_id) = 'text'),
    user_type    TEXT NOT NULL CHECK(TYPEOF(user_type) = 'text'),
    insert_time  INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer')  DEFAULT(CAST(strftime('%s', 'now') AS INTEGER)),
    UNIQUE(user_type COLLATE NOCASE, user_id COLLATE NOCASE)
);

CREATE TABLE field_history (
    id           INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')           PRIMARY KEY AUTOINCREMENT,
    user_id      INTEGER NOT NULL CHECK(TYPEOF(user_id) = 'integer')      REFERENCES user(id),
    field_type   TEXT NOT NULL CHECK(TYPEOF(field_type) = 'text'),
    field_text   TEXT NOT NULL CHECK(TYPEOF(field_text) = 'text'),
    insert_time  INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer')  DEFAULT(CAST(strftime('%s', 'now') AS INTEGER)),
    UNIQUE(user_id, field_type COLLATE NOCASE, field_text COLLATE NOCASE)
)