CREATE TABLE IF NOT EXISTS organisations (
    ods             VARCHAR(10)     NOT NULL,
    tenancy         CHAR(40)        NOT NULL,
    organisation    VARCHAR(255)    NOT NULL,
    PRIMARY KEY (ods)
);