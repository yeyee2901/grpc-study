-- public.books definition

-- Drop table

-- DROP TABLE books;

CREATE TABLE books (
	id serial4 NOT NULL,
	title varchar(255) NOT NULL,
	isbn varchar(255) NOT NULL,
	tahun int4 NOT NULL,
	CONSTRAINT books_pkey PRIMARY KEY (id)
);

/* bulk insert */
INSERT INTO books 
	(title, isbn, tahun)
VALUES
	('Buku 1', '123', 2000),
	('Buku 2', '456', 2002),
	('Buku 3', '789', 2004);


-- users definition

-- Drop table

-- DROP TABLE users;
-- WARN: for some reason time zone nya ga ke generate dari DBeaver, 
-- jadi harus manual.
-- time zone perlu untuk time formatting RFC3339 (standard timestamp)
-- 'timezone' variable di define di postgre.conf (for me: 'Asia/Jakarta')
CREATE TABLE users  (
	id			BIGSERIAL					PRIMARY KEY,
	name		VARCHAR(255)				NOT NULL,
	email		VARCHAR(100)				NOT NULL DEFAULT 'your_email@email.com',
	created_at	TIMESTAMP WITH TIME ZONE	DEFAULT CURRENT_TIMESTAMP::TIMESTAMP WITH TIME ZONE 
);

/* bulk insert */
INSERT INTO users
	(name)
VALUES
	('user1'),
	('user2'),
	('user3');
