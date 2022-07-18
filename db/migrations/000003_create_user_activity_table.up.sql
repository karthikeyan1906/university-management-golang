CREATE TABLE IF NOT EXISTS user_activity (
    id bigint NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    studentId bigint REFERENCES students(rollnumber),
    signin TIMESTAMP,
    signout TIMESTAMP
);