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

	"log"

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
	Players []*Player `json:"players,omitempty"`
	// Current player, who's turn
	CurrentPlayer *Player `json:"currentPlayer,omitempty"`
	// Final winner of the game
	FinalWinner *Player `json:"finalWinner,omitempty"`
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
	CurrentPlayer  Player `json:"currentPlayer,omitempty"`
	OpponentPlayer Player `json:"opponentPlayer,omitempty"`

	// Do
	// "login" - Login into a game
	// "initial" - Intial a game with chosen players
	// "roll" - Roll all rollable dices
	// "move" - Move the dice from board to player's area
	// "switch" - Switch the turn between players
	// "bank" - Bank the score and switch the turn of the game
	// "start" - Start a new game
	// "end" - End the current game
	Action string `json:"action"`

	//What
	// Game Id
	GameId string `json:"gameId,omitempty"`
	// dice
	DiceX Dice `json:"diceX,omitempty"`
}

type FarkelInterface interface {
	IntialGame(pname1 string, pname2 string) (Game, error)
	EndGame(gameId string) error
	Login(name string) ([]Player, error)
	// TODO:
	Logout(name string) ([]Player, error)
	//
	SearchPlayer(name string) (Player, error)
	RollDices(gameId string) (Game, error)
	SwitchTurn(gameId string) (Game, error)
	MoveDice(gameId, diceId, playerName string) (Game, error)
	IsMovable(gameId string) bool
	Calculate(gameId string) error
	BankScore(gameId string) (Game, error)
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
		fmt.Println(err)
		return "", err
	}
	var p interface{}
	switch fr.Action {

	case "login":
		p, err = gdata.Login(fr.CurrentPlayer.Name)

	case "initial":
	case "start":
		p, err = gdata.IntialGame(fr.CurrentPlayer.Name, fr.OpponentPlayer.Name)

	case "roll":
		p, err = gdata.RollDices(fr.GameId)

	case "move":
		p, err = gdata.MoveDice(fr.GameId, fr.DiceX.Id, fr.CurrentPlayer.Name)

	case "switch":
		p, err = gdata.SwitchTurn(fr.GameId)

	case "bank":
		p, err = gdata.BankScore(fr.GameId)

	case "end":
		err = gdata.EndGame(fr.GameId)
	}

	if err != nil {
		return "", err
	}
	b, _ := json.Marshal(p)
	return string(b), nil

}

