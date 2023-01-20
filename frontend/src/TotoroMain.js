import axios from "axios";
import React, { useEffect, useState } from "react";
import { Button, Container, Header } from "semantic-ui-react";

const TotoroMain = () => {
  const [balance, setBalance] = useState(0);
  const [msg, setMsg] = useState("");

  useEffect(() => {
    axios
      .get("/v1/token/balance")
      .then((res) => {
        setBalance(res.data.data);
        console.log(res.data.data);
      })
      .catch((err) => console.log(err));
  }, []);

  const onclick = () => {
    axios
      .get("/v1/token/welcome")
      .then((res) => {
        setBalance(res.data.data)
        console.log(res.data.data);
      })
      .catch((err) => {
        console.log(err);
        setMsg(err.response.data.data);
      });
  };
  return (
    <Container>
      <Header as={"h2"}>환영합니다!</Header>
      <Header>보유 토큰: {balance / 1e18}</Header>
      <Button basic onClick={onclick}>
        환영 토큰 받기
      </Button>
      <Header>{msg}</Header>
    </Container>
  );
};

export default TotoroMain;
