# Reservations System

A hotel reservations system.

 - Built using [Go](https://golang.org/).
 - Uses [chi router](github.com/go-chi/chi/v5)
 - Uses [alex edwards SCS session management](github.com/alexedwards/scs/v2)
 - Uses [nosurf](github.com/justinas/nosurf)


## database setup

* create a database in postgresql e.g `create database reservation_system`
* create a `database.yaml` file and populate the *database* *user* and *password*
* [install soda for migrations](https://gobuffalo.io/en/docs/db/toolbox/)
* *note*: if soda commands don't work you will need to add the following command to your .bashprofile(linux) or .zhrc file (mac).

```bash
export PATH="$HOME/go/bin:$PATH"
```

```yaml
development:
  dialect: postgres
  database: 
  user: 
  password: 
  host: 127.0.0.1
  pool: 5

test:
  url: {{envOr "TEST_DATABASE_URL" "postgres://postgres:postgres@127.0.0.1:5432/myapp_test"}}

production:
  url: {{envOr "DATABASE_URL" "postgres://postgres:postgres@127.0.0.1:5432/myapp_production"}}
```
* run `soda migrate`


## Email server

Email notifications are setup for bookings, [Mailhog](https://github.com/mailhog/MailHog) is good to use for testing. Simply install it and have the server running to view emails being sent. 

## run the application and develop
* start the application `./run.sh`

* Add the flag `-production=true` to use the template cache instead of creating a template on each render.

