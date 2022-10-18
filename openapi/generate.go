package openapi

//go:generate gobin -m -run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config ./signer.cfg.yaml  signer.yaml
///go:generate gobin -m -run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config ./signer-client.cfg.yaml  signer.yaml
//go:generate gobin -m -run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config ./nonce.cfg.yaml  nonce.yaml
//go:generate gobin -m -run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config ./signonce.cfg.yaml  signonce.yaml
