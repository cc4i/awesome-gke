package game

import (
	"container/list"
	"encoding/json"
	"fmt"

	"golang.org/x/exp/maps"
)

type Player struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
	Level  int    `json:"level,omitempty"`
	Region string `json:"region,omitempty"`
	Ip     string `json:"ip,omitempty"`
	Online bool   `json:"online,omitempty"`
}

type Score struct {
	Round       int           `json:"round"`
	Selections  map[int][]int `json:"selections,omitempty"`
	RoundScore  int           `json:"roundScore,omitempty"`
	BankedScore int           `json:"bankedScore,omitempty"`
}

type Game struct {
	Id      string           `json:"id"`
	Start   int              `json:"start,omitempty"`
	End     int              `json:"end,omitempty"`
	Players []Player         `json:"withPlayer,omitempty"`
	Scores  map[string]Score `json:"scores,omitempty"`
}

type Dice struct {
	Id    string `json:"id"`
	Value int    `json:"value"`
	//Player's name
	WithPlayer string `json:"withPlayer,omitempty"`
	Fixed      bool   `json:"fixed,omitempty"`
	OnBoard    bool   `json:"onBoard,omitempty"`
}

type GameData struct {
	//name=>Player
	Players map[string]Player `json:"players,omitempty"`
	//uuid=>Game
	Games map[string]Game `json:"games,omitempty"`
	// Fisrt position -> Last position
	LeaderBoard *list.List `json:"leaderBoard,omitempty"`
}

type FarkleRequest struct {
	// "login" - Login into a game
	// "roll" - Roll all rollable dices
	// "move" - Move the dice from board to player's area
	// "switch" - Switch the turn between players
	// "bank" - Bank the score and switch the turn of the game
	// "start" - Start a new game
	// "end" - End the current game
	Action string `json:"action"`
}

type FarkelInterface interface {
	IntialGame(player1 string, player2 string) (string, error)
	Login(name string) ([]Player, error)
	SearchPlayer(name string) (Player, error)
	RollDices() error
	MoveDice(d Dice) error
}

var Gdata = &GameData{
	Players:     make(map[string]Player),
	Games:       make(map[string]Game),
	LeaderBoard: list.New(),
}

// Handle actions of Farkle
func FarkleHandler(txt string) (string, error) {
	farkleRequest := FarkleRequest{}
	err := json.Unmarshal([]byte(txt), &farkleRequest)
	if err != nil {
		return "", err
	}
	switch farkleRequest.Action {
	case "login":
		Gdata.Login("")
		return "", nil

	}

	return "", nil
}

// Intial a new game & save it into GameData
func (gd *GameData) IntialGame(player1 string, player2 string) (string, error) {
	return nil
}

// The player logins into the game.
func (gd *GameData) Login(name string) ([]Player, error) {
	player := Player{
		Name:   name,
		Online: true,
	}
	if p, ok := gd.Players[name]; ok {
		if p.Online {
			return nil, fmt.Errorf("%s is duplicated, the name must be unique", name)
		}
	} else {
		gd.Players[name] = player
	}

	return maps.Values(gd.Players), nil
}

// Seach a player by name
func (gd *GameData) SearchPlayer(name string) (Player, error) {
	return Player{}, nil
}

// Roll all rollable dices
func (gd *GameData) RollDices() error {
	return nil
}

// Move specific dice into player's board or other way around.
// Id is dice indicator, which is one of following: dice1, dice2, dice3, dice4, dice5, dice6
func (gd *GameData) MoveDice(id string) error {
	return nil
}
