import axios from "axios";
import React, { useCallback, useState } from "react";
import { Button, Container, Form, Input } from "semantic-ui-react";

const TotoroCreateGame = () => {
  const [data, setData] = useState({
    title: "",
    description: "",
    home: "",
    homeOdd: 0,
    away: "",
    awayOdd: 0,
    maxRewardAmount: 0,
    betEndDate: 0,
    voteEndDate: 0,
  });

  const onChange = useCallback(
    (e) => {
      const newData = {
        ...data,
        [e.target.name]: e.target.value,
      };
      setData(newData);
    },
    [data]
  );

  const onClick = useCallback(() => {
    const betEndDate = Math.floor(new Date(data.betEndDate).getTime() / 1000);
    const voteEndDate = Math.floor(new Date(data.voteEndDate).getTime() / 1000);
    const homeOdd = parseInt(Number(data.homeOdd) * 100);
    const awayOdd = parseInt(Number(data.awayOdd) * 100);
    const maxRewardAmount = Number(data.maxRewardAmount);
    const newData = {
      ...data,
      betEndDate,
      voteEndDate,
      homeOdd,
      awayOdd,
      maxRewardAmount,
    };
    console.log("send", newData);
    setData(data);

    axios
      .post("/v1/game", JSON.stringify(newData))
      .then((res) => {
        console.log(res.data.data);
        alert("게임 생성 트랜잭션 성공: " +  res.data.data);
      })
      .catch((err) => console.log(err));
  }, [data]);

  return (
    <Container textAlign="left" style={{ padding: 10 }}>
      <Form>
        <Form.Field required>
          <label>주제</label>
          <Input name="title" onChange={onChange} placeholder="title" />
        </Form.Field>
        <Form.Field required>
          <label>설명</label>
          <Input
            name="description"
            onChange={onChange}
            placeholder="description"
          />
        </Form.Field>
        <Form.Field required>
          <label>홈팀</label>
          <Input name="home" onChange={onChange} placeholder="home" />
        </Form.Field>
        <Form.Field>
          <label>홈팀 배당률</label>
          <Input name="homeOdd" onChange={onChange} placeholder="x.xx" />
        </Form.Field>
        <Form.Field required>
          <label>원정팀</label>
          <Input name="away" onChange={onChange} placeholder="away" />
        </Form.Field>
        <Form.Field>
          <label>원정팀 배당률</label>
          <Input name="awayOdd" onChange={onChange} placeholder="x.xx" />
        </Form.Field>
        <Form.Field required>
          <label>최대 지급 가능 상금</label>
          <Input
            type="number"
            name="maxRewardAmount"
            onChange={onChange}
            placeholder="minimum: 1000000000000000000"
          />
        </Form.Field>
        <Form.Field required>
          <label>베팅 마감 일자</label>
          <Input
            name="betEndDate"
            onChange={onChange}
            placeholder="yyyy-mm-ddThh:mm:ss"
          />
        </Form.Field>
        <Form.Field required>
          <label>투표 마감 일자</label>
          <Input
            name="voteEndDate"
            onChange={onChange}
            placeholder="yyyy-mm-ddThh:mm:ss"
          />
        </Form.Field>
        <Form.Field>
          <Button basic onClick={onClick}>
            제출
          </Button>
        </Form.Field>
      </Form>
    </Container>
  );
};

export default TotoroCreateGame;
