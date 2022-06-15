-- Creates the tables needed for the maker
-- Made for PostgreSQL
-- Jacobo Soffer Levy | A01028653
-- 18/05/2022
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    joined TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

ALTER TABLE users ADD COLUMN IF NOT EXISTS ts tsvector
    GENERATED ALWAYS AS (to_tsvector('spanish', name)) STORED;

CREATE INDEX IF NOT EXISTS ts_user_idx ON users USING GIN (ts);

CREATE TABLE IF NOT EXISTS collection (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    description VARCHAR(200),
    uid INT NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES users(id)
);

ALTER TABLE collection ADD COLUMN IF NOT EXISTS ts tsvector
    GENERATED ALWAYS AS (setweight(to_tsvector('spanish', coalesce(name, '')), 'A') ||
    setweight(to_tsvector('spanish', coalesce(description, '')), 'B')) STORED;

CREATE INDEX IF NOT EXISTS ts_collection_idx ON collection USING GIN (ts);

CREATE TABLE IF NOT EXISTS levels (
    id SERIAL PRIMARY KEY,
    difficulty INT NOT NULL,
    name VARCHAR(20) NOT NULL,
    description VARCHAR(200) NOT NULL,
    uid INT NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP,
    theme INT,
    FOREIGN KEY (uid) REFERENCES users(id)
);

ALTER TABLE levels ADD COLUMN IF NOT EXISTS ts tsvector
    GENERATED ALWAYS AS (setweight(to_tsvector('spanish', coalesce(name, '')), 'A') || 
    setweight(to_tsvector('spanish', coalesce(description, '')), 'B')) STORED;

CREATE INDEX IF NOT EXISTS ts_level_idx ON levels USING GIN (ts);

CREATE TABLE IF NOT EXISTS collection_levels (
    id SERIAL PRIMARY KEY,
    collection_id INT NOT NULL,
    level_id INT NOT NULL,
    FOREIGN KEY (collection_id) REFERENCES collection(id),
    FOREIGN KEY (level_id) REFERENCES levels(id)
);
-- CREATE TABLE IF NOT EXISTS level_stats (
--     id SERIAL PRIMARY KEY,
--     level_id INT NOT NULL,
--     uid INT NOT NULL,
--     time INT NOT NULL,
--     timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY (level_id) REFERENCES levels(id),
--     FOREIGN KEY (uid) REFERENCES users(id)
-- );
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(50) NOT NULL,
    level_id INT,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    uid INT,
    time INT,
    draft_id INT,
    body jsonb,
    state VARCHAR(50),
    FOREIGN KEY (uid) REFERENCES users(id),
    FOREIGN KEY (level_id) REFERENCES levels(id)
);
CREATE TABLE IF NOT EXISTS course_data (
    id SERIAL PRIMARY KEY,
    level_id INT UNIQUE NOT NULL,
    map_data jsonb NOT NULL,
    FOREIGN KEY (level_id) REFERENCES LEVELS(id)
);
CREATE TABLE IF NOT EXISTS drafts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    level_id INT UNIQUE,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP,
    theme INT DEFAULT 1,
    course_data jsonb,
    uid INT NOT NULL,
    FOREIGN KEY (level_id) REFERENCES levels(id),
    FOREIGN KEY (uid) REFERENCES users(id)
);

ALTER TABLE drafts ADD COLUMN IF NOT EXISTS car INT, ADD COLUMN IF NOT EXISTS soundtrack INT;
ALTER TABLE levels ADD COLUMN IF NOT EXISTS car INT, ADD COLUMN IF NOT EXISTS soundtrack INT;

-- TRIGGERS
-- 1. Delete drafts on level
-- creation
CREATE OR REPLACE FUNCTION on_level_create()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL
    AS 
$$ 
BEGIN
DELETE FROM drafts
WHERE level_id = NEW.id;
RETURN NEW;
END;
$$ 
;

DROP TRIGGER IF EXISTS level_insert_trigger ON levels;

CREATE TRIGGER level_insert_trigger
AFTER INSERT
ON levels
FOR EACH ROW
EXECUTE PROCEDURE on_level_create();

-- -- 2. Delete draft on course_data update
CREATE OR REPLACE FUNCTION on_course_data_update()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL 
    AS 
$$ 
BEGIN
DELETE FROM drafts
WHERE level_id = NEW.level_id;
RETURN NEW;
END;
$$ 
;

DROP TRIGGER IF EXISTS course_data_update_trigger ON course_data;

CREATE TRIGGER course_data_update_trigger
    AFTER UPDATE
    ON course_data
    FOR EACH ROW
    EXECUTE PROCEDURE on_course_data_update();
-- -- 3. Set the updated column on draft update
CREATE OR REPLACE FUNCTION on_draft_update()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL
    AS
$$
BEGIN
NEW.updated := CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$
;

DROP TRIGGER IF EXISTS draft_update_trigger ON drafts;

CREATE TRIGGER draft_update_trigger
AFTER UPDATE
ON drafts
FOR EACH ROW
EXECUTE PROCEDURE on_draft_update();

-- == Stored Procedures ==
-- Note that although Postgres does support
-- stored procedures, we're using functions instead,
-- as we don't need any transactions

-- 1. Query GIN Index
-- sample usage SELECT * FROM query_gin(null::levels, 'update');
CREATE OR REPLACE FUNCTION query_gin(_rowtype anyelement, q TEXT)
    RETURNS SETOF anyelement
    LANGUAGE PLPGSQL
    AS
$$
DECLARE
    query_str TEXT;
    query_ts tsquery;
BEGIN
    SELECT REPLACE(q, ' ', ' <2> ') INTO query_str;
    EXECUTE format('SELECT to_tsquery(%L, %L)', 'spanish', query_str) INTO query_ts;
    RETURN QUERY EXECUTE format('SELECT * FROM %s WHERE ts @@ %L ORDER BY ts_rank(ts, %L) DESC', pg_typeof(_rowtype), query_ts, query_ts);
    --RETURN QUERY EXECUTE format('SELECT * FROM %s WHERE ts @@ to_tsquery(%L, %L) ORDER BY ts_rank(ts, to_tsquery(%L, %L)) DESC', pg_typeof(_rowtype), 'spanish', REPLACE(q, ' ', ' <2> '), 'spanish', REPLACE(q, ' ', ' <2> '));
END
$$;

-- == Views ==

-- 1. Leaderboard
CREATE OR REPLACE VIEW leaderboard AS
    SELECT e.level_id, e.uid, e.time, e.timestamp, u.name FROM events e
        INNER JOIN users u ON e.uid = u.id
        WHERE e.event_type = 'game_finish'
        ORDER BY e.time; 

-- 2. Most popular levels
CREATE OR REPLACE VIEW trending_levels AS
    SELECT DISTINCT count(e.id) OVER (
        PARTITION BY e.level_id
    ) plays, l.* FROM events e
    INNER JOIN levels l ON e.level_id = l.id
    WHERE event_type = 'game_start'
    --GROUP BY e.level_id
    ORDER BY plays DESC;

-- 3. Most popular collections
CREATE OR REPLACE VIEW trending_collections AS
    SELECT DISTINCT c.*,
        sum(tl.plays) OVER (
            PARTITION BY c.id
        ) collection_plays FROM trending_levels tl
        RIGHT JOIN collection_levels cl ON tl.id = cl.level_id
        INNER JOIN collection c ON cl.collection_id = c.id
        --GROUP BY c.id
        ORDER BY collection_plays DESC;
