import { useCallback, useEffect, useState } from "react";

export type Position = { r: number; c: number };

export type CommitPositionsFn = (
  nextPlayer: Position | null,
  nextEnemy: Position | null,
  startPlayer: Position,
  startEnemy: Position
) => void;

export type OddEvenResult = "Player won" | "Enemy won" | null;

type UseOddEvenGameParams = {
  blocked?: boolean;
  commitPositions: CommitPositionsFn;
  enemyPos: Position;
  isAdjacent: (enemy: Position, player: Position) => boolean;
  isSameTile: (a: Position, b: Position) => boolean;
  playerPos: Position;
};

const DEFAULT_PROMPT = "Guess whether the number is odd or even.";
const COUNTER_PROMPT = "Counterattack! Defend by guessing odd or even.";
const QUEUED_COUNTER_PROMPT = "Attack failed! Close this popup to face a counterattack.";

export function useOddEvenGame({
  blocked = false,
  commitPositions,
  enemyPos,
  isAdjacent,
  isSameTile,
  playerPos,
}: UseOddEvenGameParams) {
  const [showGame, setShowGame] = useState(false);
  const [result, setResult] = useState<OddEvenResult>(null);
  const [rolled, setRolled] = useState<number | null>(null);

  const [playerCommencedAttack, setPlayerCommencedAttack] = useState(false);
  const [counterMinigameQueued, setCounterMinigameQueued] = useState(false);
  const [counterStartArmed, setCounterStartArmed] = useState(false);

  const [gamePrompt, setGamePrompt] = useState<string>(DEFAULT_PROMPT);

  const startGame = useCallback(
    (playerStartedAttack: boolean) => {
      if (!isAdjacent(enemyPos, playerPos) && !isSameTile(enemyPos, playerPos)) return;

      setResult(null);
      setRolled(null);
      setPlayerCommencedAttack(playerStartedAttack);

      setCounterMinigameQueued(false);
      setCounterStartArmed(false);

      setGamePrompt(DEFAULT_PROMPT);
      setShowGame(true);
    },
    [enemyPos, isAdjacent, isSameTile, playerPos]
  );

  const startCounterMinigame = useCallback(() => {
    if (!isAdjacent(enemyPos, playerPos)) {
      setCounterMinigameQueued(false);
      setCounterStartArmed(false);
      return;
    }

    setPlayerCommencedAttack(false);

    setResult(null);
    setRolled(null);

    setGamePrompt(COUNTER_PROMPT);
    setShowGame(true);
  }, [enemyPos, isAdjacent, playerPos]);

  const closeGamePopup = useCallback(() => {
    if (counterMinigameQueued) {
      setShowGame(false);
      setCounterStartArmed(true);
      return;
    }

    setShowGame(false);
    setCounterMinigameQueued(false);
    setCounterStartArmed(false);
  }, [counterMinigameQueued]);

  useEffect(() => {
    if (!counterStartArmed || showGame) return;

    setCounterStartArmed(false);
    setCounterMinigameQueued(false);
    startCounterMinigame();
  }, [counterStartArmed, showGame, startCounterMinigame]);

  useEffect(() => {
    if (showGame || blocked) return;
    if (isSameTile(playerPos, enemyPos)) {
      startGame(false);
    }
  }, [blocked, enemyPos, isSameTile, playerPos, showGame, startGame]);

  const play = useCallback(
    (selected: "odd" | "even") => {
      if (!showGame || result) return;

      const startPlayer = { ...playerPos };
      const startEnemy = { ...enemyPos };

      const n = Math.floor(Math.random() * 10) + 1; // 1..10
      setRolled(n);

      const isEven = n % 2 === 0;
      const playerCorrect =
        (selected === "even" && isEven) || (selected === "odd" && !isEven);

      const nextResult: Exclude<OddEvenResult, null> = playerCorrect
        ? "Player won"
        : "Enemy won";
      setResult(nextResult);

      if (nextResult === "Enemy won") {
        if (!playerCommencedAttack) {
          commitPositions({ r: 0, c: 0 }, null, startPlayer, startEnemy);
        } else {
          setCounterMinigameQueued(true);
          setGamePrompt(QUEUED_COUNTER_PROMPT);
        }
      } else {
        commitPositions(null, { r: 4, c: 4 }, startPlayer, startEnemy);
      }
    },
    [commitPositions, enemyPos, playerCommencedAttack, playerPos, result, showGame]
  );

  return {
    closeGamePopup,
    gamePrompt,
    play,
    result,
    rolled,
    showGame,
    startGame,
  };
}
