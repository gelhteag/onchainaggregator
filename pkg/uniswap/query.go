package uniswap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/machinebox/graphql"
)

// ██   ██ ███████ ██      ██████  ███████ ██████      ███████ ██    ██ ███    ██  ██████ ████████ ██  ██████  ███    ██ ███████
// ██   ██ ██      ██      ██   ██ ██      ██   ██     ██      ██    ██ ████   ██ ██         ██    ██ ██    ██ ████   ██ ██
// ███████ █████   ██      ██████  █████   ██████      █████   ██    ██ ██ ██  ██ ██         ██    ██ ██    ██ ██ ██  ██ ███████
// ██   ██ ██      ██      ██      ██      ██   ██     ██      ██    ██ ██  ██ ██ ██         ██    ██ ██    ██ ██  ██ ██      ██
// ██   ██ ███████ ███████ ██      ███████ ██   ██     ██       ██████  ██   ████  ██████    ██    ██  ██████  ██   ████ ███████

func generateQueryFromStruct(obj interface{}, queryName string, id string, args map[string]interface{}) string {
	if args == nil {
		return fmt.Sprintf(`{ %s(id: "%s"){ %s } }`, queryName, id, BuildFields(obj))
	}
	return fmt.Sprintf(`{ %s(id: "%s", %s){ %s } }`, queryName, id, BuildArgs(args), BuildFields(obj))
}

func BuildArgs(args map[string]interface{}) string {
	var argStrings []string
	for k, v := range args {
		argValue := ""
		switch v := v.(type) {
		case string:
			argValue = fmt.Sprintf(`"%s"`, v)
		case int:
			argValue = fmt.Sprintf("%d", v)
		case float64:
			argValue = fmt.Sprintf("%f", v)
		case map[string]interface{}:
			argValue = fmt.Sprintf("{%s}", BuildArgs(v))
		default:
			argValue = fmt.Sprintf("%v", v)
		}
		argStrings = append(argStrings, fmt.Sprintf("%s: %s", k, argValue))
	}
	return strings.Join(argStrings, ", ")
}
func BuildFields(obj interface{}) string {
	var value reflect.Value
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		value = reflect.ValueOf(obj).Elem()
	} else {
		value = reflect.ValueOf(obj)
	}

	var fields []string
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		tag := field.Tag.Get("graphql")
		if tag != "" {
			fields = append(fields, tag)
		}
	}
	return strings.Join(fields, "\n")
}

