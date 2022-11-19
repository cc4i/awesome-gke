// The game server follows rules from https://cardgames.io/farkle/, which
// uses gameing frameworks to demonstrate their capbilities on GKE.
//
// :: Game rules ::
//
//	Ones: Any die depicting a one. Worth 100 points each.
//	Fives: Any die depicting a five. Worth 50 points each.
//	Three Ones: A set of three dice depicting a one. worth 1000 points
//	Three Twos: A set of three dice depicting a two. worth 200 points
//	Three Threes: A set of three dice depicting a three. worth 300 points
//	Three Fours: A set of three dice depicting a four. worth 400 points
//	Three Fives: A set of three dice depicting a five. worth 500 points
//	Three Sixes: A set of three dice depicting a six. worth 600 points
//	Four of a kind: Any set of four dice depicting the same value. Worth 1000 points
//	Five of a kind: Any set of five dice depicting the same value. Worth 2000 points
//	Six of a kind: Any set of six dice depicting the same value. Worth 3000 points
//	Three Pairs: Any three sets of two pairs of dice. Includes having a four of a kind plus a pair. Worth 1500 points
//	Run: Six dice in a sequence (1,2,3,4,5,6). Worth 2500 points
package main

import "fmt"

type Meld struct {
}

func main() {
	fmt.Println("Dedicated Game Server for Farkel")
}
