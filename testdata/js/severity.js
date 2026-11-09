const RANK = { critical: 0, warning: 1, info: 2 }

export function compareSeverity(a, b) {
  return (RANK[a] ?? RANK.info) - (RANK[b] ?? RANK.info)
}
