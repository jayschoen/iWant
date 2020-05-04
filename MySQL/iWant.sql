Use iWant_db;

Create table `whatsup` (
	slack_id int,
	status varchar(20),
	wants varchar(50),
	created datetime,
	target_time datetime
);

