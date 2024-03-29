import Head from 'next/head'
import Image from 'next/image'
// import styles from '../styles/Home.module.css'
import {Container, Grid, Text, Button, Input, Badge, Avatar, Spacer, Textarea, Tooltip, Modal, Row, Checkbox} from '@nextui-org/react';
import React, { Component, useEffect } from 'react';
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

//Picked dices for the player
interface PickedDices {
  who: string
  dices: DiceState[]
}

// ////////declare global/////////
// Six dices 1-6
var Dice1:DiceState = {id: "dice1", isOnBoard: false, value: 1, imageUrl:"/dice1.png"}
var Dice2:DiceState = {id: "dice2", isOnBoard: false, value: 2, imageUrl:"/dice2.png"}
var Dice3:DiceState = {id: "dice3", isOnBoard: false, value: 3, imageUrl:"/dice3.png"}
var Dice4:DiceState = {id: "dice4", isOnBoard: false, value: 4, imageUrl:"/dice4.png"}
var Dice5:DiceState = {id: "dice5", isOnBoard: false, value: 5, imageUrl:"/dice5.png"}
var Dice6:DiceState = {id: "dice6", isOnBoard: false, value: 6, imageUrl:"/dice6.png"}
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
  const [downPlayer, setDownPlayer] = React.useState<Player>(DownPlayer)
  const [onlinePlayers, setOnlinePlayers] = React.useState<string[]>([])
  const [visible, setVisible] = React.useState(false);
  const [upScore, setUpScore] = React.useState<number>(0)
  const [downScore, setDownScore] = React.useState<number>(0)
  const [upTotal, setUpTotal] = React.useState<number>(0)
  const [downTotal, setDownTotal] = React.useState<number>(0)

  const closeHandler = () => {
    setVisible(false);
  };

  const login = () => {
    setVisible(true);
    
  };

  // TODO: Initial a game, need to update after integrating with game server
  function initGames() {
    //Inital players' name
    setDownPlayer(DownPlayer);

    if (DownPlayer.name != "Chuan") {
      closeHandler()
      console.log("upPlayer => ", UpPlayer)
      console.log("downPlayer => ", DownPlayer)
  
      if (UpPlayer != undefined && UpPlayer != undefined) {
        //Setup players
        TheGame.playerResults.set(UpPlayer.id, {selectionNumber: 1, selections: new Map(), total: 0})
        TheGame.playerResults.set(DownPlayer.id, {selectionNumber: 1, selections: new Map(), total: 0})
        TheGame.currentPlayerId = DownPlayer.id
  
        //Turn on dices
        initalDices()
  
        //Check-in status
        TheGame.isStart = true
        TheGame.isEnd = false     
  
        let copyGame = structuredClone(TheGame) 
        setGame(copyGame)
  
        //Set status
        console.log("game.isStart=> ", game.isStart, " game.isEnd=>", game.isEnd)

        //Add player to list
        onlinePlayers.push(DownPlayer.name)
        setOnlinePlayers(onlinePlayers)
      }      
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

  function initalDices() {
    // place dices back && reset values
    for (var x of [Dice1, Dice2, Dice3, Dice4, Dice5, Dice6]) {
      x.isOnBoard=true
      x.withUpPlayer=false
      x.withDownPlayer=false
      x.isFixed=false
      
      let copy = structuredClone(x)
      switch (x.id) {
        case "dice1":
          copy.value=1
          copy.imageUrl="/dice"+copy.value+".png"
          setD1(copy);
          break;
        case "dice2":
          copy.value=2
          copy.imageUrl="/dice"+copy.value+".png"
          setD2(copy);
          break;
        case "dice3":
          copy.value=3
          copy.imageUrl="/dice"+copy.value+".png"
          setD3(copy);
          break;
        case "dice4":
          copy.value=4
          copy.imageUrl="/dice"+copy.value+".png"
          setD4(copy);
          break;
        case "dice5":
          copy.value=5
          copy.imageUrl="/dice"+copy.value+".png"
          setD5(copy);
          break;
        case "dice6":
          copy.value=6
          copy.imageUrl="/dice"+copy.value+".png"
          setD6(copy);
          break;
      }
    }
    
  }

  function switchTurn() {
    console.log("Switch the turn to opponent player...")
    if (TheGame.currentPlayerId != undefined) {
      // reset results
      TheGame.playerResults.forEach((s,k)=>{
        s.selectionNumber=1
        s.selections=new Map()
      })
      
      //reset dices
      initalDices()

      // switch turn
      TheGame.currentPlayerId=(TheGame.currentPlayerId==UpPlayer.id?DownPlayer.id:UpPlayer.id)
      console.log("opponent player is =>", TheGame.currentPlayerId==UpPlayer.id?"Up":"Down", "id => ", TheGame.currentPlayerId)
      let copyGame = structuredClone(TheGame)
      setGame(copyGame)

      setPb({who: uuidv4(), dices:[]})
      setPm({who: uuidv4(), dices:[]})

      setUpScore(calculate(UpPlayer.id,false))
      setDownScore(calculate(DownPlayer.id,false))

    }

  }

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
    
                if ( !x.isOnBoard) {
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

    //Calculate
    setUpScore(calculate(UpPlayer.id, false))
    setDownScore(calculate(DownPlayer.id, false))
    setUpTotal(calculate(UpPlayer.id, true))
    setDownTotal(calculate(DownPlayer.id, true))
  
  }

  // Calculating score
  // All values in array equal to specfic number
  function allEqual(ns:number[], num: number) {
    if (num == -1) {
      return ns.every(val => val === ns[0])
    } else {
      return ns.every(val => val === ns[0]) && ns[0]==num
    }
    
  }

  function calculator(dVs:number[]):number {
    // 3x1 => 1000
    if (dVs.length==3 && allEqual(dVs, 1)) {
      return 1000
    }
    // 3x2 => 200
    if (dVs.length==3 && allEqual(dVs, 2)) {
      return 200
    }
    // 3x3 => 300
    if (dVs.length==3 && allEqual(dVs, 3)) {
      return 300
    }
    // 3x4 => 400
    if (dVs.length==3 && allEqual(dVs, 4)) {
      return 400
    }
    // 3x5 => 500
    if (dVs.length==3 && allEqual(dVs, 5)) {
      return 500
    }
    // 3x6 => 600
    if (dVs.length==3 && allEqual(dVs, 6)) {
      return 600
    }
    // 4x? => 1000
    if (dVs.length==4 && allEqual(dVs, -1)) {
      return 1000
    }
    // 5x? => 2000
    if (dVs.length==5 && allEqual(dVs, -1)) {
      return 2000
    }
    // 6x? => 3000
    if (dVs.length==6 && allEqual(dVs, -1)) {
      return 3000
    }
    dVs.sort(((a, b) => a-b))
    // 3x?? => 1500
    if (dVs.length==6 && allEqual(dVs.slice(0,3), -1) && allEqual(dVs.slice(2,4), -1) && allEqual(dVs.slice(4,6), -1)) {
      return 1500
    }
    // 1,2,3,4,5,6 => 2500
    if (dVs.length==6 && dVs.toString()=="123456") {
      return 3000
    }
    
    // x1 => x100
    // x5 => x50
    if (dVs.length>0) {
      let a:number = 0
      dVs.map((s) => {
        if (s!=1 && s!=5) {
          return 0
        } else {
          if (s==1) {
            a+=100
          }
          if (s==5) {
            a+=50
          }
        }
      })
      return a
    }
  
    return 0
  }

  function calculate(playerId:string, isTotal:boolean):number {
    let s:number = 0
    
    let score = TheGame.playerResults.get(playerId)
    if (score!=undefined) {
      let sn = score.selectionNumber
      score.selections.forEach((ds, k)=>{
        if (k == sn || isTotal) {
          let dv:number[] = []
          ds.map((d)=>{
            if (d.value!=undefined) {
              dv.push(d.value)
            }
          })
          s += calculator(dv)
        }
      })
    }
    
    console.log("score=== ", s)
    return s
  }

  useEffect(()=>{
    console.log("Finally :: pb.length=> ", pb?.dices.length, " pm.length=> ", pm?.dices.length)
    console.log("Finally :: game=> ", game)
  })

  return (
    <Container fluid css={{ height: '1400px', }}>
      <Head>
        <title>Farkle</title>
        <meta name="description" content="Farkle built by CC" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Grid.Container id="mainGame" gap={2} justify="center">
        <Grid id="name" xs={12}>
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
        <Grid id="up" xs={12}>
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
              {/* <Input id='upScore' labelLeft="Score" placeholder="0" value={upScore} /> */}
              <Input bordered label="Score" placeholder="0" color="primary" width="100px" value={upScore} disabled />
              <Input bordered label="Banked" placeholder="0" color="success" width="100px" value={upTotal} disabled />
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
        <Grid key="mainBoard" xs={12} css={{ alignItems: 'center', }}>
          <Grid.Container gap={2} xs={12} css={{
            background: 'green',
            alignContent: 'center',
            alignItems: 'center',
          }}>
            <Grid>
              { d1.isOnBoard && <Avatar src={d1.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d1)} as="button"/> }
            </Grid>
            <Grid>
            { d2.isOnBoard && <Avatar src={d2.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d2)} as="button"/> }
            </Grid>
            <Grid>
            { d3.isOnBoard && <Avatar src={d3.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d3)} as="button"/> }
            </Grid>
            <Grid>
            { d4.isOnBoard && <Avatar src={d4.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d4)} as="button"/> }
            </Grid>
            <Grid>
            { d5.isOnBoard && <Avatar src={d5.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d5)} as="button"/> }
            </Grid>
            <Grid>
            { d6.isOnBoard && <Avatar src={d6.imageUrl} size="xl" squared zoomed bordered onClick={ moveDice(d6)} as="button"/> }
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
                <Button aria-label='bank' css={{
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
        <Grid id="down" xs={12}>
            <Grid>
            <Badge disableOutline placement="bottom-right" content={downPlayer.name}>
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
              {/* <Input id='downScore' labelLeft="Score" placeholder="0" value={downScore}/> */}
              <Input bordered label="Score" placeholder="0" color="primary" width="100px" value={downScore} disabled />
              <Input bordered label="Banked" placeholder="0" color="success" width="100px" value={downTotal} disabled />

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
      <Grid.Container id="controller">
          <Input id='server' labelLeft="Server" placeholder="127.0.0.1:9000" />   
          <Spacer y={2}></Spacer>     
          <Button aria-label='init' css={{
                  background: 'white',
                  color: 'Black',
                  borderColor: 'Beige',
                  fontWeight: '$bold',
                  fontStyle: 'italic'
                }}
                bordered
                ghost
                onPress={login}
          >
            Start Game
          </Button>
          <Spacer y={2}></Spacer>
          <>
            <Modal
              closeButton
              preventClose
              aria-labelledby="modal-title"
              open={visible}
              onClose={closeHandler}
            >
              <Modal.Header>
                <Text id="modal-title" size={18}>
                  Welcome to &nbsp;
                  <Text b i size={18}>
                    Farkle World
                  </Text>
                </Text>
              </Modal.Header>
              <Modal.Body>
                <Input
                  clearable
                  bordered
                  fullWidth
                  color="primary"
                  size="lg"
                  labelLeft="Name"
                  onChange={e=>{ DownPlayer.name=e.target.value; }}
                />
                <Row justify="space-between">
                  <Checkbox>
                    <Text size={14}>Remember me</Text>
                  </Checkbox>
                </Row>
              </Modal.Body>
              <Modal.Footer>
                <Button auto flat color="error" onClick={closeHandler}>
                  Close
                </Button>
                <Button auto onClick={initGames}>
                  Login
                </Button>
              </Modal.Footer>
            </Modal>
          </>
          <Button aria-label='reload' css={{
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
      <Grid.Container id="onlinePlayers">
        <Spacer y={1}></Spacer> 
        <Grid xs={12}>
          <Text size="$2xl" css={{textGradient: "45deg, $blue600 -20%, $pink600 50%", fontStyle: "italic"}}>#Online Players# </Text>
          <Spacer></Spacer> 
          <Avatar.Group count={onlinePlayers.length}>
            {onlinePlayers.map((name, index) => (
              <Tooltip
              key={index}
              color="primary"
              content={"Click to challenge :: [ "+name+" ]"}
              placement="topStart"
              >
                <Avatar key={index} size="lg" pointer text={name} stacked />
              </Tooltip>
            ))}
          </Avatar.Group>
        </Grid>
        <Spacer y={1}></Spacer> 
      </Grid.Container>
      <Grid.Container id="remark">
        <Textarea
          id='rules'
          width='98%'
          size='lg'
          status="warning"
          color="warning"
          helperColor="warning"
          initialValue="..."
          value={
            "::Game's rules::\n\n🚀#1. Ones => 100\n🚀#2. Fives => 50\n🚀#3. Three Ones => 1000\n🚀#4. Three Twos => 200\n🚀#5. Three Threes => 300\n🚀#6. Three Fours => 400\n🚀#7. Three Fives => 500\n🚀#8. Three Sixes => 600\n🚀#9. Four of a kind => 1000\n🚀#10. Five of a kind => 2000\n🚀#11. Six of a kind => 3000\n🚀#12. Three Pairs => 1500\n🚀#13. Run => 2500\n"
          }
          maxRows={20}
          readOnly
        />
      </Grid.Container>
    </Container>
    
  )
}
