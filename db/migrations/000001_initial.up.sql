CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "consumers" (
    "id" uuid PRIMARY KEY,
    "api_key" uuid NOT NULL,
    "name" varchar(255) NOT NULL,
    "created_at" timestamp NOT NULL
);

CREATE UNIQUE INDEX api_key_idx ON consumers(api_key);

CREATE TYPE "transaction_status" AS ENUM ('created', 'processing', 'processed', 'canceled');
CREATE TYPE "transaction_type" AS ENUM ('withdraw', 'replenish');

CREATE TABLE "transactions" (
    "id" uuid PRIMARY KEY,
    "amount" decimal(12,2) NOT NULL,
    "status" transaction_status NOT NULL,
    "type" transaction_type NOT NULL,
    "consumer_id" uuid NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp
);

ALTER TABLE "transactions" ADD CONSTRAINT "fk_consumer_id" FOREIGN KEY ("consumer_id") REFERENCES "consumers"("id");