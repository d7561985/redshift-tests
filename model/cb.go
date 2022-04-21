package model

import (
	"math/rand"
	"time"

	"github.com/d7561985/tel/v2"
	"github.com/google/uuid"
	"github.com/icrowley/fake"
)

type CBet struct {
	ID             uuid.UUID `json:"id"`
	License        string    `json:"license"`
	PlayerId       int       `json:"playerId"`
	GameName       string    `json:"gameName"`
	GameType       string    `json:"gameType"`
	GameID         int       `json:"gameID"`
	BonusID        int       `json:"bonusID"`
	Bet            float64   `json:"bet"`
	WinLose        float64   `json:"winLose"`
	Purse          string    `json:"purse"`
	CurrencyCode   string    `json:"currencyCode"`
	GameProvider   string    `json:"gameProvider"`
	GameRoundID    string    `json:"gameRoundID"`
	TranID         string    `json:"tranId"`
	Date           time.Time `json:"date"`
	CreateUnixNano int64     `json:"createUnixNano"`
	UpdateUnixNano int64     `json:"updateUnixNano"`
	Rollback       bool      `json:"rollback"`
	Status         string    `json:"status"`
	Error          string    `json:"error"`
	Hall           string    `json:"hall"`
	System         string    `json:"system"`
	BetInfo        string    `json:"betInfo"`
	Agent          int       `json:"agent"`
	Domain         string    `json:"domain"`
	Webview        bool      `json:"webview"`
	IsTournament   bool      `json:"isTournament"`
}

func NewCBet(playerID int) *CBet {
	bet := rand.Intn(100) + 1

	var win float64
	if rand.Intn(3) == 0 {
		win = float64(rand.Intn(bet * 10))
	}

	date := AnyDate(1970)

	return &CBet{
		ID:             uuid.New(),
		License:        License(),
		PlayerId:       playerID,
		GameName:       fake.FirstName(),
		GameType:       "SLOT",
		GameID:         rand.Intn(100000),
		BonusID:        0,
		Bet:            float64(bet),
		WinLose:        win,
		Purse:          fake.Phone(),
		CurrencyCode:   fake.CurrencyCode(),
		GameProvider:   fake.Brand(),
		GameRoundID:    uuid.New().String(),
		TranID:         uuid.New().String(),
		Date:           date,
		CreateUnixNano: date.UnixNano(),
		UpdateUnixNano: date.UnixNano(),
		Rollback:       rand.Intn(3) == 0,
		Status:         randStatus(),
		Error:          "",
		Hall:           fake.Country(),
		System:         fake.Continent(),
		BetInfo:        "",
		Agent:          rand.Intn(9999),
		Domain:         fake.DomainName(),
		Webview:        rand.Intn(3) == 0,
		IsTournament:   rand.Intn(3) == 0,
	}
}

func (c *CBet) OK() {
	ok("license", c.License, 45)
	ok("gameName", c.GameName, 30)
	ok("gameType", c.GameType, 10)
	ok("purse", c.Purse, 20)
	ok("currencyCode", c.CurrencyCode, 5)
	ok("gameProvider", c.GameProvider, 15)
	ok("gameRoundID", c.GameRoundID, 40)
	ok("tranID", c.TranID, 40)
	ok("status", c.Status, 30)
	ok("error", c.Error, 100)
	ok("hall", c.Hall, 45)
	ok("system", c.System, 20)
	ok("betInfo", c.BetInfo, 30)
	ok("domain", c.Domain, 20)
}

func ok(field, str string, l int) {
	if len(str) > l {
		tel.Global().Fatal("validation", tel.String("field", field))
	}
}
