# Overview
The scope of this service is to insert transaction (credit/debit) and calculate user's balance

## Usage
1. Rename `config.example.yaml` to `config.local.yaml`, and fill the secret value such as postgres user, password, etc

2. Use these sql scripts to create table `account` & `account_transaction` 
```
   create table account (
	id serial primary key,
	user_id bigint not null,
	balance float8 default 0,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz DEFAULT NULL
   );

   create type transaction_type as enum('debit','credit');

   create table account_transaction (
       id bigserial primary key,
       account_id bigint not null,
       transaction_type transaction_type not null,
       remark text not null,
       amount float8 not null,
       created_at timestamptz NOT NULL DEFAULT now(),
       constraint fk_account_transaction foreign key(account_id) references public.account(id)
   );
```

3. Run your application using the command in the terminal:
   `make run`

## Libraries
echo - https://echo.labstack.com/

sqlx - https://github.com/jmoiron/sqlx
