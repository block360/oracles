package feeder

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	e2e "github.com/makerdao/setzer-e2e"

	"github.com/chronicleprotocol/infestor"
	"github.com/chronicleprotocol/infestor/origin"
	"github.com/stretchr/testify/suite"
)

func TestFeedBaseBehaviourE2ESuite(t *testing.T) {
	suite.Run(t, new(FeedBaseBehaviourE2ESuite))
}

type FeedBaseBehaviourE2ESuite struct {
	e2e.SmockerAPISuite
}

func (s *FeedBaseBehaviourE2ESuite) TestPartialInvalidPricesLessThanMin() {
	ctx, cancel := context.WithTimeout(context.Background(), e2e.OmniaDefaultTimeout)
	defer cancel()

	s.Omnia = e2e.NewOmniaFeedProcess(ctx)

	// Setup price for BTC/USD
	err := infestor.NewMocksBuilder().
		Reset().
		Add(origin.NewExchange("bitstamp").WithSymbol("BTC/USD").WithPrice(1)).
		Add(origin.NewExchange("bittrex").WithSymbol("BTC/USD").WithStatusCode(http.StatusNotFound)).
		Add(origin.NewExchange("coinbase").WithSymbol("BTC/USD").WithStatusCode(http.StatusConflict)).
		Add(origin.NewExchange("gemini").WithSymbol("BTC/USD").WithStatusCode(http.StatusConflict)).
		Add(origin.NewExchange("kraken").WithSymbol("XXBT/ZUSD").WithStatusCode(http.StatusConflict)).
		Deploy(s.API)

	s.Require().NoError(err)

	err = s.Omnia.Start()
	s.Assert().NoError(err)
	// fmt.Println(s.Omnia.StdoutString())

	time.Sleep(3 * time.Second)

	err = s.Omnia.Stop()
	s.Assert().NoError(err)

	empty, err := s.Transport.IsEmpty()
	s.Assert().NoError(err)
	s.Assert().True(empty, "E2E Transport send message, but should not")
}

func (s *FeedBaseBehaviourE2ESuite) TestAllInvalidPrices() {
	ctx, cancel := context.WithTimeout(context.Background(), e2e.OmniaDefaultTimeout)
	defer cancel()

	s.Omnia = e2e.NewOmniaFeedProcess(ctx)

	// Setup price for BTC/USD
	err := infestor.NewMocksBuilder().
		Reset().
		Add(origin.NewExchange("bitstamp").WithSymbol("BTC/USD").WithStatusCode(http.StatusConflict)).
		Add(origin.NewExchange("bittrex").WithSymbol("BTC/USD").WithStatusCode(http.StatusConflict)).
		Add(origin.NewExchange("coinbase").WithSymbol("BTC/USD").WithStatusCode(http.StatusConflict)).
		Add(origin.NewExchange("gemini").WithSymbol("BTC/USD").WithStatusCode(http.StatusConflict)).
		Add(origin.NewExchange("kraken").WithSymbol("XXBT/ZUSD").WithStatusCode(http.StatusConflict)).
		Deploy(s.API)

	s.Require().NoError(err)

	err = s.Omnia.Start()
	s.Assert().NoError(err)
	// TODO: sleep ?
	// fmt.Println(s.Omnia.StdoutString())

	time.Sleep(3 * time.Second)

	err = s.Omnia.Stop()
	s.Assert().NoError(err)

	s.Assert().True(s.Transport.IsEmpty())
}

func (s *FeedBaseBehaviourE2ESuite) TestMinValuablePrices() {
	ctx, cancel := context.WithTimeout(context.Background(), e2e.OmniaDefaultTimeout)
	defer cancel()

	s.Omnia = e2e.NewOmniaFeedProcess(ctx)

	// Setup price for BTC/USD
	err := infestor.NewMocksBuilder().
		Reset().
		Add(origin.NewExchange("bitstamp").WithSymbol("BTC/USD").WithPrice(1)).
		Add(origin.NewExchange("bittrex").WithSymbol("BTC/USD").WithPrice(1)).
		Add(origin.NewExchange("coinbase").WithSymbol("BTC/USD").WithPrice(1)).
		Add(origin.NewExchange("gemini").WithSymbol("BTC/USD").WithStatusCode(http.StatusConflict)).
		Add(origin.NewExchange("kraken").WithSymbol("XXBT/ZUSD").WithStatusCode(http.StatusConflict)).
		Deploy(s.API)

	s.Require().NoError(err)

	err = s.Omnia.Start()
	s.Assert().NoError(err)
	// fmt.Println(s.Omnia.StdoutString())

	ch, err := s.Transport.ReadChan()
	s.Require().NoError(err)
	s.Require().NotNil(ch)

	msg, err := s.Transport.WaitMsg(15 * time.Second)
	s.Require().NoError(err)

	var price e2e.PriceMessage

	err = json.Unmarshal([]byte(msg), &price)
	s.Require().NoError(err)

	s.Assert().Equal("BTCUSD", price.Price.Wat)
	s.Assert().Equal("1000000000000000000", price.Price.Val)
	s.Assert().Greater(time.Now().Unix(), price.Price.Age)
	s.Assert().NotEmpty(price.Price.R)
	s.Assert().NotEmpty(price.Price.S)
	s.Assert().NotEmpty(price.Price.V)

	// Check trace ?
	s.Assert().Equal("1.0000000000", price.Trace["bitstamp"])
	s.Assert().Equal("1.0000000000", price.Trace["bittrex"])
	s.Assert().Equal("1.0000000000", price.Trace["coinbase"])
	// Should not be in trace list
	_, ok := price.Trace["gemini"]
	s.Assert().False(ok)
	_, ok = price.Trace["kraken"]
	s.Assert().False(ok)

}

