-- doc-registry-go database initialize script

CREATE TYPE entity_state_t AS ENUM ('ENABLED', 'DISABLED', 'DELETED');

CREATE TABLE tbl_organization
(
    id         serial PRIMARY KEY,
    name       VARCHAR(300)                NOT NULL,
    label      VARCHAR(512)                NOT NULL,
    url        VARCHAR(900)                NOT NULL,
    public_key TEXT                        NOT NULL,
    state      entity_state_t              NOT NULL,
    create_ts  TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    update_ts  TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    version    INT                         NOT NULL
);

CREATE UNIQUE INDEX uq_organization_name ON tbl_organization (name)
    WHERE state = 'ENABLED'::entity_state_t;

CREATE UNIQUE INDEX uq_organization_url ON tbl_organization (url)
    WHERE state = 'ENABLED'::entity_state_t;

CREATE UNIQUE INDEX uq_organization_public_key ON tbl_organization (public_key)
    WHERE state = 'ENABLED'::entity_state_t;
