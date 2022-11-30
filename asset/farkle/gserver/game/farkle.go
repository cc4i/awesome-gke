package game

import (
	"container/list"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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

type Dice struct {
	// The ID of dice, which is either of "dice1", "dice2", "dice3", "dice4", "dice5", "dice6"
	Id string `json:"id"`
	// The value of dice, which is either of 1, 2, 3, 4, 5, 6
	Value int `json:"value"`
	// The name of player
	WithPlayer string `json:"withPlayer,omitempty"`
	// True means its value is fixed
	Fixed bool `json:"fixed,omitempty"`
	// True is on the board, otherwise it's False when with players
	OnBoard bool `json:"onBoard"`
}

type Game struct {
	// Game Id
	Id string `json:"id"`
	// Start time of the game
	Start int64 `json:"start,omitempty"`
	// End time of the game
	End int64 `json:"end,omitempty"`
	// Players list
	Players []*Player `json:"withPlayer,omitempty"`
	// Current player, who's turn
	CurrentPlayer *Player `json:"currentPlayer,omitempty"`
	// player name => Score
	Scores map[string]*Score `json:"scores,omitempty"`
	// All dices
	Dices map[string]*Dice `json:"dices"`
}

type GameData struct {
	// player name => Player
	Players map[string]*Player `json:"players,omitempty"`
	// uuid => Game
	Games map[string]*Game `json:"games,omitempty"`
	// Fisrt position -> Last position
	LeaderBoard *list.List `json:"leaderBoard,omitempty"`
}

type FarkleRequest struct {
	// Who
	Who Player `json:"who"`

	// Do
	// "login" - Login into a game
	// "roll" - Roll all rollable dices
	// "move" - Move the dice from board to player's area
	// "switch" - Switch the turn between players
	// "bank" - Bank the score and switch the turn of the game
	// "start" - Start a new game
	// "end" - End the current game
	Action string `json:"action"`

	//What
	// dice1
	Dice1 Dice `json:"dice1,omitempty"`
	// dice2
	Dice2 Dice `json:"dice2,omitempty"`
	// dice3
	Dice3 Dice `json:"dice3,omitempty"`
	// dice4
	Dice4 Dice `json:"dice4,omitempty"`
	// dice5
	Dice5 Dice `json:"dice5,omitempty"`
	// dice6
	Dice6 Dice `json:"dice6,omitempty"`
}

type FarkelInterface interface {
	IntialGame(player1 string, player2 string) (string, error)
	Login(name string) ([]Player, error)
	SearchPlayer(name string) (Player, error)
	RollDices() ([]Dice, error)
	MoveDice(d Dice) (Score, error)
	Calculate() error
	BankScore() (Score, error)
	IsMovable() bool
}

var gdata = &GameData{
	Players:     make(map[string]*Player),
	Games:       make(map[string]*Game),
	LeaderBoard: list.New(),
}

// Handle actions of Farkle
func FarkleHandler(txt string) (string, error) {
	fr := FarkleRequest{}
	err := json.Unmarshal([]byte(txt), &fr)
	if err != nil {
		return "", err
	}
	switch fr.Action {
	case "login":
		ps, err := gdata.Login(fr.Who.Name)
		if err != nil {
			return "", err
		}
		b, _ := json.Marshal(ps)
		return string(b), nil

	}

	return "", nil
}

// Intial a new game & save it into GameData
func (gd *GameData) IntialGame(pname1 string, pname2 string) (Game, error) {
	nGame := &Game{
		Id:    uuid.New().String(),
		Start: time.Now().UnixMilli(),
	}
	player1, ok1 := gd.Players[pname1]
	player2, ok2 := gd.Players[pname2]
	if ok1 && ok2 && !player1.Online && !player2.Online {
		// Inital players
		nGame.Players = append(nGame.Players, player1)
		nGame.Players = append(nGame.Players, player2)
		nGame.CurrentPlayer = player1
		// Inital score for players
		nGame.Scores[pname1] = &Score{Round: 1}
		nGame.Scores[pname2] = &Score{Round: 1}
		// Intial dices
		for i := 1; i <= 6; i++ {
			nGame.Dices["dice"+string(i)] = &Dice{
				Id:      "dice" + string(i),
				Value:   i,
				OnBoard: true,
			}
		}

		gd.Games[nGame.Id] = nGame
	} else {
		return *nGame, fmt.Errorf("%s and %s must be logined", pname1, pname2)
	}
	return *nGame, nil
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
		gd.Players[name] = &player
	}
	var oPlayers []Player
	for _, p := range gd.Players {
		oPlayers = append(oPlayers, *p)
	}
	return oPlayers, nil
}

// Seach a player by name
func (gd *GameData) SearchPlayer(name string) (Player, error) {
	if player, ok := gd.Players[name]; ok {
		return *player, nil
	}
	return Player{}, fmt.Errorf("%s isn't existed", name)
}

