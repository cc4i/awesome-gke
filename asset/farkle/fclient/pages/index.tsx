import Head from 'next/head'
import Image from 'next/image'
// import styles from '../styles/Home.module.css'
import {Container, Grid, Card, Text, Row, Button, Table, Input, Badge, Avatar, Spacer} from '@nextui-org/react';
import React from 'react';




export default function Home() {
  const points = ['1', '2', '3', '4', '5', '6']
  const [dice1, setDice1] = React.useState('1');
  const [dice2, setDice2] = React.useState('2');
  const [dice3, setDice3] = React.useState('3');
  const [dice4, setDice4] = React.useState('4');
  const [dice5, setDice5] = React.useState('5');
  const [dice6, setDice6] = React.useState('6');
  

  function getRandomInt(min:number , max:number) : number{
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min; 
  }

  function rollDice() {
    console.log("roll the Dice1")
    setDice1(getRandomInt(1,6).toString())
    setDice2(getRandomInt(1,6).toString())
    setDice3(getRandomInt(1,6).toString())
    setDice4(getRandomInt(1,6).toString())
    setDice5(getRandomInt(1,6).toString())
    setDice6(getRandomInt(1,6).toString())
  }
  return (
    <Container fluid css={{ height: '1400px', }}>
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
            Frakle
          </Text>
        </Grid>
        <Grid xs={12}>
            <Grid>
              <Badge
                content=""
                isSquared
                color="primary"
                placement="bottom-right"
                variant="points"
                size="md"
              >
                <Avatar
                  squared
                  size="lg"
                  src="https://i.pravatar.cc/300?u=a042581f4e29026707e"
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
              <Avatar text={dice1} size="xl" squared zoomed/>
            </Grid>
            <Grid>
              <Avatar text={dice2} size="xl" squared/>
            </Grid>
            <Grid>
              <Avatar text={dice3} size="xl" squared/>
            </Grid>
            <Grid>
              <Avatar text={dice4} size="xl" squared/>
            </Grid>
            <Grid>
              <Avatar text={dice5} size="xl" squared/>
            </Grid>
            <Grid>
              <Avatar text={dice6} size="xl" squared/>
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
              <Badge
                content=""
                isSquared
                color="primary"
                placement="bottom-right"
                variant="points"
                size="md"
              >
                <Avatar
                  squared
                  size="lg"
                  src="https://i.pravatar.cc/300?u=a042581f4e29026704f"
                />
              </Badge>
            </Grid>
            <Grid>
              <Input readOnly labelLeft="Score" placeholder="0" />
            </Grid>
        </Grid>
      </Grid.Container>
      <Grid.Container>
      <Table
      aria-label="Static Meld table"
      css={{
        height: "auto",
        minWidth: "100%",
        backgroundColor: "White",
      }}
      bordered
      >
      <Table.Header>
        <Table.Column>Meld</Table.Column>
        <Table.Column>Values</Table.Column>
      </Table.Header>
      <Table.Body>
        <Table.Row key="1">
          <Table.Cell>Ones</Table.Cell>
          <Table.Cell> 100</Table.Cell>
        </Table.Row>
        <Table.Row key="2">
          <Table.Cell>Fives</Table.Cell>
          <Table.Cell> 50</Table.Cell>
        </Table.Row>
        <Table.Row key="3">
          <Table.Cell>Three Ones</Table.Cell>
          <Table.Cell> 1000</Table.Cell>
        </Table.Row>
        <Table.Row key="4">
          <Table.Cell>Three Twos</Table.Cell>
          <Table.Cell> 200</Table.Cell>
        </Table.Row>
        <Table.Row key="5">
          <Table.Cell>Three Threes</Table.Cell>
          <Table.Cell> 300</Table.Cell>
        </Table.Row>
        <Table.Row key="6">
          <Table.Cell>Three Fours</Table.Cell>
          <Table.Cell> 400</Table.Cell>
        </Table.Row>
        <Table.Row key="7">
          <Table.Cell>Three Fives</Table.Cell>
          <Table.Cell> 500</Table.Cell>
        </Table.Row>
        <Table.Row key="8">
          <Table.Cell>Three Sixes</Table.Cell>
          <Table.Cell> 600</Table.Cell>
        </Table.Row>
        <Table.Row key="9">
          <Table.Cell>Four of a kind</Table.Cell>
          <Table.Cell> 1000</Table.Cell>
        </Table.Row>
        <Table.Row key="10">
          <Table.Cell>Five of a kind</Table.Cell>
          <Table.Cell> 2000</Table.Cell>
        </Table.Row>
        <Table.Row key="11">
          <Table.Cell>Six of a kind</Table.Cell>
          <Table.Cell> 3000</Table.Cell>
        </Table.Row>
        <Table.Row key="12">
          <Table.Cell>Three Pairs</Table.Cell>
          <Table.Cell> 1500</Table.Cell>
        </Table.Row>
        <Table.Row key="13">
          <Table.Cell>Run</Table.Cell>
          <Table.Cell> 2500</Table.Cell>
        </Table.Row>
        </Table.Body>
      </Table>
      </Grid.Container>
    </Container>
    
  )
}
