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

CREATE TABLE category (
    id             INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')           PRIMARY KEY AUTOINCREMENT,
    category_name  TEXT NOT NULL CHECK(TYPEOF(category_name) = 'text'),
    related_id     TEXT NOT NULL CHECK(TYPEOF(related_id) = 'text')         DEFAULT('28'),
    insert_time    INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer')  DEFAULT(CAST(strftime('%s', 'now') AS INTEGER)),
    UNIQUE(category_name COLLATE NOCASE)
);

INSERT INTO category (category_name) VALUES("Garry's Mod");

CREATE TABLE media_recording (
    id           INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')                                   PRIMARY KEY AUTOINCREMENT,
    file_name    TEXT NOT NULL CHECK(TYPEOF(file_name) = 'text'),
    extension    TEXT NOT NULL CHECK(TYPEOF(extension) = 'text' AND extension IN ('mkv', 'mp4'))  DEFAULT('mkv'),
    directory    TEXT NOT NULL CHECK(TYPEOF(directory) = 'text'),
    start_time   INTEGER NOT NULL CHECK(TYPEOF(start_time) = 'integer')                           DEFAULT(CAST(strftime('%s', 'now') AS INTEGER)),
    end_time     INTEGER,
    uploaded     INTEGER NOT NULL CHECK(TYPEOF(uploaded) = 'integer' AND uploaded IN (0, 1))      DEFAULT(0),
    insert_time  INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer')                          DEFAULT(CAST(strftime('%s', 'now') AS INTEGER))
);

CREATE TABLE subtitles (
    id           INTEGER NOT NULL CHECK(TYPEOF(id) = 'integer')           PRIMARY KEY AUTOINCREMENT,
    subtitle     TEXT NOT NULL CHECK(TYPEOF(subtitle) = 'text'),
    duration     REAL NOT NULL CHECK(TYPEOF(duration) = 'real'),
    insert_time  INTEGER NOT NULL CHECK(TYPEOF(insert_time) = 'integer')  DEFAULT(CAST(strftime('%s', 'now') AS INTEGER))
);