CREATE TABLE withdrawals
(
    id           uuid        DEFAULT gen_random_uuid() PRIMARY KEY,
    num_order    varchar(40) NOT NULL,
    sum          float,
    user_id      uuid references users (id),
    processed_at timestamptz DEFAULT CURRENT_TIMESTAMP
);