
sszgen:
	sszgen --path structs.go --exclude-objs Root,Signature,Uint256
	sszgen --path ./http/builder.go --objs RegisterValidatorRequest --output ./http/builder_encoding.go

get-spec-tests:
	./scripts/download-spec-tests.sh v1.1.10

abigen-deposit:
	ethgo abigen --source ./internal/deposit/deposit.abi --package deposit --output ./internal/deposit/
