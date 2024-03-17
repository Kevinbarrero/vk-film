
CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL CHECK (LENGTH(name) > 0),
    description VARCHAR(1000) NOT NULL,
    release_date DATE NOT NULL,
    rating DECIMAL(3, 1) NOT NULL CHECK (rating >= 0 AND rating <= 10)
);


CREATE TABLE actors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    gender VARCHAR(6) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    birthday DATE NOT NULL
);
