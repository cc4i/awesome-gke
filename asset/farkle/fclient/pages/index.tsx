import Head from 'next/head'
import Image from 'next/image'
// import styles from '../styles/Home.module.css'
import {Container, Grid, Card, Text, Row, Button, Table, Input, Badge, Avatar, Spacer, Textarea} from '@nextui-org/react';
import React, { useEffect } from 'react';
import {v4 as uuidv4} from 'uuid';
import Router from 'next/router'


// Player of game
interface Player {
  id: string
  name: string
  avatarUrl: string
}

// Status of dices
interface DiceState {
  id: string
  value?: number

  withWho?: Player
  withUpPlayer?: boolean
  withDownPlayer?: boolean
  isFixed?: boolean

  isOnBoard: boolean
  imageUrl: string
}

// The result for each play in a Game
interface Score {
  selectionNumber: number
  // key -> selection number
  selections: Map<number,DiceState[]>
  total: number
}

// Each game 
interface Game {
  start?: Date
  end?: Date
  isStart?: boolean
  isEnd?: boolean
  currentPlayerId?: string
  // key -> player id
  playerResults: Map<string,Score>
}

interface PickedDices {
  who: string
  dices: DiceState[]
}

// ////////declare global/////////
// Six dices 1-6
var Dice1:DiceState = {id: "dice1", isOnBoard: true, value: 1, imageUrl:"/dice1.png"}
var Dice2:DiceState = {id: "dice2", isOnBoard: true, value: 2, imageUrl:"/dice2.png"}
var Dice3:DiceState = {id: "dice3", isOnBoard: true, value: 3, imageUrl:"/dice3.png"}
var Dice4:DiceState = {id: "dice4", isOnBoard: true, value: 4, imageUrl:"/dice4.png"}
var Dice5:DiceState = {id: "dice5", isOnBoard: true, value: 5, imageUrl:"/dice5.png"}
var Dice6:DiceState = {id: "dice6", isOnBoard: true, value: 6, imageUrl:"/dice6.png"}
// The game
var TheGame:Game = {playerResults: new Map()}
// 2 Players
var UpPlayer:Player = {id: uuidv4(), name: "Robot", avatarUrl: "/robot.png"}
var DownPlayer:Player = {id: uuidv4(), name: "Chuan", avatarUrl: "/coolman1.png"}
//

