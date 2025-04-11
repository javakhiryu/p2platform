CREATE TABLE "sell_requests" (
  "sell_req_id" INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "sell_total_amount" bigint NOT NULL,
  "sell_money_source" VARCHAR NOT NULL CHECK (sell_money_source IN ('cash', 'card')),
  "currency_from" varchar NOT NULL,
  "currency_to" varchar NOT NULL,
  "tg_username" varchar NOT NULL,
  "sell_by_card" bool,
  "sell_amount_by_card" bigint,
  "sell_by_cash" bool,
  "sell_amount_by_cash" bigint,
  "sell_exchange_rate" bigint,
  "is_actual" bool DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "is_deleted" bool DEFAULT false,
  "comment" varchar NOT NULL
);

CREATE TABLE "buy_requests" (
  "buy_req_id" uuid UNIQUE PRIMARY KEY,
  "sell_req_id" integer NOT NULL REFERENCES "sell_requests" ("sell_req_id") ON DELETE CASCADE,
  "buy_total_amount" bigint NOT NULL,
  "tg_username" varchar NOT NULL,
  "buy_by_card" bool,
  "buy_amount_by_card" bigint,
  "buy_by_cash" bool,
  "buy_amount_by_cash" bigint,
  "close_confirm_by_seller" bool DEFAULT false,
  "close_confirm_by_buyer" bool DEFAULT false,
  "seller_confirmed_at" timestamptz,
  "buyer_confirmed_at" timestamptz,
  "is_closed" bool DEFAULT false,
  "closed_at" timestamptz,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "expires_at" timestamptz NOT NULL DEFAULT (now() + interval '1 hour')
);

CREATE INDEX ON "sell_requests" ("sell_req_id");
CREATE INDEX ON "buy_requests" ("sell_req_id");