"use client";

import { use, useEffect, useRef, useState, useMemo } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

type row = [string, string, string, string, string];
type grid = [row, row, row, row, row];

type Mode = "move" | "attack";
type Direction = "up" | "down" | "left" | "right" | null;

function cloneGrid(g: grid): grid {
  // Deduction: need new arrays so React sees changes + no shared references
  return g.map((r) => [...r] as row) as grid;
}

function inBounds(r: number, c: number) {
  return r >= 0 && r <= 4 && c >= 0 && c <= 4;
}

const DELTAS: Record<Exclude<Direction, null>, [number, number]> = {
  up: [-1, 0],
  down: [1, 0],
  left: [0, -1],
  right: [0, 1],
};

const ENEMY_DIRS: Array<[number, number]> = [
  [-1, 0],
  [1, 0],
  [0, -1],
  [0, 1],
];

function shuffle<T>(arr: T[]) {
  // Deduction: Fisher–Yates shuffle for unbiased random selection
  const a = [...arr];
  for (let i = a.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [a[i], a[j]] = [a[j], a[i]];
  }
  return a;
}

function isAdjacent4(a: { r: number; c: number }, b: { r: number; c: number }) {
  const dr = Math.abs(a.r - b.r);
  const dc = Math.abs(a.c - b.c);
  return (dr === 1 && dc === 0) || (dr === 0 && dc === 1);
}

function computeEnemyMove(enemy: { r: number; c: number }) {
  // Deduction: enemy movement is 1 random step in-bounds
  const dirs = shuffle(ENEMY_DIRS);
  for (const [dr, dc] of dirs) {
    const nr = enemy.r + dr;
    const nc = enemy.c + dc;
    if (!inBounds(nr, nc)) continue;
    return { r: nr, c: nc };
  }
  return enemy;
}

