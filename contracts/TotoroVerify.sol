// SPDX-License-Identifier: MIT

pragma solidity ^0.8.17;

import "./TotoroBet.sol";

contract TotoroVerify is TotoroBet {
    // 검증 투표 상수
    uint8 constant VERIFY_VOTE_HOME = 0; // 홈팀 승리
    uint8 constant VERIFY_VOTE_AWAY = 1; // 원정팀 승리
    uint8 constant VERIFY_VOTE_VOID = 2; // 무효

    struct Verify {
        uint gameId;
        address verifier;
        uint8 vote;
    }

    event EvVerify(uint verifyId, Verify verify);
    event EvCalculate(uint gameId);

    Verify[] verifys;
    // gameId => verifyId 매핑
    mapping (uint => uint[]) gameIdVerifyIds;
    // verifyId => 계정 주소 매핑
    mapping (uint => address) verifyOwner;
    // 계정 주소 => verifyId[] 매핑
    mapping (address => uint[]) ownerVerifys;

    modifier verifyValidCheck(uint _gameId, uint8 _vote) {
        uint32 currentTime = uint32(block.timestamp);
        // 게임 아이디 유효성 체크
        require(_gameId != 0);
        // 베팅 마감 날짜 체크
        require(currentTime > games[_gameId].betEndDate);
        // 자신이 베팅한 게임은 검증할 수 없음
        uint[] memory gameBets = gameIdbetIds[_gameId];
        for (uint i=0; i<gameBets.length; i++) {
            uint betId = gameBets[i];
            require(betOwner[betId] != msg.sender);
        }
        // 검증 마감 날짜가 지난 경우 : 정산 처리
        if (currentTime > games[_gameId].verifyEndDate) {
            calculate(_gameId);
            return;
        }
        _;
    }

    function voteHome(uint _gameId) external verifyValidCheck(_gameId, VERIFY_VOTE_HOME) {
        games[_gameId].verifyHomeCount++;
        verify(_gameId, VERIFY_VOTE_HOME);
    }

    function voteAway(uint _gameId) external verifyValidCheck(_gameId, VERIFY_VOTE_AWAY) {
        games[_gameId].verifyAwayCount++;
        verify(_gameId, VERIFY_VOTE_AWAY);
    }

    function voteVoid(uint _gameId) external verifyValidCheck(_gameId, VERIFY_VOTE_VOID) {
        games[_gameId].verifyVoidCount++;
        verify(_gameId, VERIFY_VOTE_VOID);
    }

    function verify(uint _gameId, uint8 _vote) internal returns (bool) {
        // 검증 생성
        verifys.push(Verify(_gameId, msg.sender, _vote));
        // 검증 아이디 생성
        uint newVerifyId = verifys.length - 1;
        // 게임 아이디 => 검증 아이디 매핑
        gameIdVerifyIds[_gameId].push(newVerifyId); 
        // 검증 아이디 => 검증자 매핑
        verifyOwner[newVerifyId] = msg.sender;
        // 검증자 => 검증 아이디 추가
        ownerVerifys[msg.sender].push(newVerifyId);
        // 검증 성공 이벤트
        emit EvVerify(newVerifyId, verifys[newVerifyId]);
        // 검증자에게 보상 지급
        balanceOf[msg.sender] += 10000;

        return true;
    }

    function findMaxIdx(uint[] memory array) internal pure returns(uint) {
        uint maxValue = 0;
        uint maxValueIdx = 0;
        for (uint i=0; i<array.length; i++) {
            if (maxValue < array[i]) {
                maxValue = array[i];
                maxValueIdx = i;
            }
        }
        return maxValueIdx;
    }

    // 홈팀 승리 처리 함수
    function winHome(uint _gameId) internal {
        uint odds = games[_gameId].homeOdd;
        uint accReward = games[_gameId].homeAccReward;

        // 홈 승리에 베팅한 베터에게 보상 지급
        for (uint8 i=0; i<gameIdbetIds[_gameId].length; i++) {
            Bet memory bet = bets[gameIdbetIds[_gameId][i]];
            if (bet.betTarget == VERIFY_VOTE_HOME) {
                uint reward = bet.amount * odds;
                // 베팅 적중에 따른 보상
                balanceOf[bet.bettor] += reward;
                accReward -= reward;
            }
        }
        // 남은 상금은 게임 생성자에게 전달
        balanceOf[games[_gameId].creator] += accReward;
        balanceOf[games[_gameId].creator] += games[_gameId].awayAccReward;
    }

    // 원정팀 승리 처리 함수
    function winAway(uint _gameId) internal {
        uint odds = games[_gameId].awayOdd;
        uint accReward = games[_gameId].awayAccReward;

        // 득표에 따른 보상 처리
        for (uint8 i=0; i<gameIdbetIds[_gameId].length; i++) {
            Bet memory bet = bets[gameIdbetIds[_gameId][i]];
            if (bet.betTarget == VERIFY_VOTE_AWAY) {
                // 베팅 적중에 따른 보상
                uint reward = bet.amount * odds;
                balanceOf[bet.bettor] += reward;
                accReward -= reward;
            }
        }
        // 남은 상금은 게임 생성자에게 전달
        balanceOf[games[_gameId].creator] += accReward;
        balanceOf[games[_gameId].creator] += games[_gameId].awayAccReward;
    }

    // 무효 처리 함수
    function winVoid(uint _gameId) internal {
        // 게임 생성자의 동결된 자금 반환
        uint reward = rewardLock[_gameId];
        balanceOf[games[_gameId].creator] = reward;

        // 베팅 참여자에게 베팅 금액 반환
        uint[] memory gameBets = gameIdbetIds[_gameId];
        for (uint i=0; i<gameBets.length; i++) {
            Bet memory bet = bets[gameBets[i]];
            balanceOf[bet.bettor] += bet.amount;
        }
    }

    // 정산 처리 함수
    function calculate(uint _gameId) internal {
        Game memory game = games[_gameId];
        // 가장 많은 검증 득표를 받은 항목 찾기
        uint[] memory votes = new uint[](3);
        votes[0] = game.verifyHomeCount;
        votes[1] = game.verifyAwayCount;
        votes[2] = game.verifyVoidCount;
        uint winVoteIdx = findMaxIdx(votes);

        // 홈팀 승리 처리
        if (winVoteIdx == VERIFY_VOTE_HOME) {
            winHome(_gameId);
        }
        // 원정팀 승리 처리
        else if(winVoteIdx == VERIFY_VOTE_AWAY) {
            winAway(_gameId);
        }
        // 게임 무효 처리
        else {
            winVoid(_gameId);
        }

        // 정산 성공 이벤트
        emit EvCalculate(_gameId);
    }
}