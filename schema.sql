CREATE TABLE Client (
    id       SERIAL PRIMARY KEY,
    login    VARCHAR(30) NOT NULL,
    password VARCHAR(30) NOT NULL,
    name     VARCHAR(30) NOT NULL,
    type     VARCHAR(30) NOT NULL
);

CREATE TABLE Message (
    id        SERIAL PRIMARY KEY,
    clientId  INTEGER,
    text      VARCHAR(300) NOT NULL,
    FOREIGN KEY (clientId) REFERENCES Client (id)
);

INSERT INTO Client (login, password, name, type)
VALUES 
('admin', '111', 'Админ', 'ADMIN'),
('anton', '111', 'Антон', 'CLIENT')