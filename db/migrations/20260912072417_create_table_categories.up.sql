create table categories (
    id varchar(100) primary key,
    name varchar(255) not null,
    parent_id varchar(100),
    foreign key (parent_id) references categories (id)
) engine = InnoDB;
