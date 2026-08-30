CREATE TABLE IF NOT EXISTS user_machine_logs (
    transaction_id  CHAR(26)        NOT NULL,
    cellar          CHAR(26)        NOT NULL,
    subject         CHAR(50)        NOT NULL,
    hostname        VARCHAR(63)     NOT NULL,
    last_accessed   TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (transaction_id)
);

CREATE INDEX ASYNC IF NOT EXISTS idx_user_machine_logs_cellar_subject ON user_machine_logs (cellar, subject, last_accessed);

COMMENT ON  TABLE user_machine_logs IS '
// As Amazon Aurora DSQL does not support foreign key constraints, adding
// DBML (https://dbml.dbdiagram.io/docs) hints as comments to indicate
// relationships between tables for documentation purposes.

Ref: users.subject > user_machine_logs.subject
';