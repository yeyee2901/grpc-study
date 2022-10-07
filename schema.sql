create table books (
	id 		serial 			primary key,
	title 	varchar(255)	not null,
	isbn	varchar(255)	not null,
	tahun	int				not null
);
