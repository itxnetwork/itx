package keeper_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	evm "github.com/itxnetwork/itx/v12/x/evm/types"

	"github.com/itxnetwork/itx/v12/app"
	"github.com/itxnetwork/itx/v12/x/inflation/types"
)

var denomMint = types.DefaultInflationDenom

type KeeperTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	app            *app.Itx
	queryClientEvm evm.QueryClient
	queryClient    types.QueryClient
	consAddress    sdk.ConsAddress
}

var s *KeeperTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(KeeperTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.DoSetupTest()
}
