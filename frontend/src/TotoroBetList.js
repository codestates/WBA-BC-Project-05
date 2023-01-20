import React, { useEffect, useState } from "react";
import { Container, Grid, Header, Segment } from "semantic-ui-react";
import axios from "axios";

const TotoroBetList = ({ filter }) => {
  const [bets, setBets] = useState([]);

  useEffect(() => {
    setBets(null);
    axios
      .get("/v1/bet?address=0x077917d175fABAa2f53c1Efed05Ecbd3EeFFBf0f")
      .then((res) => {
        setBets(res.data.data);
        console.log(res.data.data);
      })
      .catch((err) => console.log(err));
  }, []);

  return (
    <>
      {!bets ? (
        <div>No list..</div>
      ) : (
        <Container>
          {bets.map((bet) => (
            <Container>
              <Segment>
                <Grid columns={4} divided style={{ margin: 0, padding: 0 }}>
                  <Grid.Row style={{ margin: 0, padding: 0 }}>
                    <Grid.Column>
                      <Header as={"h3"}>{bet.gameId}</Header>
                    </Grid.Column>
                    <Grid.Column width={5}>
                      <Header>{bet.target === 0 ? "" : ""}</Header>
                    </Grid.Column>
                    <Grid.Column width={2}>
                      <Header as={"h4"}>{bet.target / 100}</Header>
                    </Grid.Column>
                    <Grid.Column width={2}>
                      <Header as={"h4"}>:</Header>
                    </Grid.Column>
                    <Grid.Column width={2}>
                      <Header as={"h4"}>{bet.admount / 100}</Header>
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

export default TotoroBetList;
