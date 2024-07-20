CREATE TABLE orders
(
    id      uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    num     varchar(40) NOT NULL,
    user_id uuid references users (id),
    uploaded_at timestamptz DEFAULT CURRENT_TIMESTAMP
);