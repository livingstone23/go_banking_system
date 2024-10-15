CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar(255) NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar(3) NOT NULL, -- Asumiendo que las monedas son de 3 caracteres (ej. USD, EUR)
  "created_at" timestamp with time zone NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamp with time zone NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamp with time zone NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "entries" 
ADD CONSTRAINT fk_account 
FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") 
ON DELETE CASCADE;

ALTER TABLE "transfers" 
ADD CONSTRAINT fk_from_account 
FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") 
ON DELETE CASCADE;

ALTER TABLE "transfers" 
ADD CONSTRAINT fk_to_account 
FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id") 
ON DELETE CASCADE;