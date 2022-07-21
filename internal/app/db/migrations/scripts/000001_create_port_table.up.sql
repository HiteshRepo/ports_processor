CREATE TABLE ports
(
    name        VARCHAR(100) NOT NULL,
    city        VARCHAR(100) NOT NULL,
    country     VARCHAR(100) NOT NULL,
    alias       text[],
    regions     text[],
    coordinates float8[]     NOT NULL,
    province    VARCHAR(100) NOT NULL,
    timezone    VARCHAR(100) NOT NULL,
    unlocs      text[],
    code        VARCHAR(10)
        CONSTRAINT ports_pkey PRIMARY KEY
);