CREATE TABLE business
(
    business_id SERIAL PRIMARY KEY,
    name        TEXT NOT NULL
);

CREATE INDEX business_name_index ON business (name);

CREATE TABLE opening_hours
(
    opening_hours_id SERIAL PRIMARY KEY,
    business_id      INTEGER NOT NULL,
    day              INTEGER NOT NULL,
    opens            TIME    NOT NULL,
    closes           TIME    NOT NULL,
    FOREIGN KEY (business_id) REFERENCES business (business_id) ON UPDATE CASCADE ON DELETE CASCADE
);
