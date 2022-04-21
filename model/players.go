package model

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/icrowley/fake"
)

type Player struct {
	ID                uint64    `json:"id"`
	GUID              string    `json:"guid"`
	License           string    `json:"license"`
	PlayerID          uint64    `json:"playerID"`
	ClickID           string    `json:"clickID"`
	RegisterDate      time.Time `json:"registerDate"`
	Language          string    `json:"language"`
	Email             string    `json:"email"`
	IsEmailVerify     bool      `json:"isEmailVerify"`
	Phone             string    `json:"phone"`
	IsPhoneVerify     bool      `json:"isPhoneVerify"`
	IsMultiAccount    bool      `json:"isMultiAccount"`
	Birthday          string    `json:"birthday"`
	AccountVerifyTime time.Time `json:"accountVerifyTime"`
	LastLoginTime     time.Time `json:"lastLoginTime"`
	Country           string    `json:"country"`
	City              string    `json:"city"`
	Currency          string    `json:"currency"`
	Sex               string    `json:"sex"`
	IsTest            bool      `json:"isTest"`
	IsBot             bool      `json:"isBot"`
	Project           string    `json:"project"`
	ActivateStatus    string    `json:"activateStatus"`
	DepositStatus     string    `json:"depositStatus"`
	SmsStatus         string    `json:"smsStatus"`
	Domain            string    `json:"domain"`
	Webview           bool      `json:"webview"`
	IpAddress         string    `json:"ipAddress"`
	UserAgent         string    `json:"userAgent"`
	CreateUnixNano    int64     `json:"createUnixNano"`
	UpdateUnixNano    int64     `json:"updateUnixNano"`
}

func NewPlayer(id uint64) *Player {
	reg := AnyDate(1970)
	ua := fake.UserAgent()
	if len(ua) > 100 {
		ua = ua[:100]
	}

	return &Player{
		ID:                id,
		GUID:              fake.UserName(),
		License:           License(),
		PlayerID:          id,
		ClickID:           uuid.New().String(),
		RegisterDate:      reg,
		Language:          fake.Language(),
		Email:             fake.EmailAddress(),
		IsEmailVerify:     rand.Intn(3) == 0,
		Phone:             fake.Phone(),
		IsPhoneVerify:     rand.Intn(3) == 0,
		IsMultiAccount:    rand.Intn(3) == 0,
		Birthday:          fake.DomainZone(),
		AccountVerifyTime: AnyDate(reg.Year()),
		LastLoginTime:     AnyDate(reg.Year()),
		Country:           fake.Country(),
		City:              fake.City(),
		Currency:          fake.CurrencyCode(),
		Sex:               fake.Gender(),
		IsTest:            rand.Intn(3) == 0,
		IsBot:             rand.Intn(3) == 0,
		Project:           Project(),
		ActivateStatus:    randStatus(),
		DepositStatus:     randStatus(),
		SmsStatus:         randStatus(),
		Domain:            fake.DomainZone(),
		Webview:           rand.Intn(3) == 0,
		IpAddress:         fake.IPv4(),
		UserAgent:         ua,
		CreateUnixNano:    reg.UnixNano(),
		UpdateUnixNano:    AnyDate(reg.Year()).UnixNano(),
	}
}

func AnyDate(after int) time.Time {
	if after == 2022 {
		after -= 1
	}

	return time.Date(
		fake.Year(after, 2022),
		time.Month(fake.MonthNum()),
		fake.Day(),
		rand.Intn(24),
		rand.Intn(60),
		rand.Intn(60),
		0,
		time.UTC,
	)
}

func (p *Player) OK() {
	if len(p.GUID) > 40 {
		panic("guid")
	}
	if len(p.License) > 45 {
		panic("license")
	}

	if len(p.ClickID) > 40 {
		panic("click")
	}

	if len(p.Language) > 25 {
		panic("lang")
	}

	if len(p.Email) > 50 {
		panic("email")
	}

	if len(p.Phone) > 40 {
		panic("phone")
	}

	if len(p.Birthday) > 20 {
		panic("birth")
	}

	if len(p.Country) > 45 {
		panic("country")
	}

	if len(p.City) > 25 {
		panic("city")
	}
	if len(p.Currency) > 5 {
		panic("currency")
	}
	if len(p.Sex) > 7 {
		panic("sex")
	}
	if len(p.Project) > 10 {
		panic("project")
	}
	if len(p.ActivateStatus) > 30 {
		panic("actstatus")
	}
	if len(p.DepositStatus) > 30 {
		panic("depstatus")
	}
	if len(p.SmsStatus) > 30 {
		panic("smsstatus")
	}
	if len(p.IpAddress) > 15 {
		panic("ip")
	}
	if len(p.UserAgent) > 100 {
		panic("user agent")
	}
}

func License() string {
	return []string{"com", fake.Country()}[rand.Intn(2)]
}

const (
	Undefined = "undefined"
	Sport     = "sport"
	Casino    = "casino"
)

func Project() string {
	return []string{Sport, Casino}[rand.Intn(2)]
}