// Roll all rollable dices
func (gd *GameData) RollDices(gameId string, playerName string) (Game, error) {

	rand.Seed(time.Now().UnixNano())

	if game, ok := gd.Games[gameId]; ok {
		for _, d := range game.Dices {
			// 1. Roll the dices - generate random value between 1 and 6 for dices (not with players)
			if d.OnBoard {
				d.Value = rand.Intn(6) + 1
			}
			// 2. Checking the picked dices in players' area and set 'fixed=true'
			if !d.OnBoard && d.WithPlayer != "" {
				d.Fixed = true
			}
		}

		// 3. Switch turn
		if !gd.IsMovable(gameId) {
			//3.1 Reset RoundScore=0
			score := game.Scores[playerName]
			score.RoundScore = 0

			//3.2 Move to next round
			score.Round++

			// 3.3 Swith to other player
			for _, p := range game.Players {
				if p.Name != playerName {
					game.CurrentPlayer = p
				}
			}

			//3.2 Reset dices
			for _, d := range game.Dices {
				d.OnBoard = true
				d.Fixed = false
				d.Value, _ = strconv.Atoi(strings.TrimLeft(d.Id, "dice"))
				d.WithPlayer = ""
			}
			//
		}
		return *game, nil
	}
	return Game{}, fmt.Errorf("%s is invalid game id", gameId)
}

// Move specific dice into player's board or other way around.
// Id is dice indicator, which is one of following: dice1, dice2, dice3, dice4, dice5, dice6
func (gd *GameData) MoveDice(gameId string, diceId string, playerName string) (Game, error) {

	// 1. Move a dice in or out players' area if possible
	if game, ok := gd.Games[gameId]; ok {

		// Move out
		if playerName == "" {
			if dice, ok := game.Dices[diceId]; ok {

				if !dice.Fixed {
					dice.OnBoard = true
					dice.WithPlayer = ""

					score := game.Scores[game.CurrentPlayer.Name]
					if vals, ok := score.Selections[score.Round]; ok {
						var nvals []int
						for _, v := range vals {
							if v != dice.Value {
								nvals = append(nvals, v)
							}
						}
						if len(nvals) > 0 {
							score.Selections[score.Round] = nvals
						} else {
							delete(score.Selections, score.Round)
						}
					}
				}
			}
		}
		// Move in
		if gd.IsMovable(gameId) && playerName != "" {
			if dice, ok := game.Dices[diceId]; ok {
				dice.WithPlayer = playerName
				dice.OnBoard = false

				score := game.Scores[game.CurrentPlayer.Name]
				if vals, ok := score.Selections[score.Round]; ok {
					vals = append(vals, dice.Value)
					score.Selections[score.Round] = vals
				} else {
					score.Selections[score.Round] = []int{dice.Value}
				}
			}
		}
		// Calculate score
		return *game, gd.Calculate(gameId)
	}
	return Game{}, fmt.Errorf("%s is invalid game id", gameId)
}

// Check if the dice is movable
func (gd *GameData) IsMovable(gameId string) bool {
	isMovable := false
	if game, ok := gd.Games[gameId]; ok {
		// 1; 5; 3x2; 3x3; 3x4; 3x5; 3x6; 4x?; 5x?; 6x?; 3 pairs; 1/2/3/4/5/6;
		str := ""
		for _, dice := range game.Dices {
			if dice.OnBoard {
				str = str + strconv.Itoa(dice.Value)
			}
		}
		if strings.Count(str, "1") >= 1 ||
			strings.Count(str, "5") >= 1 ||
			strings.Count(str, "2") >= 3 ||
			strings.Count(str, "3") >= 3 ||
			strings.Count(str, "4") >= 3 ||
			strings.Count(str, "5") >= 3 ||
			strings.Count(str, "6") >= 3 {
			isMovable = true
		}

		s := strings.Split(str, "")
		sort.Strings(s)
		if len(s) == 6 && s[0] == s[1] && s[2] == s[3] && s[4] == s[5] {
			isMovable = true
		}

	}
	return isMovable
}

func allEqualTo(ns []int, num int) bool {
	for i := 1; i < len(ns); i++ {
		if ns[i] != ns[0] {
			return false
		}
	}
	if num != -1 {
		return len(ns) > 0 && ns[0] == num
	}
	return true

}

// TODO:
func calculator(ns []int) int {
	return 0
}

func (gd *GameData) Calculate(gameId string) error {
	if game, ok := gd.Games[gameId]; ok {
		score := game.Scores[game.CurrentPlayer.Name]
		if sls, ok := score.Selections[score.Round]; ok {
			score.RoundScore = calculator(sls)
		}
	}
	return nil
}

func (gd *GameData) BankScore() (Score, error) {
	// 1. Add the score of current round to total
	// 2. Switch turn
	return Score{}, nil
}
