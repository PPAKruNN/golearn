CREATE TABLE IF NOT EXISTS "Transfer" (
	"id" bigint GENERATED ALWAYS AS IDENTITY NOT NULL UNIQUE,
	"account_origin_id" bigint NOT NULL,
	"account_destination_id" bigint NOT NULL,
	"amount" bigint NOT NULL,
	"created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "Account" (
	"id" bigint GENERATED ALWAYS AS IDENTITY NOT NULL UNIQUE,
	"name" text NOT NULL,
	"cpf" text NOT NULL,
	"secret" text NOT NULL,
	"balance" bigint NOT NULL,
	"created_at" timestamp with time zone NOT NULL DEFAULT NOW(),
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "Auth" (
	"id" bigint NOT NULL UNIQUE,
  "token" uuid NOT NULL,
	PRIMARY KEY ("id")
);

ALTER TABLE "Transfer" ADD CONSTRAINT "Transfer_fk1" FOREIGN KEY ("account_origin_id") REFERENCES "Account"("id");

ALTER TABLE "Transfer" ADD CONSTRAINT "Transfer_fk2" FOREIGN KEY ("account_destination_id") REFERENCES "Account"("id");
