CREATE TABLE users
(
    id         uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    login      VARCHAR(40) UNIQUE NOT NULL,
    password   VARCHAR(40)        NOT NULL,
    accruals   float4 default 0,
    withdrawal float4 default 0
);