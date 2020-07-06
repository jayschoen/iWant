Use iWant_db;

Create table `whatsup` (
	id int NOT NULL AUTO_INCREMENT,
	slackID int NOT NULL,
	status varchar(20),
	wants varchar(50),
	created datetime NOT NULL,
	targetTime datetime,
	PRIMARY KEY (id)
);

