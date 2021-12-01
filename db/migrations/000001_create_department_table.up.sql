CREATE SEQUENCE department_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS department (
    id bigint NOT NULL DEFAULT nextval('department_id_seq'::regclass),
    name text
);