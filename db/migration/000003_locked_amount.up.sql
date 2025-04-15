CREATE TABLE locked_amounts (
  id SERIAL PRIMARY KEY,
  sell_req_id INTEGER NOT NULL REFERENCES sell_requests(sell_req_id),
  buy_req_id UUID NOT NULL REFERENCES buy_requests(buy_req_id) ON DELETE CASCADE,
  locked_total_amount BIGINT NOT NULL,
  locked_by_card BIGINT DEFAULT 0,
  locked_by_cash BIGINT DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  is_released BOOLEAN DEFAULT false,
  released_at TIMESTAMPTZ
);

CREATE INDEX ON locked_amounts (sell_req_id);
CREATE INDEX ON locked_amounts (buy_req_id);