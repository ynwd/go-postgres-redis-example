create table tasks (
  id serial primary key,
  description text not null
);

INSERT INTO "tasks" ("id", "description") VALUES
(1,	'learn go'),
(2,	'learn python');

