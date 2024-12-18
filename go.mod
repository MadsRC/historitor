module github.com/MadsRC/historitor

go 1.23

require (
	github.com/plar/go-adaptive-radix-tree/v2 v2.0.3
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
    v1.0.0 // Published accidentally.
)
