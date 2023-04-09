# Onchain Aggregator

Onchain Aggregator is a Go package that implements queries for uniswapV2, as defined in their API reference. These queries retrieve data from the Uniswap smart contract running on the Ethereum blockchain.
[uniswapV2 API reference](https://docs.uniswap.org/contracts/v2/reference/API/queries)
## About The Graph

[The Graph](https://thegraph.com/docs/en/about/) is a decentralized protocol for indexing and querying blockchain data. The Graph makes it possible to query data that is difficult to query directly. Projects with complex smart contracts like Uniswap store data on the Ethereum blockchain, making it really difficult to read anything other than basic data directly from the blockchain.
The Graph solves this with a decentralized protocol that indexes and enables the performant and efficient querying of blockchain data. These APIs (indexed "subgraphs") can then be queried with a standard GraphQL API.

To test subgraph queries, you can use the Uniswap v2 subgraph on the hosted service provided by The Graph, which can be accessed on the [playground](https://thegraph.com/hosted-service/subgraph/uniswap/uniswap-v2)

## Functions

Onchain Aggregator provides the following functions, which implement queries defined by Uniswap and can be found in their API reference:

Global Data
- QueryGlobalStats
- QueryGlobalHistoricalLookup

Pair Data
- QueryRecentSwapsFromPair
- QueryPairOverview
- QueryAllUniswapPairs
- QueryMostLiquidPairs
- QueryRecentSwapsFromPair
- QueryPairDailyAggregated

Token Data
- QueryTokenOverview
- QueryTokenData
- QueryAllUniswapTokens
- QueryTokenTransactions
- Query



## Requirements
- Go v1.16 or higher
- The following Go packages:
  - github.com/machinebox/graphql

## About

Onchain Aggregator is a Go package that retrieves onchain data from the Uniswap smart contract and stores it in a database. The raw data will subsequently undergo transformation to extract meaningful information. At present, it only gathers data from Uniswap, but this package may expand to include other blockchains and platforms in the future.

If you use this package, please provide a reference to the original source.

## License

This package is released under the MIT License. See the LICENSE file for details.