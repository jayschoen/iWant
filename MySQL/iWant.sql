Use iWant_db;

Create table `whatsup` (
	id int NOT NULL AUTO_INCREMENT,
	slackName varchar(20) NOT NULL,
	status varchar(20),
	wants varchar(50),
	created datetime NOT NULL,
	targetTime datetime,
	PRIMARY KEY (id)
);