// Intial a new game & save it into GameData
func (gd *GameData) IntialGame(pname1 string, pname2 string) (Game, error) {
	nGame := &Game{
		Id:     uuid.New().String(),
		Start:  time.Now().UnixMilli(),
		Scores: make(map[string]*Score),
		Dices:  make(map[string]*Dice),
	}
	if pname1 == pname2 {
		return *nGame, fmt.Errorf("%s, please do not play with yourself :)", pname1)
	}
	player1, ok1 := gd.Players[pname1]
	player2, ok2 := gd.Players[pname2]
	log.Printf("player1 -> %v, player2 -> %v", player1, player2)
	if ok1 && ok2 && player1.Online && player2.Online {
		// Inital players
		nGame.Players = append(nGame.Players, player1)
		nGame.Players = append(nGame.Players, player2)
		nGame.CurrentPlayer = player1
		// Inital score for players
		nGame.Scores[pname1] = &Score{Round: 1}
		nGame.Scores[pname2] = &Score{Round: 1}
		// Intial dices
		for i := 1; i <= 6; i++ {
			nGame.Dices["dice"+strconv.Itoa(i)] = &Dice{
				Id:      "dice" + strconv.Itoa(i),
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

// End current game
func (gd *GameData) EndGame(gameId string) error {
	if game, ok := gd.Games[gameId]; ok {
		game.End = time.Now().UnixMilli()

		// who's the winner
		winner := ""
		winnerScore := 0
		for n, s := range game.Scores {

			if s.BankedScore > winnerScore {
				winner = n
				winnerScore = s.BankedScore
			}

		}
		game.FinalWinner = gd.Players[winner]
		return nil
	}
	return fmt.Errorf("%s is invalid game id", gameId)
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
func (gd *GameData) RollDices(gameId string) (Game, error) {

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
			return gd.SwitchTurn(gameId)
		}
		return *game, nil
	}
	return Game{}, fmt.Errorf("%s is invalid game id", gameId)
}

func (gd *GameData) SwitchTurn(gameId string) (Game, error) {
	if game, ok := gd.Games[gameId]; ok {
		//3.1 Reset RoundScore=0
		score := game.Scores[game.CurrentPlayer.Name]
		score.RoundScore = 0

		//3.2 Move to next round
		score.Round++

		// 3.3 Swith to other player
		for _, p := range game.Players {
			if p.Name != game.CurrentPlayer.Name {
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
		return *game, nil
	}
	return Game{}, fmt.Errorf("%s is invalid game id", gameId)
}

// Move specific dice into player's board or other way around.
// Id is dice indicator, which is one of following: dice1, dice2, dice3, dice4, dice5, dice6
func (gd *GameData) MoveDice(gameId, diceId, playerName string) (Game, error) {

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

// Score calculator
func calculator(ns []int) int {
	rt := 0

	// 3x1 => 1000
	if len(ns) == 3 && allEqualTo(ns, 1) {
		rt += 1000
	}
	// 3x2 => 200
	if len(ns) == 3 && allEqualTo(ns, 2) {
		rt += 200
	}
	// 3x3 => 300
	if len(ns) == 3 && allEqualTo(ns, 3) {
		rt += 300
	}
	// 3x4 => 400
	if len(ns) == 3 && allEqualTo(ns, 4) {
		rt += 400
	}
	// 3x5 => 500
	if len(ns) == 3 && allEqualTo(ns, 5) {
		rt += 500
	}
	// 3x6 => 600
	if len(ns) == 3 && allEqualTo(ns, 6) {
		rt += 600
	}
	// 4x? => 1000
	if len(ns) == 4 && allEqualTo(ns, -1) {
		rt += 1000
	}
	// 5x? => 2000
	if len(ns) == 5 && allEqualTo(ns, -1) {
		rt += 2000
	}
	// 6x? => 3000
	if len(ns) == 6 && allEqualTo(ns, -1) {
		rt += 3000
	}

	// 3x?? => 1500
	sort.Ints(ns)
	if len(ns) == 6 && allEqualTo(ns[:3], -1) && allEqualTo(ns[2:4], -1) && allEqualTo(ns[4:6], -1) {
		rt += 1500
	}
	// 1,2,3,4,5,6 => 2500
	if len(ns) == 6 && strings.Trim(strings.Replace(fmt.Sprint(ns), " ", ",", -1), "[]") == "123456" {
		rt += 3000
	}

	// x1 => x100
	// x5 => x50
	if len(ns) > 0 {
		for _, n := range ns {

			if n == 1 {
				rt += 100
			}
			if n == 5 {
				rt += 50
			}
		}
	}

	return rt
}

// Calculate the score for curent round
func (gd *GameData) Calculate(gameId string) error {
	if game, ok := gd.Games[gameId]; ok {
		if score, ok := game.Scores[game.CurrentPlayer.Name]; ok {
			if sls, ok := score.Selections[score.Round]; ok {
				score.RoundScore = calculator(sls)
			}
		}
	}
	return nil
}

// Bank the score
func (gd *GameData) BankScore(gameId string) (Game, error) {

	if game, ok := gd.Games[gameId]; ok {
		// 1. Add the score of current round to total
		if score, ok := game.Scores[game.CurrentPlayer.Name]; ok {
			score.BankedScore += score.RoundScore
			score.RoundScore = 0
		}

		// 2. Switch turn
		_, err := gd.SwitchTurn(gameId)
		return *game, err
	}

	return Game{}, fmt.Errorf("%s is invalid game id", gameId)
}
