import React, { useEffect, useState } from "react";
import {
  Button,
  Confirm,
  Container,
  Grid,
  Header,
  Segment,
} from "semantic-ui-react";
import axios from "axios";

const TotoroGameList = ({ filter }) => {
  const [games, setGames] = useState([]);

  useEffect(() => {
    setGames(null);
    axios
      .get("/v1/game/list?filter=" + filter)
      .then((res) => {
        setGames(res.data.data);
        console.log(res.data.data);
      })
      .catch((err) => console.log(err));
  }, [filter]);

  const betHome = (e, game) => {
    if (filter === "bet") {
      const amount = prompt("얼마나 베팅하시겠습니까?");
      if (amount === null || amount === "0") return;
      axios
        .post(
          "/v1/bet/home",
          JSON.stringify({
            gameId: Number(game.gameId),
            amount: Number(amount),
          })
        )
        .then((res) => {
          alert("베팅 트랜잭션 성공: " + res.data.data);
        })
        .catch((err) => console.log(err));
    } else if (filter === "vote" || filter === "end") {
      if (window.confirm(game.home + "에 투표하시겠습니까?") === false) return;
      axios
        .post(
          "/v1/vote/home",
          JSON.stringify({
            gameId: Number(game.gameId),
          })
        )
        .then((res) => {
          alert("투표 트랜잭션 성공: " + res.data.data);
        })
        .catch((err) => console.log(err));
    }
  };

  const betAway = (e, game) => {
    if (filter === "bet") {
      const amount = prompt("얼마나 베팅하시겠습니까?");
      console.log(amount);
      if (amount === null || amount === "0") return;
      axios
        .post(
          "/v1/bet/away",
          JSON.stringify({
            gameId: Number(game.gameId),
            amount: Number(amount),
          })
        )
        .then((res) => {
          alert("베팅 트랜잭션 성공: " + res.data.data);
        })
        .catch((err) => console.log(err));
    } else if (filter === "vote" || filter === "end") {
      if (window.confirm(game.away + "에 투표하시겠습니까?") === false) return;
      axios
        .post(
          "/v1/vote/away",
          JSON.stringify({
            gameId: Number(game.gameId),
          })
        )
        .then((res) => {
          alert("투표 트랜잭션 성공: " + res.data.data);
        })
        .catch((err) => console.log(err));
    }
  };

  return (
    <>
      {!games ? (
        <div>No list..</div>
      ) : (
        <Container>
          {games.map((game) => (
            <Container>
              <Segment>
                <Grid columns={5} divided style={{ margin: 0, padding: 0 }}>
                  <Grid.Row columns={1} style={{ paddingTop: 0 }}>
                    <Grid.Column>
                      <Header as={"h3"}>{game.title}</Header>
                    </Grid.Column>
                  </Grid.Row>
                  <Grid.Row style={{ margin: 0, padding: 0 }}>
                    <Grid.Column width={5}>
                      <Button
                        basic
                        style={{ margin: 0, padding: 0 }}
                        onClick={(e) => betHome(e, game)}
                      >
                        <Header as={"h4"}>{game.home}</Header>
                      </Button>
                    </Grid.Column>
                    <Grid.Column width={2}>
                      <Header as={"h4"}>{game.homeOdd / 100}</Header>
                    </Grid.Column>
                    <Grid.Column width={2}>
                      <Header as={"h4"}>:</Header>
                    </Grid.Column>
                    <Grid.Column width={2}>
                      <Header as={"h4"}>{game.awayOdd / 100}</Header>
                    </Grid.Column>
                    <Grid.Column width={5}>
                      <Button
                        basic
                        style={{ margin: 0, padding: 0 }}
                        onClick={(e) => betAway(e, game)}
                      >
                        <Header as={"h4"}>{game.away}</Header>
                      </Button>
                    </Grid.Column>
                  </Grid.Row>
                </Grid>
              </Segment>
            </Container>
          ))}
        </Container>
      )}
    </>
  );
};

export default TotoroGameList;
