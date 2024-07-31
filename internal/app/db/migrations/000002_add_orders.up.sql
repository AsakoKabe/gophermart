CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED', 'REGISTERED');

CREATE TABLE orders
(
    id          uuid         DEFAULT gen_random_uuid() PRIMARY KEY,
    num         varchar(40) NOT NULL,
    status      order_status default 'NEW',
    accrual     float4       default 0,
    user_id     uuid references users (id),
    uploaded_at timestamptz  DEFAULT CURRENT_TIMESTAMP
);

