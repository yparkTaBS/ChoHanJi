"use client";

import { Button } from "@/components/ui/button";

export type Direction = "up" | "down" | "left" | "right";

type DirectionalControlsProps = {
  centerLabel?: string;
  gridSize: { rows: number; cols: number };
  onSelect: (dir: Direction) => void;
  position: { r: number; c: number };
  showDirection: (dir: Direction) => boolean;
};

export default function DirectionalControls({
  centerLabel = "You",
  gridSize,
  onSelect,
  position,
  showDirection,
}: DirectionalControlsProps) {
  return (
    <div
      className="pointer-events-none absolute inset-0 grid place-items-center"
      style={{
        gridTemplateColumns: `repeat(${gridSize.cols}, minmax(0, 1fr))`,
        gridTemplateRows: `repeat(${gridSize.rows}, minmax(0, 1fr))`,
      }}
    >
      {/* Up */}
      {showDirection("up") ? (
        <div
          className="pointer-events-auto flex items-center justify-center p-1"
          style={{
            gridColumnStart: position.c + 1,
            gridColumnEnd: position.c + 2,
            gridRowStart: position.r,
            gridRowEnd: position.r + 1,
          }}
        >
          <Button
            size="sm"
            variant="ghost"
            className="w-full bg-background/50 hover:bg-background/70"
            onClick={() => onSelect("up")}
            aria-label="Move up"
          >
            ↑
          </Button>
        </div>
      ) : null}

      {/* Left */}
      {showDirection("left") ? (
        <div
          className="pointer-events-auto flex items-center justify-center p-1"
          style={{
            gridColumnStart: position.c,
            gridColumnEnd: position.c + 1,
            gridRowStart: position.r + 1,
            gridRowEnd: position.r + 2,
          }}
        >
          <Button
            size="sm"
            variant="ghost"
            className="w-full bg-background/50 hover:bg-background/70"
            onClick={() => onSelect("left")}
            aria-label="Move left"
          >
            ←
          </Button>
        </div>
      ) : null}

      {/* Center marker */}
      <div
        className="pointer-events-none flex items-center justify-center p-1"
        style={{
          gridColumnStart: position.c + 1,
          gridColumnEnd: position.c + 2,
          gridRowStart: position.r + 1,
          gridRowEnd: position.r + 2,
        }}
      >
        <div className="flex h-full items-center justify-center rounded-md border border-dashed text-xs font-medium text-muted-foreground">
          {centerLabel}
        </div>
      </div>

      {/* Right */}
      {showDirection("right") ? (
        <div
          className="pointer-events-auto flex items-center justify-center p-1"
          style={{
            gridColumnStart: position.c + 2,
            gridColumnEnd: position.c + 3,
            gridRowStart: position.r + 1,
            gridRowEnd: position.r + 2,
          }}
        >
          <Button
            size="sm"
            variant="ghost"
            className="w-full bg-background/50 hover:bg-background/70"
            onClick={() => onSelect("right")}
            aria-label="Move right"
          >
            →
          </Button>
        </div>
      ) : null}

      {/* Down */}
      {showDirection("down") ? (
        <div
          className="pointer-events-auto flex items-center justify-center p-1"
          style={{
            gridColumnStart: position.c + 1,
            gridColumnEnd: position.c + 2,
            gridRowStart: position.r + 2,
            gridRowEnd: position.r + 3,
          }}
        >
          <Button
            size="sm"
            variant="ghost"
            className="w-full bg-background/50 hover:bg-background/70"
            onClick={() => onSelect("down")}
            aria-label="Move down"
          >
            ↓
          </Button>
        </div>
      ) : null}
    </div>
  );
}
