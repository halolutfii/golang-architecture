delete from categories where parent_id is not null;

delete from categories where parent_id is null;