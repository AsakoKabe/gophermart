CREATE TABLE orders
(
    id      uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    num     integer NOT NULL,
    user_id uuid references users (id)
);