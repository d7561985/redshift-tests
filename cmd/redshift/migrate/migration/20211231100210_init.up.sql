CREATE TABLE IF NOT EXISTS "balance"
(
    "accountId"      INT8   NOT NULL PRIMARY KEY,
    "balance"        float4 NOT NULL,
    "depositAllSum"  float4 NOT NULL,
    "depositCount"   INT    NOT NULL,
    "pincoinBalance" float4 NOT NULL,
    "pincoinAllSum"  float4 NOT NULL
);

CREATE TABLE IF NOT EXISTS "journal"
(
    "id"             VARBYTE                    NOT NULL,
    "transactionId"  VARBYTE                    NOT NULL,
    "accountId"      integer                     NOT NULL,
    "created_at"     TIMESTAMP WITH TIME ZONE NOT NULL,
    "balance"        FLOAT8   DEFAULT NULL,
    "pincoinBalance" FLOAT8   DEFAULT NULL,
    "change"         FLOAT4   DEFAULT NULL,
    "pincoinChange"  FLOAT4   DEFAULT NULL,
    "currency"       SMALLINT DEFAULT NULL,
    "project"        VARCHAR(10)                 NOT NULL,
    "type"           VARCHAR(20)                 NOT NULL,
    "revert"         BOOLEAN  DEFAULT NULL,
    primary key (id),
    FOREIGN KEY ("accountId") REFERENCES balance ("accountId")
) distkey(accountId);