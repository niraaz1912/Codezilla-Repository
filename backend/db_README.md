#db.db README

## TABLES
### users
- sessionID [string, PRIMARY KEY]
- username [string]
- role [string]
- passhash [string]

### location
- sessionID [string, FOREIGN KEY]
- longitude [floating point]
- latitude [floating point]

## QUERIES

ALL USERS
`SELECT * FROM users;`

LOCATION DATA FOR ALL SESSIONS
`SELECT * FROM location;`

ALL USERS AND THEIR LOCATIONS
`SELECT * FROM users INNER JOIN location on location.sessionID = users.sessionID;`