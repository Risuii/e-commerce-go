include .env
.PHONY: migration ecommerce-migration

migration: ecommerce-migration

run:
	go mod tidy
	go run .

wire-dep:
	cd wire && ~/go/bin/wire .

test:
	go test -race -coverprofile cover.out ./internal/...
	go tool cover -html=cover.out

ecommerce-migration:
	@chmod -R 777 scripts
	@./scripts/migration.sh	$(CURDIR)/migrations/ecommerce postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable $(DB_LAST_MIGRATION_VERSION) $(target-version)

generate-mig-dev:
	@chmod +x export_git_files_to_csv_development.sh
	@./export_git_files_to_csv_development.sh

mockgen.config:
	~/go/bin/mockgen -source=config/config.go -destination=internal/auth/tests/config/init.go

mockgen.library:
	~/go/bin/mockgen -source=library/library.go  -destination=internal/auth/tests/library/mockLibrary.go

mockgen.log:
	~/go/bin/mockgen -source=pkg/logger/logger.go  -destination=internal/auth/tests/log/mockLog.go

mockgen.handler:
	~/go/bin/mockgen -source=internal/auth/domain/usecase/register.go  -destination=internal/auth/tests/delivery/presenter/http/mock/init.go

mockgen.usecase:
	~/go/bin/mockgen -source=internal/auth/domain/repository/user.go  -destination=internal/auth/tests/domain/usecase/mock/init.go

mockgen.pkg:
	~/go/bin/mockgen -source=pkg/bcrypt/bcrypt.go  -destination=internal/auth/tests/pkg/bcrypt/init.go
	~/go/bin/mockgen -source=pkg/crypto/crypto.go  -destination=internal/auth/tests/pkg/crypto/init.go