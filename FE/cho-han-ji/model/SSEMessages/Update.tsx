export type PlayerChangePayload = {
  X: number;
  Y: number;
  PrevX: number;
  PrevY: number;
  Id: string;
  ItemId?: string | null;
};

export type ItemChangePayload = {
  X: number;
  Y: number;
  PrevX: number;
  PrevY: number;
  ItemId: string;
};

export type UpdateMessagePayload = {
  PlayerChanges?: Record<string, PlayerChangePayload>;
  ItemChanges?: Record<string, ItemChangePayload>;
};

type EnvelopeLike = {
  Message?: unknown;
};

const normalizeNumber = (value: unknown): number | null => {
  if (typeof value === "number" && Number.isFinite(value)) return value;
  if (typeof value === "string" && value.trim() !== "") {
    const parsed = Number(value);
    return Number.isFinite(parsed) ? parsed : null;
  }
  return null;
};

const parsePlayerChange = (id: string, change: unknown): PlayerChangePayload | null => {
  if (!change || typeof change !== "object") return null;
  const raw = change as Partial<PlayerChangePayload>;

  const X = normalizeNumber(raw.X);
  const Y = normalizeNumber(raw.Y);
  const PrevX = normalizeNumber(raw.PrevX);
  const PrevY = normalizeNumber(raw.PrevY);
  const Id = typeof raw.Id === "string" && raw.Id !== "" ? raw.Id : id;

  if (X === null || Y === null || PrevX === null || PrevY === null || !Id) return null;

  return {
    X,
    Y,
    PrevX,
    PrevY,
    Id,
    ItemId: raw.ItemId ?? null,
  };
};

const parseItemChange = (id: string, change: unknown): ItemChangePayload | null => {
  if (!change || typeof change !== "object") return null;
  const raw = change as Partial<ItemChangePayload>;

  const X = normalizeNumber(raw.X);
  const Y = normalizeNumber(raw.Y);
  const PrevX = normalizeNumber(raw.PrevX);
  const PrevY = normalizeNumber(raw.PrevY);
  const ItemId = typeof raw.ItemId === "string" && raw.ItemId !== "" ? raw.ItemId : id;

  if (X === null || Y === null || PrevX === null || PrevY === null || !ItemId) return null;

  return {
    X,
    Y,
    PrevX,
    PrevY,
    ItemId,
  };
};

const parsePayload = (message: unknown): UpdateMessagePayload | null => {
  if (typeof message === "string") {
    try {
      const parsed = JSON.parse(message);
      if (parsed && typeof parsed === "object") {
        return parsed as UpdateMessagePayload;
      }
      return null;
    } catch {
      return null;
    }
  }

  if (!message || typeof message !== "object") return null;

  return message as UpdateMessagePayload;
};

const unwrapEnvelope = (value: unknown): UpdateMessagePayload | null => {
  let current = parsePayload(value);

  // Unwrap envelopes like { MessageType, Message } until we reach a payload
  // that actually contains the change collections.
  for (let depth = 0; depth < 3 && current; depth += 1) {
    if (current.PlayerChanges || current.ItemChanges) break;
    if (!("Message" in current)) break;

    const envelope = current as EnvelopeLike;
    const next = parsePayload(envelope.Message);
    if (!next || next === current) break;

    current = next;
  }

  return current;
};

export const parseUpdateMessage = (message: unknown) => {
  const payload = unwrapEnvelope(message);

  if (!payload) {
    return { playerChanges: [] as PlayerChangePayload[], itemChanges: [] as ItemChangePayload[] };
  }

  const playerChanges = Object.entries(payload.PlayerChanges ?? {})
    .map(([id, change]) => parsePlayerChange(id, change))
    .filter((change): change is PlayerChangePayload => Boolean(change));

  const itemChanges = Object.entries(payload.ItemChanges ?? {})
    .map(([id, change]) => parseItemChange(id, change))
    .filter((change): change is ItemChangePayload => Boolean(change));

  return { playerChanges, itemChanges };
};
