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

CREATE TABLE IF NOT EXISTS payment_anomalies (
                                                 id SERIAL PRIMARY KEY,
                                                 order_id BIGINT,
                                                 external_id TEXT,
                                                 anomaly_type INTEGER,
                                                 notes TEXT,
                                                 status INTEGER,
                                                 create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                 update_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS failed_events (
                                             id SERIAL PRIMARY KEY,
                                             order_id BIGINT,
                                             external_id TEXT,
                                             failed_type INTEGER,
                                             notes TEXT,
                                             status INTEGER,
                                             create_time TIMESTAMP,
                                             update_time TIMESTAMP
);


INSERT INTO status (name) VALUES
                              ('pending'),
                              ('success'),
                              ('failed'),
                              ('refunded');