func RunGraphQLQuery(query string) (map[string]interface{}, error) {
	// Set up a new GraphQL client
	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2")

	// Create a new request with the GraphQL query
	req := graphql.NewRequest(query)

	// Execute the request
	var responseData map[string]interface{}
	err := client.Run(context.Background(), req, &responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

// ██████   ██       ██████  ██████   █████  ██          ██████   █████  ████████  █████
// ██       ██      ██    ██ ██   ██ ██   ██ ██          ██   ██ ██   ██    ██    ██   ██
// ██   ███ ██      ██    ██ ██████  ███████ ██          ██   ██ ███████    ██    ███████
// ██    ██ ██      ██    ██ ██   ██ ██   ██ ██          ██   ██ ██   ██    ██    ██   ██
//	██████  ███████  ██████  ██████  ██   ██ ███████     ██████  ██   ██    ██    ██   ██

type GlobalStats struct {
	TotalVolumeUSD    string `graphql:"totalVolumeUSD"`
	TotalLiquidityUSD string `graphql:"totalLiquidityUSD"`
	TxCount           string `graphql:"txCount"`
}

// All time volume in USD, total liquidity in USD, all time transaction count.
func QueryGlobalStats(factoryID string) string {
	globalStatsQuery := GlobalStats{}
	query := generateQueryFromStruct(&globalStatsQuery, "uniswapFactory", factoryID, nil)
	return query
}

// To get a snapshot of past state, use The Graph's block query feature and query at a previous block.
// See this post https://blocklytics.org/blog/ethereum-blocks-subgraph-made-for-time-travel/
// to get more information about fetching block numbers from timestamps.
// This can be used to calculate things like 24hr volume.

func QueryGlobalHistoricalLookup(factoryID string, blockNumber int) string {
	globalHistoricalLookup := GlobalStats{}
	queryArgs := map[string]interface{}{
		"block": map[string]interface{}{
			"number": blockNumber,
		},
	}
	query := generateQueryFromStruct(&globalHistoricalLookup, "uniswapFactory", factoryID, queryArgs)
	return query
}

// ██████   █████  ██ ██████      ██████   █████  ████████  █████
// ██   ██ ██   ██ ██ ██   ██     ██   ██ ██   ██    ██    ██   ██
// ██████  ███████ ██ ██████      ██   ██ ███████    ██    ███████
// ██      ██   ██ ██ ██   ██     ██   ██ ██   ██    ██    ██   ██
// ██      ██   ██ ██ ██   ██     ██████  ██   ██    ██    ██   ██

type Token struct {
	ID         string `graphql:"id"`
	Symbol     string `graphql:"symbol"`
	Name       string `graphql:"name"`
	DerivedETH string `graphql:"derivedETH"`
}

type PairData struct {
	Token0     *Token `graphql:"token0"`
	Token1     *Token `graphql:"token1"`
	Reserve0   string `graphql:"reserve0"`
	Reserve1   string `graphql:"reserve1"`
	ReserveUSD string `graphql:"reserveUSD"`
	VolumeUSD  string `graphql:"volumeUSD"`
	TxCount    string `graphql:"txCount"`
}
type Pairs struct {
	ID string `graphql:"id"`
}
type RecentSwapsFromPairQuery struct {
	Pair struct {
		Token0 struct {
			Symbol string `graphql:"symbol"`
		} `graphql:"token0"`
		Token1 struct {
			Symbol string `graphql:"symbol"`
		} `graphql:"token1"`
	} `graphql:"pair"`
	Amount0In  string `graphql:"amount0In"`
	Amount0Out string `graphql:"amount0Out"`
	Amount1In  string `graphql:"amount1In"`
	Amount1Out string `graphql:"amount1Out"`
	AmountUSD  string `graphql:"amountUSD"`
	To         string `graphql:"to"`
}
type PairDailyAggregated struct {
	Date              string `graphql:"date"`
	DailyVolumeToken0 string `graphql:"dailyVolumeToken0"`
	DailyVolumeToken1 string `graphql:"dailyVolumeToken1"`
	DailyVolumeUSD    string `graphql:"dailyVolumeUSD"`
	ReserveUSD        string `graphql:"reserveUSD"`
}

func QueryPairOverview(pairID string) (*PairData, error) {
	pairQuery := PairData{}
	query := generateQueryFromStruct(&pairQuery, "pair", pairID, nil)
	pairDataResponse, err := RunGraphQLQuery(query)
	if err != nil {
		log.Fatalf("Error executing  GraphQL request: %v", err)
	}
	pairData, ok := pairDataResponse["pair"].(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected response format")
	}
	var pairOverview *PairData
	pairDataJSON, err := json.Marshal(pairData)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(pairDataJSON, &pairOverview)
	if err != nil {
		return nil, err
	}

	return pairOverview, nil
}

// The Graph limits entity return amounts to 1000 per query as of now.
// To get all pairs on Uniswap use a loop and graphql skip query to fetch multiple chunks of 1000 pairs.
func QueryAllUniswapPairs(skip int) (*[]Pairs, error) {
	pairQuery := Pairs{}
	args := map[string]interface{}{
		"skip": skip,
	}
	query := generateQueryFromStruct(&pairQuery, "pairs", "", args)
	pairsResponse, err := RunGraphQLQuery(query)
	if err != nil {
		log.Fatalf("Error executing  GraphQL request: %v", err)
	}

	// pairsData, ok := pairsResponse["pairs"].(map[string]interface{})
	// if !ok {
	// 	return nil, errors.New("unexpected response format")
	// }
	// fmt.Printf("pairsResponse: %v", pairsResponse["pairs"])
	var pairs *[]Pairs
	pairsJSON, err := json.Marshal(pairsResponse["pairs"])
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(pairsJSON, &pairs)
	if err != nil {
		return nil, err
	}
	return pairs, nil
}

// Order by liquidity to get the most liquid pairs in Uniswap.

// QueryMostLiquidPairs generates a GraphQL query to fetch the most liquid pairs.
// The function takes a map as an argument, which contains the following default values:
//
//	args := map[string]interface{}{
//		"first":          1000,                 // Fetch the first 1000 pairs
//		"orderBy":        "reserveUSD",         // Order the pairs by their reserveUSD
//		"orderDirection": "desc",               // Sort the pairs in descending order (highest liquidity first)
//	}
func QueryMostLiquidPairs(args map[string]interface{}) (*[]Pairs, error) {
	pairQuery := Pairs{}
	query := generateQueryFromStruct(&pairQuery, "pairs", "", args)
	pairsResponse, err := RunGraphQLQuery(query)
	if err != nil {
		log.Fatalf("Error executing  GraphQL request: %v", err)
	}

	// pairsData, ok := pairsResponse["pairs"].(map[string]interface{})
	// if !ok {
	// 	return nil, errors.New("unexpected response format")
	// }
	// fmt.Printf("pairsResponse: %v", pairsResponse["pairs"])
	var pairs *[]Pairs
	pairsJSON, err := json.Marshal(pairsResponse["pairs"])
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(pairsJSON, &pairs)
	if err != nil {
		return nil, err
	}
	return pairs, nil
}

// Get the last 100 swaps on a pair by fetching Swap events and passing in the pair address.
// You'll often want token information as well.

// QueryRecentSwapsFromPair generates a GraphQL query to fetch recent swaps for a specific pair.
// The function takes a map as an argument, which contains the following default values:
//
//	args := map[string]interface{}{
//		"orderBy":        "timestamp",          // Order the swaps by their timestamp
//		"orderDirection": "desc",               // Sort the swaps in descending order (latest first)
//		"where": map[string]interface{}{
//			"pair": pairID,                      // Filter the swaps based on the pairID
//		},
//	}
func QueryRecentSwapsFromPair(args map[string]interface{}) (*[]RecentSwapsFromPairQuery, error) {
	queryStruct := RecentSwapsFromPairQuery{}
	query := generateQueryFromStruct(&queryStruct, "swaps", "", args)
	recentSwapsFromPairQueryResponse, err := RunGraphQLQuery(query)
	if err != nil {
		log.Fatalf("Error executing  GraphQL request: %v", err)
	}
	var recentSwapsFromPairQuery *[]RecentSwapsFromPairQuery
	pairsJSON, err := json.Marshal(recentSwapsFromPairQueryResponse["swaps"])
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(pairsJSON, &recentSwapsFromPairQuery)
	if err != nil {
		return nil, err
	}
	return recentSwapsFromPairQuery, nil
}

// QueryPairDailyAggregated generates a GraphQL query to fetch daily aggregated data for a pair.
// The function takes a map as an argument, which contains the following default values:
//
//	args := map[string]interface{}{
//		"first":          100,                  // Fetch the first 100 data points
//		"orderBy":        "date",               // Order the data points by date
//		"orderDirection": "asc",                // Sort the data points in ascending order (oldest first)
//		"where": map[string]interface{}{
//			"pairAddress": pairAddress,         // The pair address to filter by
//			"date_gt":     timestamp,           // Fetch data points with a date greater than the provided timestamp
//		},
//	}
func QueryPairDailyAggregated(args map[string]interface{}) (*[]PairDailyAggregated, error) {
	queryStruct := PairDailyAggregated{}
	query := generateQueryFromStruct(&queryStruct, "pairDayDatas", "", args)
	pairDailyAggregatedResponse, err := RunGraphQLQuery(query)
	if err != nil {
		log.Fatalf("Error executing  GraphQL request: %v", err)
	}
	var pairDailyAggregated *[]PairDailyAggregated
	pairsJSON, err := json.Marshal(pairDailyAggregatedResponse["swaps"])
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(pairsJSON, &pairDailyAggregated)
	if err != nil {
		return nil, err
	}
	return pairDailyAggregated, nil
}

// ████████  ██████  ██   ██ ███████ ███    ██     ██████   █████  ████████  █████
//    ██    ██    ██ ██  ██  ██      ████   ██     ██   ██ ██   ██    ██    ██   ██
//    ██    ██    ██ █████   █████   ██ ██  ██     ██   ██ ███████    ██    ███████
//    ██    ██    ██ ██  ██  ██      ██  ██ ██     ██   ██ ██   ██    ██    ██   ██
//    ██     ██████  ██   ██ ███████ ██   ████     ██████  ██   ██    ██    ██   ██

// Token data can be fetched using the token contract address as an ID.
// Token data is aggregated across all pairs the token is included in.
// Any token that is included in some pair in Uniswap can be queried.
type TokenOverview struct {
	Name           string `graphql:"name"`
	Symbol         string `graphql:"symbol"`
	Decimals       string `graphql:"decimals"`
	DerivedETH     string `graphql:"derivedETH"`
	TradeVolumeUSD string `graphql:"tradeVolumeUSD"`
	TotalLiquidity string `graphql:"totalLiquidity"`
}
type TokenData struct {
	ID         string `graphql:"id"`
	Symbol     string `graphql:"symbol"`
	Name       string `graphql:"name"`
	DerivedETH string `graphql:"derivedETH"`
}
type TokenDayData struct {
	ID                string `graphql:"id"`
	Date              int    `graphql:"date"`
	PriceUSD          string `graphql:"priceUSD"`
	TotalLiquidity    string `graphql:"totalLiquidityToken"`
	TotalLiquidityUSD string `graphql:"totalLiquidityUSD"`
	TotalLiquidityETH string `graphql:"totalLiquidityETH"`
	DailyVolumeETH    string `graphql:"dailyVolumeETH"`
	DailyVolume       string `graphql:"dailyVolumeToken"`
	DailyVolumeUSD    string `graphql:"dailyVolumeUSD"`
}
type Transaction struct {
	ID        string `graphql:"id"`
	Timestamp string `graphql:"timestamp"`
}

type Mint struct {
	Transaction *Transaction `graphql:"transaction"`
	To          string       `graphql:"to"`
	Liquidity   string       `graphql:"liquidity"`
	Amount0     string       `graphql:"amount0"`
	Amount1     string       `graphql:"amount1"`
	AmountUSD   string       `graphql:"amountUSD"`
}

type Burn struct {
	Transaction *Transaction `graphql:"transaction"`
	To          string       `graphql:"to"`
	Liquidity   string       `graphql:"liquidity"`
	Amount0     string       `graphql:"amount0"`
	Amount1     string       `graphql:"amount1"`
	AmountUSD   string       `graphql:"amountUSD"`
}

type Swap struct {
	Transaction *Transaction `graphql:"transaction"`
	Amount0In   string       `graphql:"amount0In"`
	Amount0Out  string       `graphql:"amount0Out"`
	Amount1In   string       `graphql:"amount1In"`
	Amount1Out  string       `graphql:"amount1Out"`
	AmountUSD   string       `graphql:"amountUSD"`
	To          string       `graphql:"to"`
}

type TokenTransactions struct {
	Mints []*Mint `graphql:"mints"`
	Burns []*Burn `graphql:"burns"`
	Swaps []*Swap `graphql:"swaps"`
}

// Get a snapshot of the current stats on a token in Uniswap.
// This query fetches current stats of the given Token.
func QueryTokenOverview(tokenID string) (*TokenOverview, error) {
	// Define a new instance of TokenOverview to store the response
	queryStruct := &TokenOverview{}
	// Generate the GraphQL query string from the TokenOverview struct
	query := generateQueryFromStruct(queryStruct, "token", tokenID, nil)
	tokenOverviewResponse, err := RunGraphQLQuery(query)
	if err != nil {
		log.Fatalf("Error executing GraphQL request: %v", err)
	}
	tokenOverview, ok := tokenOverviewResponse["token"].(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected response format")
	}
	tokenOverviewJSON, err := json.Marshal(tokenOverview)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(tokenOverviewJSON, queryStruct)
	if err != nil {
		return nil, err
	}
	return queryStruct, nil
}

// QueryTokenData queries the Uniswap GraphQL API for data about a specific token on the Uniswap exchange. It
// returns a TokenData struct containing information about the token's name, symbol, ID, and derived ETH, or an
// error if the query fails.
func QueryTokenData(tokenID string) (*TokenData, error) {
	// Define a new instance of TokenData to store the response
	queryStruct := &TokenData{}

	// Generate the GraphQL query string from the TokenData struct
	query := generateQueryFromStruct(queryStruct, "token", tokenID, nil)

	tokenDataResponse, err := RunGraphQLQuery(query)
	if err != nil {
		log.Fatalf("Error executing  GraphQL request: %v", err)
	}
	tokenData, ok := tokenDataResponse["token"].(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected response format")
	}
	// var tokenDataInterface *Token
	tokenDataJSON, err := json.Marshal(tokenData)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(tokenDataJSON, queryStruct)
	if err != nil {
		return nil, err
	}

	return queryStruct, nil
}

// QueryAllUniswapTokens queries the Uniswap GraphQL API for all tokens on the Uniswap exchange. It returns
// three slices containing the IDs, names, and symbols of all tokens found. If an error occurs while executing
// the query, the function returns an error and the slices will be nil.
func QueryAllUniswapTokens(skip int) (*[]TokenData, error) {
	// Define a new instance of TokenData to store the response
	queryStruct := &TokenData{}
	// Set up arguments for the query
	args := map[string]interface{}{
		"skip": skip,
	}
	// Generate the GraphQL query string from the TokenData struct
	query := generateQueryFromStruct(queryStruct, "tokens", "", args)
	allUniswapTokensResponse, err := RunGraphQLQuery(query)
	if err != nil {
		log.Fatalf("Error executing  GraphQL request: %v", err)
	}
	var allUniswapTokens *[]TokenData
	allUniswapTokensJSON, err := json.Marshal(allUniswapTokensResponse["tokens"])
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(allUniswapTokensJSON, &allUniswapTokens)
	if err != nil {
		return nil, err
	}
	return allUniswapTokens, nil
}

// QueryTokenTransactions queries the Uniswap GraphQL API for mints, burns, and swaps transactions
// for a list of given pairs, up to a specified number of transactions per pair. It returns the
// transactions separated by type (mints, burns, and swaps) in three slices of their respective types.
// If an error occurs while executing the query, the function returns an error and the transaction slices
// will be nil
func QueryTokenTransactions(client *graphql.Client, allPairs []string, first int) ([]*Mint, []*Burn, []*Swap, error) {
	query := `
	query($allPairs: [String!]) {
	  mints(first: $first, where: { pair_in: $allPairs }, orderBy: timestamp, orderDirection: desc) {
		transaction {
		  id
		  timestamp
		}
		to
		liquidity
		amount0
		amount1
		amountUSD
	  }
	  burns(first: $first, where: { pair_in: $allPairs }, orderBy: timestamp, orderDirection: desc) {
		transaction {
		  id
		  timestamp
		}
		to
		liquidity
		amount0
		amount1
		amountUSD
	  }
	  swaps(first: $first, where: { pair_in: $allPairs }, orderBy: timestamp, orderDirection: desc) {
		transaction {
		  id
		  timestamp
		}
		amount0In
		amount0Out
		amount1In
		amount1Out
		amountUSD
		to
	  }
	}
	`

	req := graphql.NewRequest(query)
	req.Var("allPairs", allPairs)
	req.Var("first", first)
	fmt.Printf("req: %v", req)
	var responseData struct {
		Mints []*Mint `graphql:"mints"`
		Burns []*Burn `graphql:"burns"`
		Swaps []*Swap `graphql:"swaps"`
	}
	err := client.Run(context.Background(), req, &responseData)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Error executing token transactions GraphQL request: %v", err)
	}

	return responseData.Mints, responseData.Burns, responseData.Swaps, nil
}

// Like pair and global daily lookups, tokens have daily entities that can be queries as well.
// This query gets daily information for DAI.
// Note that you may want to sort in ascending order to receive your days from oldest to most recent in the return array.

// QueryTokenDailyData retrieves daily aggregated data for a specific token
// from the Uniswap GraphQL API using the provided arguments as query variables.
// Returns a slice of TokenDayData or an error if the query fails.
//
//	args := map[string]interface{}{
//				"orderBy":        "date",          // Order  by daily entities
//				"orderDirection": "asc",               // Sort the daily entities in ascending order (first latest)
//				"where": map[string]interface{}{
//					"token": pairID,
//				},
//			}
func QueryTokenDailyData(args map[string]interface{}) (*[]TokenDayData, error) {

	query := &TokenDayData{}

	// Generate the GraphQL query string from the TokenData struct
	tokenDailyData := generateQueryFromStruct(query, "tokenDayDatas", "", args)
	tokenDailyDataResponse, err := RunGraphQLQuery(tokenDailyData)
	if err != nil {
		log.Fatalf("Error executing  GraphQL request: %v", err)
	}
	var tokenDailyDataSlice *[]TokenDayData
	tokenDailyDataJSON, err := json.Marshal(tokenDailyDataResponse["tokenDayDatas"])
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(tokenDailyDataJSON, &tokenDailyDataSlice)
	if err != nil {
		return nil, err
	}

	return tokenDailyDataSlice, nil
}
