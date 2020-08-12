Use iWant_db;

Create table `whatsup` (
	id int NOT NULL AUTO_INCREMENT,
	slackName varchar(20) NOT NULL,
	urgency varchar(30),
	wants varchar(50),
	created datetime NOT NULL,
	appointmentTime datetime,
	PRIMARY KEY (id)
);

