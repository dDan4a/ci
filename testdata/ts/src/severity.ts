export type Severity = "critical" | "warning" | "info"

const rank: Record<Severity, number> = { critical: 0, warning: 1, info: 2 }

export function compareSeverity(a: Severity, b: Severity): number {
  return rank[a] - rank[b]
}
