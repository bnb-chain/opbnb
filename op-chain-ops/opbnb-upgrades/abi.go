package opbnb_upgrades

//go:generate ./abigen --abi old-contracts/L1CrossDomainMessenger.json --pkg L1CrossDomainMessenger --out old-contracts/bindings/L1CrossDomainMessenger.go
//go:generate ./abigen --abi old-contracts/L1ERC721Bridge.json --pkg L1ERC721Bridge --out old-contracts/bindings/L1ERC721Bridge.go
//go:generate ./abigen --abi old-contracts/L1StandardBridge.json --pkg L1StandardBridge --out old-contracts/bindings/L1StandardBridge.go
//go:generate ./abigen --abi old-contracts/L2OutputOracle.json --pkg L2OutputOracle --out old-contracts/bindings/L2OutputOracle.go
//go:generate ./abigen --abi old-contracts/OptimismMintableERC20Factory.json --pkg OptimismMintableERC20Factory --out old-contracts/bindings/OptimismMintableERC20Factory.go
//go:generate ./abigen --abi old-contracts/OptimismPortal.json --pkg OptimismPortal --out old-contracts/bindings/OptimismPortal.go
//go:generate ./abigen --abi old-contracts/SystemConfig.json --pkg SystemConfig --out old-contracts/bindings/SystemConfig.go
//go:generate ./abigen --abi old-contracts/Semver.json --pkg Semver --out old-contracts/bindings/Semver.go
//go:generate ./abigen --abi old-contracts/ProxyAdmin.json --pkg ProxyAdmin --out old-contracts/bindings/ProxyAdmin.go

//go:generate ./abigen --abi new-contracts/SuperChainConfig.json --pkg SuperChainConfig --out new-contracts/bindings/SuperChainConfig.go
//go:generate ./abigen --abi new-contracts/StorageSetter.json --pkg StorageSetter --out new-contracts/bindings/StorageSetter.go
//go:generate ./abigen --abi new-contracts/L1CrossDomainMessenger.json --pkg L1CrossDomainMessenger --out new-contracts/bindings/L1CrossDomainMessenger.go
//go:generate ./abigen --abi new-contracts/L1ERC721Bridge.json --pkg L1ERC721Bridge --out new-contracts/bindings/L1ERC721Bridge.go
//go:generate ./abigen --abi new-contracts/L1StandardBridge.json --pkg L1StandardBridge --out new-contracts/bindings/L1StandardBridge.go
//go:generate ./abigen --abi new-contracts/L2OutputOracle.json --pkg L2OutputOracle --out new-contracts/bindings/L2OutputOracle.go
//go:generate ./abigen --abi new-contracts/OptimismMintableERC20Factory.json --pkg OptimismMintableERC20Factory --out new-contracts/bindings/OptimismMintableERC20Factory.go
//go:generate ./abigen --abi new-contracts/OptimismPortal.json --pkg OptimismPortal --out new-contracts/bindings/OptimismPortal.go
//go:generate ./abigen --abi new-contracts/SystemConfig.json --pkg SystemConfig --out new-contracts/bindings/SystemConfig.go