func (s *FeedBaseBehaviourE2ESuite) TestBaseSuccessBehaviour() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	s.Omnia = e2e.NewOmniaFeedProcess(ctx)

	// Setup price for BTC/USD
	err := infestor.NewMocksBuilder().
		Reset().
		Add(origin.NewExchange("bitstamp").WithSymbol("BTC/USD").WithPrice(1)).
		Add(origin.NewExchange("bittrex").WithSymbol("BTC/USD").WithPrice(1)).
		Add(origin.NewExchange("coinbase").WithSymbol("BTC/USD").WithPrice(1)).
		Add(origin.NewExchange("gemini").WithSymbol("BTC/USD").WithPrice(1)).
		Add(origin.NewExchange("kraken").WithSymbol("XXBT/ZUSD").WithPrice(1)).
		Deploy(s.API)

	s.Require().NoError(err)

	err = s.Omnia.Start()
	s.Require().NoError(err)

	ch, err := s.Transport.ReadChan()
	s.Require().NoError(err)
	s.Require().NotNil(ch)

	msg, err := s.Transport.WaitMsg(15 * time.Second)
	s.Require().NoError(err)

	var price e2e.PriceMessage

	err = json.Unmarshal([]byte(msg), &price)
	s.Require().NoError(err)

	s.Assert().Equal("BTCUSD", price.Price.Wat)
	s.Assert().Equal("1000000000000000000", price.Price.Val)
	s.Assert().Greater(time.Now().Unix(), price.Price.Age)
	s.Assert().NotEmpty(price.Price.R)
	s.Assert().NotEmpty(price.Price.S)
	s.Assert().NotEmpty(price.Price.V)

	// Check trace ?
	s.Assert().Equal("1.0000000000", price.Trace["bitstamp"])
	s.Assert().Equal("1.0000000000", price.Trace["bittrex"])
	s.Assert().Equal("1.0000000000", price.Trace["coinbase"])
	s.Assert().Equal("1.0000000000", price.Trace["gemini"])
	s.Assert().Equal("1.0000000000", price.Trace["kraken"])

	// Next call
	// Setup price for BTC/USD
	err = infestor.NewMocksBuilder().
		Reset().
		Add(origin.NewExchange("bitstamp").WithSymbol("BTC/USD").WithPrice(2)).
		Add(origin.NewExchange("bittrex").WithSymbol("BTC/USD").WithPrice(2)).
		Add(origin.NewExchange("coinbase").WithSymbol("BTC/USD").WithPrice(2)).
		Add(origin.NewExchange("gemini").WithSymbol("BTC/USD").WithPrice(2)).
		Add(origin.NewExchange("kraken").WithSymbol("XXBT/ZUSD").WithPrice(2)).
		Deploy(s.API)

	s.Require().NoError(err)

	msg, err = s.Transport.WaitMsg(15 * time.Second)
	s.Require().NoError(err)

	err = json.Unmarshal([]byte(msg), &price)
	s.Require().NoError(err)

	s.Assert().Equal("BTCUSD", price.Price.Wat)
	s.Assert().Equal("2000000000000000000", price.Price.Val)
	s.Assert().Greater(time.Now().Unix(), price.Price.Age)
	s.Assert().NotEmpty(price.Price.R)
	s.Assert().NotEmpty(price.Price.S)
	s.Assert().NotEmpty(price.Price.V)

	// Check trace ?
	s.Assert().Equal("2.0000000000", price.Trace["bitstamp"])
	s.Assert().Equal("2.0000000000", price.Trace["bittrex"])
	s.Assert().Equal("2.0000000000", price.Trace["coinbase"])
	s.Assert().Equal("2.0000000000", price.Trace["gemini"])
	s.Assert().Equal("2.0000000000", price.Trace["kraken"])

	// 3rd step
	// Setup price for BTC/USD
	err = infestor.NewMocksBuilder().
		Reset().
		Add(origin.NewExchange("bitstamp").WithSymbol("BTC/USD").WithPrice(3)).
		Add(origin.NewExchange("bittrex").WithSymbol("BTC/USD").WithPrice(3)).
		Add(origin.NewExchange("coinbase").WithSymbol("BTC/USD").WithPrice(3)).
		Add(origin.NewExchange("gemini").WithSymbol("BTC/USD").WithPrice(3)).
		Add(origin.NewExchange("kraken").WithSymbol("XXBT/ZUSD").WithPrice(3)).
		Deploy(s.API)

	s.Require().NoError(err)

	msg, err = s.Transport.WaitMsg(15 * time.Second)
	s.Require().NoError(err)

	err = json.Unmarshal([]byte(msg), &price)
	s.Require().NoError(err)

	s.Assert().Equal("BTCUSD", price.Price.Wat)
	s.Assert().Equal("3000000000000000000", price.Price.Val)
	s.Assert().Greater(time.Now().Unix(), price.Price.Age)
	s.Assert().NotEmpty(price.Price.R)
	s.Assert().NotEmpty(price.Price.S)
	s.Assert().NotEmpty(price.Price.V)

	// Check trace ?
	s.Assert().Equal("3.0000000000", price.Trace["bitstamp"])
	s.Assert().Equal("3.0000000000", price.Trace["bittrex"])
	s.Assert().Equal("3.0000000000", price.Trace["coinbase"])
	s.Assert().Equal("3.0000000000", price.Trace["gemini"])
	s.Assert().Equal("3.0000000000", price.Trace["kraken"])
}
