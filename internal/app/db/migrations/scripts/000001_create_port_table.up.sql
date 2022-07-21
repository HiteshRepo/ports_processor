CREATE TABLE ports
(
    name        VARCHAR(100) CONSTRAINT ports_pkey PRIMARY KEY,
    city        VARCHAR(100) NOT NULL,
    country     VARCHAR(100) NOT NULL,
    alias       VARCHAR(100),
    regions     VARCHAR(100),
    coordinates VARCHAR(100) NOT NULL,
    province    VARCHAR(100) NOT NULL,
    timezone    VARCHAR(100) NOT NULL,
    unlocs      VARCHAR(100),
    code        VARCHAR(10)
);