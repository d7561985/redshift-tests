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
    "id"             VARBYTE                  NOT NULL,
    "transactionId"  VARBYTE                  NOT NULL,
    "accountId"      integer                  NOT NULL,
    "created_at"     TIMESTAMP WITH TIME ZONE NOT NULL,
    "balance"        FLOAT8   DEFAULT NULL,
    "pincoinBalance" FLOAT8   DEFAULT NULL,
    "change"         FLOAT4   DEFAULT NULL,
    "pincoinChange"  FLOAT4   DEFAULT NULL,
    "currency"       SMALLINT DEFAULT NULL,
    "project"        VARCHAR(10)              NOT NULL,
    "type"           VARCHAR(20)              NOT NULL,
    "revert"         BOOLEAN  DEFAULT NULL,
    primary key (id),
    FOREIGN KEY ("accountId") REFERENCES balance ("accountId")
) distkey (accountId);

CREATE TABLE IF NOT EXISTS "players"
(
    "id"                INT8                     NOT NULL PRIMARY KEY,
    "guid"              varchar(40),
    "license"           varchar(45)              NOT NULL,
    "playerId"          INT8                     NOT NULL,
    "clickID"           VARCHAR(40)              NOT NULL,
    "registerDate"      TIMESTAMP WITH TIME ZONE NOT NULL,
    "language"          VARCHAR(25)              NOT NULL,
    "email"             varchar(50)              NOT NULL,
    "isEmailVerify"     boolean,
    "phone"             VARCHAR(40),
    "isPhoneVerify"     boolean,
    "isMultiAccount"    boolean,
    "birthday"          varchar(20),
    "accountVerifyTime" TIMESTAMP WITH TIME ZONE,
    "lastLoginTime"     TIMESTAMP WITH TIME ZONE,
    "country"           VARCHAR(45)              NOT NULL,
    "city"              VARCHAR(25)              NOT NULL,
    "currency"          varchar(5)               NOT NULL,
    "sex"               varchar(7),
    "isTest"            boolean,
    "isBot"             boolean,
    "project"           VARCHAR(10)              NOT NULL,
    "activateStatus"    VARCHAR(30),
    "depositStatus"     VARCHAR(30),
    "smsStatus"         VARCHAR(30),
    "domain"            VARCHAR(20)              NOT NULL,
    "webview"           boolean,
    "ipAddress"         VARCHAR(15)              NOT NULL,
    "userAgent"         VARCHAR(100)             NOT NULL,
    "createUnixNano"    int8                     NOT NULL,
    "updateUnixNano"    int8                     NOT NULL
);

CREATE TABLE IF NOT EXISTS "cb"
(
    "id"             varbyte                  NOT NULL,
    "license"        VARCHAR(45)              NOT NULL,
    "playerId"       integer                  NOT NULL,
    "gameName"       VARCHAR(30)              NOT NULL,
    "gameType"       VARCHAR(10)              NOT NULL,
    "gameID"         int8                     NOT NULL,
    "bonusID"        int8                     NOT NULL,
    "bet"            float8                   NOT NULL,
    "winLose"        float8                   NOT NULL,
    "purse"          varchar(20)              NOT NULL,
    "currencyCode"   VARCHAR(5),
    "gameProvider"   VARCHAR(15),
    "gameRoundID"    varchar(40),
    "tranID"         varchar(40)              NOT NULL,
    "date"           TIMESTAMP WITH TIME ZONE NOT NULL,
    "createUnixNano" int8                     NOT NULL,
    "updateUnixNano" int8                     NOT NULL,
    "rollback"       boolean,
    "status"         VARCHAR(30)              NOT NULL,
    "error"          VARCHAR(100),
    "hall"           VARCHAR(45),
    "sstm"         VARCHAR(20),
    "betInfo"        VARCHAR(30),
    "agent"          int8                     NOT NULL,
    "domain"         VARCHAR(20)              NOT NULL,
    "webview"        boolean,
    "isTournament"   boolean,
    primary key (id),
    FOREIGN KEY ("playerId") REFERENCES players ("id")
) distkey (playerId);