// Main
export default function Home() {
  const [d1, setD1] = React.useState<DiceState>(Dice1);
  const [d2, setD2] = React.useState<DiceState>(Dice2);
  const [d3, setD3] = React.useState<DiceState>(Dice3);
  const [d4, setD4] = React.useState<DiceState>(Dice4);
  const [d5, setD5] = React.useState<DiceState>(Dice5);
  const [d6, setD6] = React.useState<DiceState>(Dice6);
  const [game, setGame] = React.useState<Game>(TheGame)
  const [pb, setPb] = React.useState<PickedDices>()
  const [pm, setPm] = React.useState<PickedDices>()
  const [upPlayer, setUpPlayer] = React.useState<Player>(UpPlayer)
  const [downPlayer, setDownUpPlayer] = React.useState<Player>(DownPlayer)

  // TODO: Initial a game, need to update after integrating with game server
  function initGames() {

    console.log("upPlayer => ", UpPlayer)
    console.log("downPlayer => ", DownPlayer)

    if (UpPlayer != undefined && UpPlayer != undefined) {
      TheGame.playerResults.set(UpPlayer.id, {selectionNumber: 1, selections: new Map(), total: 0})
      TheGame.playerResults.set(DownPlayer.id, {selectionNumber: 1, selections: new Map(), total: 0})
      TheGame.currentPlayerId = DownPlayer.id

      TheGame.isStart = true
      TheGame.isEnd = false     
      let copyGame = structuredClone(TheGame) 
      setGame(copyGame)

      //Set status
      console.log("game.isStart=> ", game.isStart, " game.isEnd=>", game.isEnd)
    }

  }

  function reload() {
    Router.reload();
    console.log("game.isStart=> ", game.isStart, " game.isEnd=>", game.isEnd)

  }

  function getRandomInt(min:number , max:number) : number{
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min; 

  }

  function switchTurn() {
    console.log("Switch the turn to opponent player...")
    if (TheGame.currentPlayerId != undefined) {
      // reset dices
      TheGame.playerResults.forEach((s,k)=>{
        s.selectionNumber=1
        s.selections=new Map()
      })
      
      // place dices back
      for (var x of [Dice1, Dice2, Dice3, Dice4, Dice5, Dice6]) {
        x.isOnBoard=true
        x.withUpPlayer=false
        x.withDownPlayer=false
        x.isFixed=false
        let copy = structuredClone(x)
        switch (x.id) {
          case "dice1":
            setD1(copy);
            break;
          case "dice2":
            setD2(copy);
            break;
          case "dice3":
            setD3(copy);
            break;
          case "dice4":
            setD4(copy);
            break;
          case "dice5":
            setD5(copy);
            break;
          case "dice6":
            setD6(copy);
            break;
        }
      }
      // switch turn
      TheGame.currentPlayerId=(TheGame.currentPlayerId==UpPlayer.id?DownPlayer.id:UpPlayer.id)
      console.log("opponent player is =>", TheGame.currentPlayerId==UpPlayer.id?"Up":"Down", "id => ", TheGame.currentPlayerId)
      let copyGame = structuredClone(TheGame)
      setGame(copyGame)

      setPb({who: uuidv4(), dices:[]})
      setPm({who: uuidv4(), dices:[]})

    }

  }

  useEffect(()=>{
    console.log("Finally :: pb.length=> ", pb?.dices.length, " pm.length=> ", pm?.dices.length)
    console.log("Finally :: game=> ", game)
  })

  function rollOneDice(ds: DiceState) {
    
    for (var x of [Dice1, Dice2, Dice3, Dice4, Dice5, Dice6]) {
      if (x.id == ds.id && x.isOnBoard) {
        let ri = getRandomInt(1,6)
        x.imageUrl = "/dice"+ri+".png"
        x.value = ri
  
        let copy = structuredClone(x)
        switch (x.id) {
          case "dice1":
            setD1(copy);
            break;
          case "dice2":
            setD2(copy);
            break;
          case "dice3":
            setD3(copy);
            break;
          case "dice4":
            setD4(copy);
            break;
          case "dice5":
            setD5(copy);
            break;
          case "dice6":
            setD6(copy);
            break;
        }
        console.log("Roll the dice => ",x)
      }
    }
    
  }

  function rollDice() {
    console.log("player is =>", TheGame.currentPlayerId==UpPlayer.id?"Up":"Down", "id => ", TheGame.currentPlayerId)
    console.log("Roll all dices")
    //Before roll dices
    if (TheGame.currentPlayerId != undefined) {
      let s = TheGame.playerResults.get(TheGame.currentPlayerId)
      if ( s != undefined) {
        s.selectionNumber = s.selections.size + 1
        pm?.dices.map((d)=>{
          for (var x of [Dice1, Dice2, Dice3, Dice4, Dice5, Dice6]) {
            if (d.id == x.id) {
              x.isFixed = true
            }
          }
          d.isFixed = true
        })
      }
    }
    
    // Roll dices
    for (var x of [Dice1, Dice2, Dice3, Dice4, Dice5, Dice6]) {
      if (x.isOnBoard) {
        rollOneDice(x)
      }
    }
    

    // Check if it's turn for other player
    let swt = true
    for (var x of [Dice1, Dice2, Dice3, Dice4, Dice5, Dice6]) {
        if (x.isOnBoard) {
          if (x.value==1 || x.value==5) {
            swt = false
          }
        }
      }
    if (swt) {
      switchTurn()
    }

  }

  const moveDice = (ds: DiceState) => () => {
    console.log("Move dice => ", ds)
    if (!game.isStart) {
      console.log("Game is not started, do nothing!")
      return
    }
    if (ds.isFixed) {
      console.log("Do nothing with fixed dice => ", ds.id)
      return
    }

    // Handle click dices
    for (var x of [Dice1, Dice2, Dice3, Dice4, Dice5, Dice6]) {
      if (x.id == ds.id) {
        ///
        x.isOnBoard= !x.isOnBoard
        if (TheGame.currentPlayerId==UpPlayer.id) {
          x.withUpPlayer=!x.withUpPlayer
        } else {
          x.withDownPlayer=!x.withDownPlayer
        }
        ///
        console.log("TheGame.currentPlayerId =>", TheGame.currentPlayerId)
        if (TheGame.currentPlayerId != undefined) {
          let pScore = TheGame.playerResults.get(TheGame.currentPlayerId)
          console.log("pScore.selectionNumbe =>", pScore?.selectionNumber)
          if (pScore?.selections.size == 0 || pScore?.selections.get(pScore?.selectionNumber)==undefined) {
            console.log("pScore.selections.size ==> 0")
            console.log("push===",x)
            let nds:DiceState[] = []
            nds.push(x)
            pScore?.selections.set(pScore.selectionNumber, nds)
            console.log(TheGame.playerResults.get(TheGame.currentPlayerId)?.selections)
          } else {
            pScore?.selections.forEach((dss: DiceState[], num: number) => {
              if (pScore?.selectionNumber == num) {
                console.log(num, dss);
    
                if ( !ds.withUpPlayer || !ds.withDownPlayer) {
                  console.log("push===",x)
                  dss.push(x)
                } else {
                  dss.forEach((d, idx)=>{
                    if (d.id == ds.id) {
                      console.log("splice===",d)
                      dss.splice(idx, 1);
                    }
                  })
                }
    
              }
            })
            console.log(TheGame.playerResults.get(TheGame.currentPlayerId)?.selections)
          }
          if (pScore!=undefined) {
            let da:PickedDices = {who: uuidv4(), dices:[]}
            pScore.selections.forEach((x,k)=>{
              x.map((d)=>{
                da.dices.push(d)
              })
              da.dices.push({id: "-", isOnBoard: false, imageUrl: ""})
            })
            console.log("da =>",da)
            if (TheGame.currentPlayerId == upPlayer?.id) {
              setPb(da)
            } else {
              setPm(da)
            }
            
          }
          
        }

        // Change states of 'isOnBoard && withPlayer'
        let copy = structuredClone(x)
        switch (ds.id) {
          case "dice1":
            setD1(copy);
            break;
          case "dice2":
            setD2(copy);
            break;
          case "dice3":
            setD3(copy);
            break;
          case "dice4":
            setD4(copy);
            break;
          case "dice5":
            setD5(copy);
            break;
          case "dice6":
            setD6(copy);
            break;
        }

      }
    }
    let copyGame = structuredClone(TheGame)
    setGame(copyGame)
    //
  
  }

  // TODO: calculating score
  function calculate() {

  }

  
  return (
    <Container fluid css={{ height: '1400px', }}>
      <Head>
        <title>Farkle</title>
        <meta name="description" content="Farkle built by CC" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Grid.Container gap={2} justify="center">
        <Grid xs={12}>
        <Text
            h1
            size={60}
            css={{
              textGradient: "45deg, $yellow600 -20%, $red600 100%",
            }}
            weight="bold"
          >
            Farkle
          </Text>
          <Text
            h1
            size={12}
            css={{
              textGradient: "45deg, $yellow600 -20%, $red600 100%",
            }}
            weight="bold"
          >
            by CC
          </Text>
        </Grid>
        <Grid xs={12}>
            <Grid>
              <Badge disableOutline placement="bottom-right" content="Robot">
                <Avatar
                  squared
                  size="xl"
                  src="/robot.png"
                  bordered={game.currentPlayerId==upPlayer.id}
                  borderWeight="extrabold"
                  color="gradient"
                />
              </Badge>
            </Grid>
            <Grid>
              <Input readOnly labelLeft="Score" placeholder="0" />
            </Grid>
            <Grid>
              <Grid.Container>
                {
                  pb?.dices.map((d, idx) => {
                    if (d.id == "-") {
                      return <Avatar key={uuidv4()} src={d.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d)} /> && <Spacer key={uuidv4()} y={2}></Spacer>
                    } else {
                      return <Avatar key={uuidv4()} src={d.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d)} />
                    }
                    
                  })
                }
                      
              </Grid.Container> 
            </Grid>
        </Grid>
        <Grid xs={12} css={{ alignItems: 'center', }}>
          <Grid.Container gap={2} xs={12} css={{
            background: 'green',
            alignContent: 'center',
            alignItems: 'center',
          }}>
            <Grid>
              { d1.isOnBoard && <Avatar src={d1.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d1)} /> }
            </Grid>
            <Grid>
            { d2.isOnBoard && <Avatar src={d2.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d2)} /> }
            </Grid>
            <Grid>
            { d3.isOnBoard && <Avatar src={d3.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d3)} /> }
            </Grid>
            <Grid>
            { d4.isOnBoard && <Avatar src={d4.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d4)} /> }
            </Grid>
            <Grid>
            { d5.isOnBoard && <Avatar src={d5.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d5)} /> }
            </Grid>
            <Grid>
            { d6.isOnBoard && <Avatar src={d6.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d6)} /> }
            </Grid>
            <Spacer y={10} />

            <Grid.Container gap={2}>
              <Grid>
                <Button aria-label='roll' css={{
                  background: 'white',
                  color: 'Black',
                  borderColor: 'Beige',
                  fontWeight: '$bold',
                  fontStyle: 'italic'
                }}
                bordered
                ghost
                onPress={rollDice}
                disabled={!game.isStart}
                >Roll Dice</Button>
              </Grid>
              <Grid>
                <Button aria-label="bank" css={{
                  background: 'white',
                  color: 'Black',
                  borderColor: 'Beige',
                  fontWeight: '$bold',
                  fontStyle: 'italic'
                }}
                bordered
                ghost
                disabled={!game.isStart}
                >Bank Score</Button>
              </Grid>
            </Grid.Container>
          </Grid.Container>
        </Grid>       
        <Grid xs={12}>
            <Grid>
            <Badge disableOutline placement="bottom-right" content="Chuan">
                <Avatar
                  squared
                  size="xl"
                  src="/coolman1.png"
                  bordered={game.currentPlayerId==downPlayer.id}
                  borderWeight="extrabold"
                  color="gradient"
                />
              </Badge>
            </Grid>
            <Grid>
              <Input readOnly labelLeft="Score" placeholder="0" />
            </Grid>
            <Grid>
              <Grid.Container>
                {
                  pm?.dices.map((d, idx) => {
                    if (d.id == "-") {
                      return <Avatar key={uuidv4()} src={d.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d)} /> && <Spacer key={uuidv4()} y={2}></Spacer>
                    } else {
                      return <Avatar key={uuidv4()} src={d.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d)} />
                    }
                    
                  })
                }
                      
              </Grid.Container> 
            </Grid>
        </Grid>
      </Grid.Container>
      <Grid.Container>
          <Input labelLeft="Server" placeholder="127.0.0.1:9000" />   
          <Spacer y={2}></Spacer>     
          <Button aria-label='roll' css={{
                  background: 'white',
                  color: 'Black',
                  borderColor: 'Beige',
                  fontWeight: '$bold',
                  fontStyle: 'italic'
                }}
                bordered
                ghost
                onPress={initGames}
          >
            Start Game
          </Button>
          <Spacer y={2}></Spacer>     
          <Button aria-label='roll' css={{
                  background: 'white',
                  color: 'Black',
                  borderColor: 'Beige',
                  fontWeight: '$bold',
                  fontStyle: 'italic'
                }}
                bordered
                ghost
                onPress={reload}
          >
            Reload
          </Button>
      </Grid.Container>
      <Grid.Container>
        <Textarea
          width='96%'
          size='lg'
          status="success"
          helperColor="success"
          initialValue="..."
          placeholder="Description"
          label="Game's rules"
          value={
            "🚀#1. Ones => 100\n🚀#2. Fives => 50\n🚀#3. Three Ones => 1000\n🚀#4. Three Twos => 200\n🚀#5. Three Threes => 300\n🚀#6. Three Fours => 400\n🚀#7. Three Fives => 500\n🚀#8. Three Sixes => 600\n🚀#9. Four of a kind => 1000\n🚀#10. Five of a kind => 2000\n🚀#11. Six of a kind => 3000\n🚀#12. Three Pairs => 1500\n🚀#13. Run => 2500\n"
          }
          rows={13}
          readOnly
        />
      </Grid.Container>
    </Container>
    
  )
}
