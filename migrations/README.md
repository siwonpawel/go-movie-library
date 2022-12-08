brew install golang-migrate

Creation of new migrations:

migrate create -seq -ext=.sql -dir ./migrations add_movies_check_constraints 

Migration:

migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up