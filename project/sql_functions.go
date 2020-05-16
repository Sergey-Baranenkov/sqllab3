package main
var dropCreateFunctions = `
create extension if not exists dblink;
create or replace function create_database() returns void as $$
    begin
        if not exists (select from pg_database where datname = 'lab4_db') then
            perform dblink_exec(
                'dbname=' || current_database(),
                'create database lab4_db with owner me'
          );
        end if;
    end;
$$ language plpgsql;

create or replace function drop_database() returns void as $$
    begin
        perform dblink_exec(
            'dbname=' || current_database(),
            'drop database if exists lab4_db;'
        );
    end;
$$ language plpgsql;
`
var sqlFunctions = `
create or replace function create_tables() returns void as $$
    begin
        create table users (
            user_id bigserial primary key not null,
            nickname text not null,
            registered_at timestamptz default current_timestamp
        ); create index if not exists nickname on users (nickname);


        create table comments (
            auth_id bigserial references users (user_id) on delete cascade on update cascade,
            comment_id bigserial primary key not null,
            message_text text not null,

            creation_time      timestamptz default current_timestamp,
            modification_time   timestamptz
        ); create index if not exists comments_auth_id_idx on comments (auth_id);


        create or replace function update__comment__modification_time() returns trigger as $u$
            begin
                new.modification_time = current_timestamp;
                return new;
            end;
        $u$ language plpgsql;

        drop trigger if exists update_comment__trigger on comments;
        create trigger update_comment__trigger
            before update on comments
        for row execute procedure update__comment__modification_time();
            end;
$$ language plpgsql;

create or replace function get_all_users() returns json as $$
	begin
		return (select json_agg(users) from users);
	end;
$$ language plpgsql;

create or replace function get_all_comments() returns json as $$
	begin
		return (select json_agg(comments) from comments);
	end;
$$ language plpgsql;

create or replace function truncate_users() returns void as $$
	begin
		truncate users cascade;
	end;
$$ language plpgsql;

create or replace function truncate_comments() returns void as $$
	begin
		truncate comments;
	end;
$$ language plpgsql;

create or replace function insert_user(_nickname text) returns json as $$
	declare output json;
begin
	insert into users (nickname) values (_nickname) returning json_build_object('user_id', user_id, 'nickname',nickname, 'registered_at',registered_at) into output;
	return output;
	end;
$$ language plpgsql;


create or replace function insert_comment(_auth_id bigint, _message_text text) returns json as $$
	declare output json;
begin
	insert into comments (auth_id, message_text)  values (_auth_id, _message_text) returning json_build_object(
			'auth_id', auth_id, 
            'comment_id', comment_id,
            'message_text', message_text, 
            'creation_time', creation_time,
            'modification_time', modification_time
            ) into output;
	return output;
	end;
$$ language plpgsql;

create or replace function update_user_nickname(_user_id bigint, new_nickname text) returns json as $$
	declare output json;
begin
	update users set nickname = new_nickname where user_id = _user_id returning json_build_object('user_id', user_id, 'nickname',nickname, 'registered_at',registered_at) into output;
		return output;
	end;
$$ language plpgsql;

create or replace function update_comment_text(_comment_id bigint, _message_text text) returns json as $$
	declare output json;
begin
	update comments set message_text = _message_text where comment_id = _comment_id returning json_build_object(
			'auth_id', auth_id, 
            'comment_id', comment_id,
            'message_text', message_text, 
            'creation_time', creation_time,
            'modification_time', modification_time
            ) into output;
		return output;
	end;
$$ language plpgsql;

create or replace function select_users_by_nickname(_nickname text) returns json as $$
	begin
		return (select json_agg(users) from users where nickname = _nickname);
	end;
$$ language plpgsql;

create or replace function select_comments_by_message(_message_text text) returns json as $$
	begin
		return (select json_agg(comments) from comments where message_text = _message_text);
	end;
$$ language plpgsql;


create or replace function delete_users_by_nickname(_nickname text) returns void as $$
	begin
		delete from users where nickname = _nickname;
	end;
$$ language plpgsql;

create or replace function delete_comments_by_message(_message_text text) returns void as $$
	begin
		delete from comments where message_text = _message_text;
	end;
$$ language plpgsql;


create or replace function delete_particular_user(_user_id bigint) returns void as $$
	begin
		delete from users where user_id = _user_id;
	end;
$$ language plpgsql;

create or replace function delete_particular_comment(_comment_id bigint) returns void as $$
	begin
		delete from comments where comment_id = _comment_id;
	end;
$$ language plpgsql;
`
