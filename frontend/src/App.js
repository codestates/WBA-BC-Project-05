import "./App.css";
import "semantic-ui-css/semantic.min.css";
import "semantic-ui-react";
import { Container } from "semantic-ui-react";
import { BrowserRouter, Routes, Route, useNavigate } from "react-router-dom";
import TotoroHeader from "./TotoroHeader";
import TotoroGameList from "./TotoroGameList";
import TotoroCreateGame from "./TotoroCreateGame";
import TotoroMain from "./TotoroMain";
import TotoroBetList from "./TotoroBetList";

function App() {
  return (
    <BrowserRouter>
      <div className="App">
        <Container>
          <TotoroHeader></TotoroHeader>
          <Routes>
            <Route path="/" element={<TotoroMain />}></Route>
            <Route
              path="/games/bet"
              element={<TotoroGameList filter={"bet"} />}
            ></Route>
            <Route
              path="/games/vote"
              element={<TotoroGameList filter={"vote"} />}
            ></Route>
            <Route
              path="/games/end"
              element={<TotoroGameList filter={"end"} />}
            ></Route>
            <Route path="/games/create" element={<TotoroCreateGame />}></Route>

            <Route path="/bets" element={<TotoroBetList />}></Route>
          </Routes>
        </Container>
      </div>
    </BrowserRouter>
  );
}

export default App;
