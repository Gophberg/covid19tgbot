create table nums (
	id serial primary key,
	timedate timestamp,
	cases integer,
	deaths integer,
	recovered integer,
	isolation integer,
	isolationathome integer,
	observation integer,
	quarantineathospital integer,
	unquarantined integer,
    Peopleathospital integer,
    Ambulanced integer
);
create table users (
	id serial primary key,
	tgid integer,
	addedat timestamp,
	actived boolean,
	subscases boolean,
	subsdeaths boolean,
	subsrecovered boolean,
	subsisolation boolean,
	subsisolationathome boolean,
	subsobservation boolean,
	subsquarantineathospital boolean,
	subsunquarantined boolean,
    subspeopleathospital boolean,
    subsambulanced boolean
);
