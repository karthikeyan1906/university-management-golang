CREATE TABLE IF NOT EXISTS attendance(
    id bigint NOT NULL,
    loginTime TIMESTAMP,
    CONSTRAINT fk_id
        FOREIGN KEY (id)
            REFERENCES student(roll)
);