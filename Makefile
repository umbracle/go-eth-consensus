
sszgen:
	sszgen --path structs.go --exclude-objs Root,Signature,Uint256
	sszgen --path ./http/validator.go --objs RegisterValidatorRequest --output ./http/builder_encoding.go

get-spec-tests:
	./scripts/download-spec-tests.sh v1.3.0-rc.2

abigen-deposit:
	ethgo abigen --source ./internal/deposit/deposit.abi --package deposit --output ./internal/deposit/
