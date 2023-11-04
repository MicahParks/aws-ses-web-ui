CREATE DATABASE aswu;
\connect aswu;

CREATE SCHEMA compose;
CREATE TABLE compose.message
(
    id            BIGSERIAL PRIMARY KEY,
    uuid          UUID      NOT NULL UNIQUE,
    from_addr     TEXT      NOT NULL,
    recipient_to  TEXT[]    NOT NULL,
    recipient_cc  TEXT[]    NOT NULL,
    recipient_bcc TEXT[]    NOT NULL,
    subject       TEXT      NOT NULL,
    body          TEXT      NOT NULL,
    ses_accepted  BOOLEAN   NOT NULL DEFAULT FALSE,
    created       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON compose.message (id);
CREATE INDEX ON compose.message (uuid);
CREATE INDEX ON compose.message (from_addr);
CREATE INDEX ON compose.message USING GIN (recipient_to);
CREATE INDEX ON compose.message USING GIN (recipient_cc);
CREATE INDEX ON compose.message USING GIN (recipient_bcc);
CREATE INDEX ON compose.message (subject);
CREATE INDEX ON compose.message (created);
