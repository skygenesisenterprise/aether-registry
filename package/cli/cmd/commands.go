package cmd

import (
	"github.com/skygenesisenterprise/aether-bank/cli/cmd/account"
	"github.com/skygenesisenterprise/aether-bank/cli/cmd/auth"
	"github.com/skygenesisenterprise/aether-bank/cli/cmd/debug"
	"github.com/skygenesisenterprise/aether-bank/cli/cmd/ledger"
	"github.com/skygenesisenterprise/aether-bank/cli/cmd/transfer"
	"github.com/skygenesisenterprise/aether-bank/cli/cmd/tx"
	"github.com/skygenesisenterprise/aether-bank/cli/cmd/user"
)

var (
	authCmd     = auth.AuthCmd
	userCmd     = user.UserCmd
	accountCmd  = account.AccountCmd
	txCmd       = tx.TxCmd
	transferCmd = transfer.TransferCmd
	ledgerCmd   = ledger.LedgerCmd
	debugCmd    = debug.DebugCmd
)
