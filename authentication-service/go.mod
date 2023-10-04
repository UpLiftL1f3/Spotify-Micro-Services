module github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service

go 1.20

replace github.com/UpLiftL1f3/Spotify-Micro-Services/shared => ../shared

require (
	github.com/UpLiftL1f3/Spotify-Micro-Services/shared v0.0.0-00010101000000-000000000000
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/google/uuid v1.3.1
	github.com/jackc/pgconn v1.14.0
	github.com/jackc/pgx/v4 v4.18.1
	github.com/lib/pq v1.10.2
	golang.org/x/crypto v0.13.0
)

require (
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/text v0.13.0 // indirect
)
