module server

go 1.18

require (
	github.com/disneystreaming/gomux v0.0.0-20200305000114-de122d6df124
	github.com/jessevdk/go-flags v1.5.0
	github.com/mattn/go-tty v0.0.5
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/mattn/go-tty v0.0.4 => github.com/voices-team/go-tty v0.0.5-0.20220124182555-19f3f2eb2a37