export default function CharacterPage({
  params,
}: {
  params: Promise<{ roomId: string; characterId: string }>;
}) {
  const esRef = useRef<EventSource | null>(null);

  const initialField: grid = [
    ["", "", "", "", ""],
    ["", "", "", "", ""],
    ["", "", "", "", ""],
    ["", "", "", "", ""],
    ["", "", "", "", ""],
  ];

  const [map] = useState<grid>(initialField);

  const [mode, setMode] = useState<Mode>("move");

  const [pos, setPos] = useState<{ r: number; c: number }>({ r: 0, c: 0 });
  const [enemyPos, setEnemyPos] = useState<{ r: number; c: number }>({
    r: 4,
    c: 4,
  });

  const { roomId, characterId } = use(params);

  const [pendingDir, setPendingDir] = useState<Direction>(null);

  // Odd/Even game popup
  const [showGame, setShowGame] = useState(false);
  const [result, setResult] = useState<"Player won" | "Enemy won" | null>(null);
  const [rolled, setRolled] = useState<number | null>(null);

  // Preserve whether the PLAYER commenced an attack (needed for outcome rules)
  const [playerCommencedAttack, setPlayerCommencedAttack] = useState(false);

  // If player's initiated attack fails, queue a second (counterattack) minigame
  const [counterMinigameQueued, setCounterMinigameQueued] = useState(false);

  // Only start counter minigame AFTER player closes the popup
  const [counterStartArmed, setCounterStartArmed] = useState(false);

  const [gamePrompt, setGamePrompt] = useState<string>(
    "Guess whether the number is odd or even."
  );

  // "Move not allowed" popup (when player tries to move into 4,4)
  const [showMoveNotAllowed, setShowMoveNotAllowed] = useState(false);

  // NEW: directional attack availability
  function canAttackDir(dir: Exclude<Direction, null>) {
    // Deduction: "enemy is occupying that block" means the adjacent tile in that direction equals enemyPos
    const [dr, dc] = DELTAS[dir];
    const target = { r: pos.r + dr, c: pos.c + dc };
    return target.r === enemyPos.r && target.c === enemyPos.c;
  }

  // NEW: if we are in attack mode and the selected direction no longer points at the enemy, clear it
  useEffect(() => {
    if (mode !== "attack") return;
    if (!pendingDir) return;
    if (!canAttackDir(pendingDir)) setPendingDir(null);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [mode, pos, enemyPos]);

  function startOddEvenGame(playerStartedAttack: boolean) {
    // Deduction: minigame should only be possible if enemy is adjacent RIGHT NOW
    if (!isAdjacent4(enemyPos, pos)) return;

    setPendingDir(null);
    setResult(null);
    setRolled(null);
    setPlayerCommencedAttack(playerStartedAttack);

    setCounterMinigameQueued(false);
    setCounterStartArmed(false);

    setGamePrompt("Guess whether the number is odd or even.");
    setShowGame(true);
  }

  function startCounterMinigame() {
    // Deduction: counter minigame should only happen if still adjacent
    if (!isAdjacent4(enemyPos, pos)) {
      setCounterMinigameQueued(false);
      setCounterStartArmed(false);
      return;
    }

    // Deduction: counterattack is a DEFENSE scenario
    setPlayerCommencedAttack(false);

    setResult(null);
    setRolled(null);

    setGamePrompt("Counterattack! Defend by guessing odd or even.");
    setShowGame(true);
  }

  function closeGamePopup() {
    // Deduction: player must close popup first; counter begins only after close if queued
    if (counterMinigameQueued) {
      setShowGame(false);
      setCounterStartArmed(true);
      return;
    }

    setShowGame(false);
    setCounterMinigameQueued(false);
    setCounterStartArmed(false);
  }

  useEffect(() => {
    if (!counterStartArmed) return;
    if (showGame) return; // must be closed first

    setCounterStartArmed(false);
    setCounterMinigameQueued(false);
    startCounterMinigame();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [counterStartArmed, showGame]);

  function toggleMode() {
    setMode((prev) => (prev === "move" ? "attack" : "move"));
  }

  const displayMap = useMemo(() => {
    const next = cloneGrid(map);
    next[enemyPos.r][enemyPos.c] = "X";
    next[pos.r][pos.c] = next[pos.r][pos.c] ? `@ ${next[pos.r][pos.c]}` : "@";
    return next;
  }, [map, pos, enemyPos]);

  function submit() {
    if (!pendingDir || showGame || showMoveNotAllowed) return;

    const dir = pendingDir;
    const startPlayer = { ...pos };
    const startEnemy = { ...enemyPos };

    // 1) Player commenced an attack => ONLY start minigame if:
    // - adjacent, AND
    // - the selected direction actually targets the enemy tile (per new directional rule)
    if (mode === "attack") {
      if (isAdjacent4(startEnemy, startPlayer) && canAttackDir(dir)) {
        startOddEvenGame(true);
      }
      setPendingDir(null);
      return;
    }

    // 2) Enemy commenced an attack => ONLY if adjacent at start-of-turn
    if (isAdjacent4(startEnemy, startPlayer)) {
      startOddEvenGame(false);
      setPendingDir(null);
      return;
    }

    // 3) Compute player movement
    let playerNext = { ...startPlayer };
    {
      const [dr, dc] = DELTAS[dir];
      const nr = startPlayer.r + dr;
      const nc = startPlayer.c + dc;

      if (inBounds(nr, nc)) {
        if (nr === 4 && nc === 4) {
          setShowMoveNotAllowed(true);
          setPendingDir(null);
          return;
        }
        playerNext = { r: nr, c: nc };
      }
    }

    // 4) Compute enemy movement
    const enemyNext = computeEnemyMove(startEnemy);

    // 5) Prevent minigame when enemy is not adjacent
    setPos(playerNext);
    setEnemyPos(enemyNext);
    setPendingDir(null);
  }

  function skipTurn() {
    if (showGame || showMoveNotAllowed) return;

    const startPlayer = { ...pos };
    const startEnemy = { ...enemyPos };

    setPendingDir(null);

    // Enemy attacks only if adjacent at start-of-turn
    if (isAdjacent4(startEnemy, startPlayer)) {
      startOddEvenGame(false);
      return;
    }

    const enemyNext = computeEnemyMove(startEnemy);

    // Deduction: even if enemy overlaps, do NOT start minigame unless adjacent
    setEnemyPos(enemyNext);
  }

  function playOddEven(selected: "odd" | "even") {
    if (!showGame || result) return;

    const n = Math.floor(Math.random() * 10) + 1; // 1..10
    setRolled(n);

    const isEven = n % 2 === 0;
    const playerCorrect =
      (selected === "even" && isEven) || (selected === "odd" && !isEven);

    const nextResult: "Player won" | "Enemy won" = playerCorrect
      ? "Player won"
      : "Enemy won";
    setResult(nextResult);

    /**
     * Outcomes:
     * - If player wins => enemy back to (4,4)
     * - If player loses:
     *   - if defending => player back to (0,0)
     *   - if attacking => player does NOT go back; BUT queue counterattack minigame
     */
    if (nextResult === "Enemy won") {
      if (!playerCommencedAttack) {
        setPos({ r: 0, c: 0 });
      } else {
        setCounterMinigameQueued(true);
        setGamePrompt("Attack failed! Close this popup to face a counterattack.");
      }
    } else {
      setEnemyPos({ r: 4, c: 4 });
    }
  }

  useEffect(() => {
    const es = new EventSource(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}api/player/event?roomId=${roomId}&playerId=${characterId}`
    );
    esRef.current = es;

    es.onopen = () => console.log("SSE connected");
    es.onmessage = (event) => console.log("SSE message:", event.data);
    es.onerror = (err) => console.error("SSE error", err);

    return () => {
      es.close();
      esRef.current = null;
    };
  }, [roomId, characterId]);

  // NEW: directional buttons should be disabled in attack mode unless enemy occupies that block
  const disableDir = (dir: Exclude<Direction, null>) =>
    showGame ||
    showMoveNotAllowed ||
    (mode === "attack" && !canAttackDir(dir));

  return (
    <Card className="w-fit">
      {/* Move not allowed Popup */}
      {showMoveNotAllowed && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="w-full max-w-sm rounded-xl border bg-background p-4 space-y-3">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold">Move not allowed</h2>
              <Button
                variant="ghost"
                onClick={() => setShowMoveNotAllowed(false)}
              >
                ✕
              </Button>
            </div>
            <p className="text-sm text-muted-foreground">
              You cannot move into the enemy spawn point
            </p>
            <div className="flex justify-end">
              <Button onClick={() => setShowMoveNotAllowed(false)}>OK</Button>
            </div>
          </div>
        </div>
      )}

      {/* Odd/Even Popup */}
      {showGame && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="w-full max-w-sm rounded-xl border bg-background p-4 space-y-3">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold">Odd / Even</h2>
              <Button variant="ghost" onClick={closeGamePopup}>
                ✕
              </Button>
            </div>

            <p className="text-sm text-muted-foreground">{gamePrompt}</p>

            <div className="flex gap-2">
              <Button
                className="flex-1"
                disabled={!!result}
                onClick={() => playOddEven("odd")}
              >
                Odd
              </Button>
              <Button
                className="flex-1"
                disabled={!!result}
                onClick={() => playOddEven("even")}
              >
                Even
              </Button>
            </div>

            {rolled !== null && (
              <p className="text-sm">
                Rolled: <strong>{rolled}</strong>
              </p>
            )}

            {result && <p className="text-base font-semibold">{result}</p>}

            <div className="flex justify-end gap-2 pt-2">
              <Button onClick={closeGamePopup}>Close</Button>
            </div>
          </div>
        </div>
      )}

      <CardHeader>
        <CardTitle>Map</CardTitle>
      </CardHeader>

      <CardContent className="space-y-6">
        {/* Tutorial (kept) */}
        <Card className="border-muted">
          <CardHeader>
            <CardTitle className="text-base">Tutorial</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2 text-sm text-muted-foreground">
            <p>
              <strong>Move Mode:</strong> Select a direction and press{" "}
              <strong>Move</strong>.
            </p>
            <p>
              <strong>Attack Mode:</strong> Select a direction and press{" "}
              <strong>Attack</strong>.
            </p>
            <p>
              <strong>Order of Actions:</strong> Attack move will always happen
              first, followed by the normal movement, bonus movement, bonus
              attack, and collision resolution.
            </p>
            <p>
              <strong>Collision:</strong> If you bump into an enemy (or an enemy bumps into you), a mini game starts. Teammates can overlap freely.
            </p>
            <p>
              <strong>Rule:</strong> Moving into enemy spawn point (4, 4) is not
              allowed.
            </p>
            <p>
              <strong>Classes:</strong> Each class has unique bonuses. Thief has bonus moves. Fighters can attack after moving. Rangers can snipe enemies from afar.
            </p>
            <p>
              <strong>Death:</strong> If you lose while defending, you move back
              to your spawn point (0,0)
            </p>
            <p>
              <strong>Skip:</strong> You can skip your turn; the enemy may still
              attack if adjacent or move otherwise.
            </p>
          </CardContent>
        </Card>

        {/* Map */}
        <div className="inline-grid grid-cols-5 rounded-xl border border-border overflow-hidden">
          {displayMap.flatMap((r, ri) =>
            r.map((cell, ci) => (
              <div
                key={`${ri}-${ci}`}
                className="flex h-20 w-20 items-center justify-center bg-background text-lg font-semibold border border-border -ml-px -mt-px"
              >
                {cell || ""}
              </div>
            ))
          )}
        </div>

        {/* Controls */}
        <div className="grid grid-cols-3 gap-2 place-items-center">
          <div />
          <Button disabled={disableDir("up")} onClick={() => setPendingDir("up")}>
            Up
          </Button>
          <div />

          <Button
            disabled={disableDir("left")}
            onClick={() => setPendingDir("left")}
          >
            Left
          </Button>

          <Button variant="secondary" onClick={toggleMode}>
            Switch to {mode === "move" ? "Attack" : "Move"} Mode
          </Button>

          <Button
            disabled={disableDir("right")}
            onClick={() => setPendingDir("right")}
          >
            Right
          </Button>

          <div />
          <Button
            disabled={disableDir("down")}
            onClick={() => setPendingDir("down")}
          >
            Down
          </Button>
          <div />

          <div />
          <div className="flex gap-2 pt-2">
            <Button
              disabled={!pendingDir || showGame || showMoveNotAllowed}
              onClick={submit}
            >
              {mode === "move" ? "Move" : "Attack"}
            </Button>

            <Button
              variant="outline"
              disabled={showGame || showMoveNotAllowed}
              onClick={skipTurn}
            >
              Skip Turn
            </Button>

            <Button
              variant="outline"
              disabled={!pendingDir || showGame || showMoveNotAllowed}
              onClick={() => setPendingDir(null)}
            >
              Clear
            </Button>
          </div>
          <div />

          <div />
          {pendingDir ? (
            <p className="text-sm text-muted-foreground">
              Selected: <strong>{pendingDir.toUpperCase()}</strong>
            </p>
          ) : (
            <div />
          )}
          <div />
        </div>
      </CardContent>
    </Card>
  );
}

