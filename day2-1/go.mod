module day2-1

go 1.16

replace example.com/helloworld => ./helloworld

replace example.com/uuid => ./uuid

require (
	example.com/helloworld v0.0.0-00010101000000-000000000000
	example.com/uuid v0.0.0-00010101000000-000000000000
)
