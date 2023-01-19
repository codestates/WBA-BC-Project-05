// SPDX-License-Identifier: MIT

pragma solidity ^0.8.17;

import "./TotoroToken.sol";

contract TotoroGame is TotoroToken {
    uint constant createFee = 1e9;
    uint constant minMaxRewardAmount = 1e18;

    struct Game {
        uint gameId;
        address creator;
        string title;
        string description;
        string home;
        string away;
        uint homeOdd;
        uint awayOdd;
        uint voteHomeCount;
        uint voteAwayCount;
        uint voteVoidCount;
        uint maxRewardAmount;
        uint maxRewardHomeAcc;
        uint maxRewardAwayAcc;
        uint homeAccReward;
        uint awayAccReward;
        uint32 createDate;
        uint32 betEndDate;
        uint32 voteEndDate;
    }

    event EvCreateGame(Game game);

    Game[] games;
    mapping (uint => address) gameOwner;
    mapping (address => uint[]) ownerGames;

    function createGame(string memory _title, string memory _description, string memory _home, string memory _away,
        uint _homeOdd, uint _awayOdd, uint _maxRewardAmount, uint32 _betEndDate, uint32 _voteEndDate) external returns (bool) {

        uint32 currentTime = uint32(block.timestamp);

        // 최소한으로 설정해야 할 최대 지급 가능 잔고 체크
        require(_maxRewardAmount >= minMaxRewardAmount, "Need minimun(1 ** deciamls token) maxRewardAmount");

        // 잔고 체크
        require(balanceOf[msg.sender] >= _maxRewardAmount + createFee, "Not enough balance");

        // 베팅 마감 날짜 체크
        require(currentTime < _betEndDate, "Invalid BetEndDate");

        // 투표 마감 날짜 체크
        require(_betEndDate < _voteEndDate, "Invalid VoteEndDate");

        // 게임 생성자에게 수수료 부과?
        _transferToOwner(msg.sender, createFee);

        // 게임 생성
        uint newGameId = games.length;
        games.push(Game(newGameId, msg.sender, _title, _description, _home, _away, _homeOdd, _awayOdd, 0, 0, 0,
            _maxRewardAmount, 0, 0, 0, 0, uint32(block.timestamp), _betEndDate, _voteEndDate));
        // 게임 아이디 => 게임 생성자 주소 매핑
        gameOwner[newGameId] = msg.sender;
        // 게임 생성자 주소 => 게임 아이디 매핑
        ownerGames[msg.sender].push(newGameId);

        // 게임 마감 시, 지급해야 할 최대 상금 전달
        _transferToOwner(msg.sender, _maxRewardAmount);

        // 게임 생성 성공 이벤트
        emit EvCreateGame(games[newGameId]);

        return true;
    }
}