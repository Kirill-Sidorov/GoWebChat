DROP TABLE IF EXISTS Client;
CREATE TABLE Client (
    id       SERIAL PRIMARY KEY,
    login    VARCHAR(30) NOT NULL,
    password VARCHAR(70) NOT NULL,
    name     VARCHAR(30) NOT NULL,
    type     VARCHAR(30) NOT NULL
);

INSERT INTO Client (login, password, name, type)
VALUES 
('admin', '$2a$14$K7y0tG1YP3IRpTZfHFfBQuxHymYKp5xm8zuhYoMic0wQW85OCq7Z2', 'Admin', 'ADMIN'),
('anton', '$2a$14$H4n4LQJUO2lSO.35R8zRHu97FE4hIBtX8xISzbQ95ov.iMpRNd.9.', 'Anton', 'CLIENT')