import Head from 'next/head'
import Image from 'next/image'
// import styles from '../styles/Home.module.css'
import {Container, Grid, Card, Text, Row, Button, Table, Input, Badge, Avatar, Spacer, Textarea} from '@nextui-org/react';
import React from 'react';
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
  withPlayer?: boolean
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
  currentPlayerId?: string
  // key -> player id
  playerResults: Map<string,Score>
}

// Main
export default function Home() {

  const [d1, setD1] = React.useState<DiceState>({id: "dice1", isOnBoard: true, value: 1, imageUrl:"/dice1.png"});
  const [d2, setD2] = React.useState<DiceState>({id: "dice2", isOnBoard: true, value: 2, imageUrl:"/dice2.png"});
  const [d3, setD3] = React.useState<DiceState>({id: "dice3", isOnBoard: true, value: 3, imageUrl:"/dice3.png"});
  const [d4, setD4] = React.useState<DiceState>({id: "dice4", isOnBoard: true, value: 4, imageUrl:"/dice4.png"});
  const [d5, setD5] = React.useState<DiceState>({id: "dice5", isOnBoard: true, value: 5, imageUrl:"/dice5.png"});
  const [d6, setD6] = React.useState<DiceState>({id: "dice6", isOnBoard: true, value: 6, imageUrl:"/dice6.png"});
  const [game, setGame] = React.useState<Game>({playerResults: new Map()})
  const [pb, setPb] = React.useState<DiceState[]>([])
  const [pm, setPm] = React.useState<DiceState[]>([])

  // TODO: Initial a game, need to update after integrating with game server
  function initGames() {

    let playerBot = {id: uuidv4(), name: "Robot", avatarUrl: "/robot.png"}
    let playerMan = {id: uuidv4(), name: "Chuan", avatarUrl: "/coolman1.png"}

    game.playerResults.set(playerBot.id, {selectionNumber: 1, selections: new Map(), total: 0})
    game.playerResults.set(playerMan.id, {selectionNumber: 1, selections: new Map(), total: 0})
    game.currentPlayerId = playerMan.id
    let ngm = {
      start: new Date(),
      playerResults: game.playerResults,
      //default starting with playerMan
      currentPlayerId: playerMan.id
    }
    setGame(ngm)
  }

  function reload() {
    Router.reload();
  }

  function getRandomInt(min:number , max:number) : number{
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min; 
  }

  function rollOneDice(ds: DiceState) {
    console.log("what->",ds)
    if (ds.isOnBoard) {
      let val = getRandomInt(1,6)
      var dsx = {
        id: ds.id, 
        isOnBoard: ds.isOnBoard, 
        imageUrl:"/dice"+val+".png",
        value: val
      }
      switch (ds.id) {
        case "dice1":
          setD1(dsx);
          break;
        case "dice2":
          setD2(dsx);
          break;
        case "dice3":
          setD3(dsx);
          break;
        case "dice4":
          setD4(dsx);
          break;
        case "dice5":
          setD5(dsx);
          break;
        case "dice6":
          setD6(dsx);
          break;
      
      }
    }
  }

  function rollDice() {
    console.log("Roll all dices")
    //Before roll dices
    let pid = game.currentPlayerId
    if (pid != undefined) {
      let s = game.playerResults.get(pid)
      if ( s != undefined) {
        s.selectionNumber = s.selections.size + 1
        pm.map((d)=>{
          d.isFixed = true
        })
      }
    }
    setGame(game)
    

    //Dices
    rollOneDice(d1)
    rollOneDice(d2)
    rollOneDice(d3)
    rollOneDice(d4)
    rollOneDice(d5)
    rollOneDice(d6)
  }

  const moveDice = (ds: DiceState) => () => {
    console.log("Move dice => ", ds.id)
    if (ds.isFixed) {
      console.log("Do nothing with fixed dice => ", ds.id)
      return
    }
    // setDice1H((dice1H) => (dice1H === id ? -1 : id))
    let dsx = {
      id: ds.id, 
      isOnBoard: !ds.isOnBoard, 
      withPlayer: !ds.withPlayer,
      imageUrl: ds.imageUrl,
      value: ds.value
    }
    console.log("ds => ",ds)
    console.log("dsx => ",dsx)
    // Handle click dices
    console.log("currentPlayerId =>", game.currentPlayerId)
    if (game.currentPlayerId != undefined) {
      let pScore = game.playerResults.get(game.currentPlayerId)
      console.log("pScore.selectionNumbe=>", pScore?.selectionNumber)
      if (pScore?.selections.size == 0 || pScore?.selections.get(pScore?.selectionNumber)==undefined) {
        console.log("pScore.selections.size ==> 0")
        console.log("push===",dsx)
        let nds:DiceState[] = []
        nds.push(dsx)
        pScore?.selections.set(pScore.selectionNumber, nds)
        console.log(game.playerResults.get(game.currentPlayerId)?.selections)
      } else {
        pScore?.selections.forEach((dss: DiceState[], num: number) => {
          if (pScore?.selectionNumber == num) {
            console.log(num, dss);

            if ( !ds.withPlayer) {
              console.log("push===",dsx)
              dss.push(dsx)
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
        console.log(game.playerResults.get(game.currentPlayerId)?.selections)
      }
      if (pScore!=undefined) {
        let da:DiceState[] = []
        pScore.selections.forEach((x,k)=>{
          x.map((d)=>{
            da.push(d)
          })
          da.push({id: "-", isOnBoard: false, imageUrl: ""})
        })
        console.log("da =>",da)
        setPm(da)
      }
      
    }
    setGame(game)

    
    //

    // Change states of 'isOnBoard && withPlayer'
    switch (ds.id) {
      case "dice1":
        setD1(dsx);
        break;
      case "dice2":
        setD2(dsx);
        break;
      case "dice3":
        setD3(dsx);
        break;
      case "dice4":
        setD4(dsx);
        break;
      case "dice5":
        setD5(dsx);
        break;
      case "dice6":
        setD6(dsx);
        break;
    }

    // Update player's own board

  
  }

  function calculate() {

    //x1=x100
    //x5=x50
    //3*1=1000
    //3*2=200
    //3*3=300
    //3*4=400
    //3*5=500
    //3*6=600
    //4*?=1000
    //5*?=2000
    //6*?=3000
    //3*??=1500
    //1,2,3,4,5,6=2500
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
                />
              </Badge>
            </Grid>
            <Grid>
              <Input readOnly labelLeft="Score" placeholder="0" />
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
                />
              </Badge>
            </Grid>
            <Grid>
              <Input readOnly labelLeft="Score" placeholder="0" />
            </Grid>
            <Grid>
              <Grid.Container>
                {
                  pm.map((d, idx) => {
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
