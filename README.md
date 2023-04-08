# Onchain Aggregator

Onchain Aggregator is a Go package that implements queries for uniswapV2, as defined in their API reference. These queries retrieve data from the Uniswap smart contract running on the Ethereum blockchain.
[uniswapV2 API reference](https://docs.uniswap.org/contracts/v2/reference/API/queries)
## About The Graph

The Graph is a decentralized protocol for indexing and querying blockchain data. The Graph makes it possible to query data that is difficult to query directly. Projects with complex smart contracts like Uniswap store data on the Ethereum blockchain, making it really difficult to read anything other than basic data directly from the blockchain.
Reference: [The Graph](https://thegraph.com/docs/en/about/)
The Graph solves this with a decentralized protocol that indexes and enables the performant and efficient querying of blockchain data. These APIs (indexed "subgraphs") can then be queried with a standard GraphQL API.

## Functions

Onchain Aggregator provides the following functions, which implement queries defined by Uniswap and can be found in their API reference:

- QueryGlobalStats
- QueryGlobalHistoricalLookup
- RecentSwapsFromPairQuery
- QueryPairOverview
- QueryAllUniswapPairs
- QueryMostLiquidPairs
- QueryRecentSwapsFromPair
- QueryPairDailyAggregated
- QueryTokenData
- QueryAllUniswapTokens
- QueryTokenTransactions
- QueryTokenDailyData


## Requirements
- Go v1.16 or higher
- The following Go packages:
  - github.com/machinebox/graphql
  
## About

Onchain Aggregator was created for the purpose of retrieving onchain data and storing it in a database. Currently, it only gathers data from Uniswap. 
However, this package may expand to include other blockchains and platforms in the future.

If you use this package, please provide a reference to the original source.

## License

This package is released under the MIT License. See the LICENSE file for details.