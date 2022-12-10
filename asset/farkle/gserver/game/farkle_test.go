package game

import (
	"container/list"
	"reflect"
	"testing"
)

func TestFarkleHandler(t *testing.T) {
	type args struct {
		txt string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "login_cc1",
			args: args{
				txt: `{"currentPlayer": {"name": "cc1"}, "action": "login"}`},
			want:    `[{"name":"cc1","online":true}]`,
			wantErr: false,
		},
		{
			name: "login_cc2_after_cc1",
			args: args{
				txt: `{"currentPlayer": {"name": "cc2"}, "action": "login"}`},
			want:    `[{"name":"cc1","online":true},{"name":"cc2","online":true}]`,
			wantErr: false,
		},
		{
			name: "initial_game_for_cc1_&_cc2",
			args: args{
				txt: `{"currentPlayer": {"name": "cc1","online":true}, "opponentPlayer":{"name": "cc2","online":true}, "action": "initial", "gameId": "06c32c52-59dc-4d4b-aa3f-323cb8a0403c", "start": 1670666709239}`,
			},
			want:    `{"id":"06c32c52-59dc-4d4b-aa3f-323cb8a0403c","start":1670666709239,"players":[{"name":"cc1","online":true},{"name":"cc2","online":true}],"currentPlayer":{"name":"cc1","online":true},"scores":{"cc1":{"round":1},"cc2":{"round":1}},"dices":{"dice1":{"id":"dice1","value":1,"onBoard":true},"dice2":{"id":"dice2","value":2,"onBoard":true},"dice3":{"id":"dice3","value":3,"onBoard":true},"dice4":{"id":"dice4","value":4,"onBoard":true},"dice5":{"id":"dice5","value":5,"onBoard":true},"dice6":{"id":"dice6","value":6,"onBoard":true}}}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FarkleHandler(tt.args.txt)
			if (err != nil) != tt.wantErr {
				t.Errorf("FarkleHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FarkleHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameData_IntialGame(t *testing.T) {
	type fields struct {
		Players     map[string]*Player
		Games       map[string]*Game
		LeaderBoard *list.List
	}
	type args struct {
		pname1 string
		pname2 string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Game
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd := &GameData{
				Players:     tt.fields.Players,
				Games:       tt.fields.Games,
				LeaderBoard: tt.fields.LeaderBoard,
			}
			got, err := gd.IntialGame("06c32c52-59dc-4d4b-aa3f-323cb8a0403c", 1670667289521, tt.args.pname1, tt.args.pname2)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameData.IntialGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameData.IntialGame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameData_EndGame(t *testing.T) {
	type fields struct {
		Players     map[string]*Player
		Games       map[string]*Game
		LeaderBoard *list.List
	}
	type args struct {
		gameId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd := &GameData{
				Players:     tt.fields.Players,
				Games:       tt.fields.Games,
				LeaderBoard: tt.fields.LeaderBoard,
			}
			if err := gd.EndGame(tt.args.gameId); (err != nil) != tt.wantErr {
				t.Errorf("GameData.EndGame() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGameData_Login(t *testing.T) {
	type fields struct {
		Players     map[string]*Player
		Games       map[string]*Game
		LeaderBoard *list.List
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Player
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd := &GameData{
				Players:     tt.fields.Players,
				Games:       tt.fields.Games,
				LeaderBoard: tt.fields.LeaderBoard,
			}
			got, err := gd.Login(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameData.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameData.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameData_SearchPlayer(t *testing.T) {
	type fields struct {
		Players     map[string]*Player
		Games       map[string]*Game
		LeaderBoard *list.List
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Player
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd := &GameData{
				Players:     tt.fields.Players,
				Games:       tt.fields.Games,
				LeaderBoard: tt.fields.LeaderBoard,
			}
			got, err := gd.SearchPlayer(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameData.SearchPlayer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameData.SearchPlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameData_RollDices(t *testing.T) {
	type fields struct {
		Players     map[string]*Player
		Games       map[string]*Game
		LeaderBoard *list.List
	}
	type args struct {
		gameId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Game
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd := &GameData{
				Players:     tt.fields.Players,
				Games:       tt.fields.Games,
				LeaderBoard: tt.fields.LeaderBoard,
			}
			got, err := gd.RollDices(tt.args.gameId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameData.RollDices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameData.RollDices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameData_SwitchTurn(t *testing.T) {
	type fields struct {
		Players     map[string]*Player
		Games       map[string]*Game
		LeaderBoard *list.List
	}
	type args struct {
		gameId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Game
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd := &GameData{
				Players:     tt.fields.Players,
				Games:       tt.fields.Games,
				LeaderBoard: tt.fields.LeaderBoard,
			}
			got, err := gd.SwitchTurn(tt.args.gameId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameData.SwitchTurn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameData.SwitchTurn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameData_MoveDice(t *testing.T) {
	type fields struct {
		Players     map[string]*Player
		Games       map[string]*Game
		LeaderBoard *list.List
	}
	type args struct {
		gameId     string
		diceId     string
		playerName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Game
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd := &GameData{
				Players:     tt.fields.Players,
				Games:       tt.fields.Games,
				LeaderBoard: tt.fields.LeaderBoard,
			}
			got, err := gd.MoveDice(tt.args.gameId, tt.args.diceId, tt.args.playerName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameData.MoveDice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameData.MoveDice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameData_IsMovable(t *testing.T) {
	type fields struct {
		Players     map[string]*Player
		Games       map[string]*Game
		LeaderBoard *list.List
	}
	type args struct {
		gameId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd := &GameData{
				Players:     tt.fields.Players,
				Games:       tt.fields.Games,
				LeaderBoard: tt.fields.LeaderBoard,
			}
			if got := gd.IsMovable(tt.args.gameId); got != tt.want {
				t.Errorf("GameData.IsMovable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_allEqualTo(t *testing.T) {
	type args struct {
		ns  []int
		num int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allEqualTo(tt.args.ns, tt.args.num); got != tt.want {
				t.Errorf("allEqualTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculator(t *testing.T) {
	type args struct {
		ns []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculator(tt.args.ns); got != tt.want {
				t.Errorf("calculator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameData_Calculate(t *testing.T) {
	type fields struct {
		Players     map[string]*Player
		Games       map[string]*Game
		LeaderBoard *list.List
	}
	type args struct {
		gameId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd := &GameData{
				Players:     tt.fields.Players,
				Games:       tt.fields.Games,
				LeaderBoard: tt.fields.LeaderBoard,
			}
			if err := gd.Calculate(tt.args.gameId); (err != nil) != tt.wantErr {
				t.Errorf("GameData.Calculate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGameData_BankScore(t *testing.T) {
	type fields struct {
		Players     map[string]*Player
		Games       map[string]*Game
		LeaderBoard *list.List
	}
	type args struct {
		gameId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Game
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd := &GameData{
				Players:     tt.fields.Players,
				Games:       tt.fields.Games,
				LeaderBoard: tt.fields.LeaderBoard,
			}
			got, err := gd.BankScore(tt.args.gameId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameData.BankScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameData.BankScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
