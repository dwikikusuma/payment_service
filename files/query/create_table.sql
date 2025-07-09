CREATE TABLE IF NOT EXISTS status (
                                      id SERIAL PRIMARY KEY,
                                      name VARCHAR UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS payments (
                                        id BIGSERIAL PRIMARY KEY,
                                        order_id BIGINT,
                                        user_id BIGINT,
                                        external_id TEXT UNIQUE NOT NULL,
                                        amount NUMERIC,
                                        status_id INTEGER REFERENCES status(id),
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

INSERT INTO status (name) VALUES
                              ('pending'),
                              ('success'),
                              ('failed'),
                              ('refunded');
