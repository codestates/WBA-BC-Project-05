// SPDX-License-Identifier: MIT

pragma solidity ^0.8.17;

import "./TotoroGame.sol";

contract TotoroBet is TotoroGame {
    // 베팅 타겟 상수
    uint8 constant BET_HOME = 0; // 홈팀 승리
    uint8 constant BET_AWAY = 1; // 원정팀 승리

    struct Bet {
        uint betId;
        uint gameId;
        uint amount;
        address bettor;
        uint8 target;
        bool hit;
    }

    event EvBet(Bet bet);

    Bet[] bets;
    // gameId => Bet 매핑
    mapping (uint => uint[]) gameIdbetIds;
    // betId => gameId 매핑
    mapping (uint => uint) betIdGameId;
    // betId => 계정 주소 매핑
    mapping (uint => address) betOwner;
    // 계정 주소 => betId[] 매핑
    mapping (address => uint[]) ownerBets;
    // gameId => 홈 승리에 베팅한 계정 주소 리스트 매핑
    mapping (uint => address[]) betHomeWinBettors;
    // gameId => 원정 승리에 베팅한 계정 주소 리스트 매핑
    mapping (uint => address[]) betAwayWinBettors;

    modifier betValidCheck(uint _gameId, uint _amount) {
        // 테스트 코드
        balanceOf[msg.sender] += 1000000;
        uint32 currentTime = uint32(block.timestamp);
        // 베터의 잔액 체크
        require(balanceOf[msg.sender] >= _amount, "Not enough balance");
        // 게임 아이디 유효성 체크
        require(games.length - 1 >= _gameId, "Invalid gameId");
        // 베팅 마감 날짜 체크
        require(currentTime < games[_gameId].betEndDate, "Aready end game");
        _;
    }

    function betHome(uint _gameId, uint _amount) external betValidCheck(_gameId, _amount) returns (bool) {
        // 현재의 베팅으로 인해 게임 생성자의 최대 상금을 초과하는지 체크
        uint odd = games[_gameId].homeOdd;
        uint maxAccReward = games[_gameId].maxRewardHomeAcc;
        uint winReward = odd * _amount;
        require(games[_gameId].maxRewardAmount >= winReward + maxAccReward, "Exceeding the maximum prize amount");

        // 베팅 생성
        uint newBetId = bets.length;
        bets.push(Bet(newBetId, _gameId, _amount, msg.sender, BET_HOME, false));
        // 게임 아이디 => 베팅 아이디 매핑
        gameIdbetIds[_gameId].push(newBetId); 
        // 베팅 아이디 -> 게임 아이디 매핑
        betIdGameId[newBetId] = _gameId;
        // 베팅 아이디 => 베팅자 매핑
        betOwner[newBetId] = msg.sender;
        // 베팅자 => 베팅 아이디 추가
        ownerBets[msg.sender].push(newBetId);
        // 베터의 잔액 차감
        balanceOf[msg.sender] -= _amount;
        // 게임의 누적 베팅 금액 증가
        games[_gameId].maxRewardHomeAcc += winReward;
        // 홈 승리 베터 리스트에 추가
        betHomeWinBettors[_gameId].push(msg.sender);
        // 베팅 성공 이벤트
        emit EvBet(bets[newBetId]);

        return true;
    }

    function betAway(uint _gameId, uint _amount) external betValidCheck(_gameId, _amount) returns (bool) {
        // 현재의 베팅으로 인해 게임 생성자의 최대 상금을 초과하는지 체크
        uint odd = games[_gameId].awayOdd;
        uint maxAccReward = games[_gameId].maxRewardAwayAcc;
        uint winReward = odd * _amount;
        require(games[_gameId].maxRewardAmount >= winReward + maxAccReward, "Exceeding the maximum prize amount");

        // 베팅 생성
        uint newBetId = bets.length;
        bets.push(Bet(newBetId, _gameId, _amount, msg.sender, BET_AWAY, false));
        // 게임 아이디 => 베팅 아이디 매핑
        gameIdbetIds[_gameId].push(newBetId); 
        // 베팅 아이디 -> 게임 아이디 매핑
        betIdGameId[newBetId] = _gameId;
        // 베팅 아이디 => 베팅자 매핑
        betOwner[newBetId] = msg.sender;
        // 베팅자 => 베팅 아이디 추가
        ownerBets[msg.sender].push(newBetId);
        // 베터의 잔액 차감
        balanceOf[msg.sender] -= _amount;
        // 게임의 누적 베팅 금액 증가
        games[_gameId].maxRewardAwayAcc += winReward;
        // 원정 승리 베터 리스트에 추가
        betAwayWinBettors[_gameId].push(msg.sender);
        // 베팅 성공 이벤트
        emit EvBet(bets[newBetId]);

        return true;
    }

}