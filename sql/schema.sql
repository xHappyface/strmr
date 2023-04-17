PRAGMA foreign_key = ON;

CREATE TABLE user (
    id           INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')           PRIMARY KEY AUTOINCREMENT,
    user_id      TEXT NOT NULL CHECK(TYPEOF(user_id) = 'text'),
    user_type    TEXT NOT NULL CHECK(TYPEOF(user_type) = 'text'),
    insert_time  INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer')  DEFAULT(CAST(strftime('%s', 'now') AS INTEGER)),
    UNIQUE(user_type COLLATE NOCASE, user_id COLLATE NOCASE)
);

CREATE TABLE stream (
    id           INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')           PRIMARY KEY AUTOINCREMENT,
    start_time   INTEGER NOT NULL CHECK(TYPEOF(start_time) = 'integer')   DEFAULT(CAST(strftime('%s', 'now') AS INTEGER)),
    end_time     INTEGER,
    insert_time  INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer')  DEFAULT(CAST(strftime('%s', 'now') AS INTEGER))
);

CREATE TABLE metadata (
    id              INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')         PRIMARY KEY AUTOINCREMENT,
    metadata_key    TEXT NOT NULL CHECK(TYPEOF(metadata_key) = 'text'),
    metadata_value  TEXT NOT NULL CHECK(TYPEOF(metadata_value) = 'text'),
    insert_time     INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer')  DEFAULT(CAST(strftime('%s', 'now') AS INTEGER))
);

CREATE TABLE media_recording (
    id           INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')           PRIMARY KEY AUTOINCREMENT,
    file_name    TEXT NOT NULL CHECK(TYPEOF(file_name) = 'text'),
    directory    TEXT NOT NULL CHECK(TYPEOF(directory) = 'text'),
    start_time   INTEGER NOT NULL CHECK(TYPEOF(start_time) = 'integer')   DEFAULT(CAST(strftime('%s', 'now') AS INTEGER)),
    end_time     INTEGER,
    insert_time  INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer')  DEFAULT(CAST(strftime('%s', 'now') AS INTEGER))
);