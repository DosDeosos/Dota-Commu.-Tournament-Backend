create table players (
  id            text primary key,
  name          text not null,
  primary_roles text not null default '[]' check (json_valid(primary_roles))
) strict;

create table ratings (
  subject_id text not null references players (id),
  rater      text not null,
  scores     text not null default '{}' check (json_valid(scores)),
  positions  text not null default '{}' check (json_valid(positions)),
  traits     text not null default '[]' check (json_valid(traits)),
  primary key (subject_id, rater)
) strict;
