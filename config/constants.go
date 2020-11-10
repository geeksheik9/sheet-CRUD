package config

const (
	port                = "PORT"
	characterDatabase   = "CHARACTER_DATABASE"
	characterCollection = "CHARACTER_COLLECTION"
	characterArchive    = "CHARACTER_ARCHIVE"
	logLevel            = "LOG_LEVEL"
	jwtTokenDecoder     = "JWT_TOKEN_DECODER"
)

const (
	defaultPort                = "3000"
	defaultCharacterDatabase   = "characters"
	defaultCharacterCollection = "sheets"
	defaultCharacterArchive    = "sheets_Archive"
	defaultlogLevel            = "trace"
	defaultJWTTokenDecoder     = "http://localhost:8080"
)
