// SPDX-License-Identifier: MIT

pragma solidity ^0.8.17;

import "./TotoroBet.sol";

contract TotoroVote is TotoroBet {
    // 투표 상수
    uint8 constant VOTE_TARGET_HOME = 0; // 홈팀 승리
    uint8 constant VOTE_TARGET_AWAY = 1; // 원정팀 승리
    uint8 constant VOTE_TARGET_VOID = 2; // 무효

    struct Vote {
        uint voteId;
        uint gameId;
        address voter;
        uint8 target;
    }

    event EvVote(Vote vote);
    event EvCalculate(uint gameId);

    Vote[] votes;
    // gameId => voteId 매핑
    mapping (uint => uint[]) gameIdVoteIds;
    // voteId => 계정 주소 매핑
    mapping (uint => address) voteOwner;
    // 계정 주소 => voteId[] 매핑
    mapping (address => uint[]) ownerVotes;

    modifier voteValidCheck(uint _gameId, uint8 _vote) {
        uint32 currentTime = uint32(block.timestamp);
        // 게임 아이디 유효성 체크
        require(games.length - 1 >= _gameId);
        // 베팅 마감 날짜 체크
        require(currentTime > games[_gameId].betEndDate);
        // 자신이 베팅한 게임은 투표할 수 없음
        uint[] memory gameBets = gameIdbetIds[_gameId];
        for (uint i=0; i<gameBets.length; i++) {
            uint betId = gameBets[i];
            require(betOwner[betId] != msg.sender);
        }
        // 투표 마감 날짜가 지난 경우 : 정산 처리
        if (currentTime > games[_gameId].voteEndDate) {
            calculate(_gameId);
            return;
        }
        _;
    }

    function voteHome(uint _gameId) external voteValidCheck(_gameId, VOTE_TARGET_HOME) {
        games[_gameId].voteHomeCount++;
        vote(_gameId, VOTE_TARGET_HOME);
    }

    function voteAway(uint _gameId) external voteValidCheck(_gameId, VOTE_TARGET_AWAY) {
        games[_gameId].voteAwayCount++;
        vote(_gameId, VOTE_TARGET_AWAY);
    }

    function voteVoid(uint _gameId) external voteValidCheck(_gameId, VOTE_TARGET_VOID) {
        games[_gameId].voteVoidCount++;
        vote(_gameId, VOTE_TARGET_VOID);
    }

    function vote(uint _gameId, uint8 _target) internal returns (bool) {
        // 투표 생성
        uint newVoteId = votes.length;
        votes.push(Vote(newVoteId, _gameId, msg.sender, _target));
        // 게임 아이디 => 투표 아이디 매핑
        gameIdVoteIds[_gameId].push(newVoteId); 
        // 투표 아이디 => 투표자 매핑
        voteOwner[newVoteId] = msg.sender;
        // 투표자 => 투표 아이디 추가
        ownerVotes[msg.sender].push(newVoteId);
        // 투표 성공 이벤트
        emit EvVote(votes[newVoteId]);
        // 투표자에게 보상 지급
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
            if (bet.target == VOTE_TARGET_HOME) {
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
            if (bet.target == VOTE_TARGET_AWAY) {
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
        // 가장 많은 득표를 받은 항목 찾기
        uint[] memory voteCounts = new uint[](3);
        voteCounts[0] = game.voteHomeCount;
        voteCounts[1] = game.voteAwayCount;
        voteCounts[2] = game.voteVoidCount;
        uint winVoteIdx = findMaxIdx(voteCounts);

        // 홈팀 승리 처리
        if (winVoteIdx == VOTE_TARGET_HOME) {
            winHome(_gameId);
        }
        // 원정팀 승리 처리
        else if(winVoteIdx == VOTE_TARGET_AWAY) {
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