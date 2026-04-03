package tools

//go:generate go tool -modfile=../../go.tools.mod swag init --dir cmd,internal/handler/v1 --output ../../docs/api/v1 --ot json --parseInternal true
