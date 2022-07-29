
sszgen:
	sszgen --path structs.go --include ./common.go

get-spec-tests:
	./scripts/download-spec-tests.sh v1.1.10

abigen-deposit:
	ethgo abigen --source ./internal/deposit/deposit.abi --package deposit --output ./internal/deposit/
