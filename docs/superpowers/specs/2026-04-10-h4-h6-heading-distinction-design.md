# Design: H4–H6 Heading Distinction

**Date:** 2026-04-10
**Branch:** h4-h6-heading-distinction
**Status:** Approved

## Problem

H4, H5, and H6 headings currently all render using the H3 style — same color, same weight, no visual distinction. Documents that use deep heading hierarchies (API references, RFCs, technical specs) lose all depth information.

## Goal

Give H4, H5, and H6 each a visually distinct style that clearly communicates their position in the heading hierarchy, while remaining cohesive with the existing H1–H3 palette.

## Approach

Add `H4`, `H5`, `H6` fields to the `Theme` struct. Populate them in `DarkTheme()` and `LightTheme()` using progressively dimmer/lower-contrast variants of the H3 color. Update the renderer to dispatch to these styles instead of reusing H3.

No prefix markers (####, #####, ######) will be added — the feature is purely visual color differentiation.

## Theme Changes

### Dark Theme

| Level | Hex color | Notes |
|-------|-----------|-------|
| H3 (existing) | `#C9A86A` | Reference point (golden) |
| H4 (new) | `#A88A55` | ~75% brightness of H3 |
| H5 (new) | `#876C42` | ~55% brightness of H3 |
| H6 (new) | `#665030` | ~40% brightness of H3 |

### Light Theme

| Level | Hex color | Notes |
|-------|-----------|-------|
| H3 (existing) | `#8D6B3F` | Reference point (golden) |
| H4 (new) | `#A07850` | Slightly lighter/muted |
| H5 (new) | `#B38A63` | Further lightened |
| H6 (new) | `#C69C78` | Clearly lighter/lower-contrast |

All H4–H6 styles use bold weight, consistent with H1–H3.

## Affected Files

| File | Change |
|------|--------|
| `internal/theme/theme.go` | Add `H4`, `H5`, `H6` fields to `Theme` struct; populate in `DarkTheme()` and `LightTheme()` |
| `internal/renderer/headings.go` | Change `default: style = r.theme.H3` to dispatch `H4`, `H5`, `H6` separately |
| `internal/renderer/headings_test.go` | Add tests for levels 4, 5, 6 verifying distinct text rendering; update `TestRenderH4UsesH3Style` |
| `testdata/full.md` | Add H4, H5, H6 sample headings |
| `internal/renderer/integration_test.go` | Add integration check for H4–H6 content in full doc render |

## Out of Scope

- Prefix marker symbols (`####`, `#####`, `######`) — explicitly excluded
- Any changes to H1–H3 rendering
- User-configurable theme colors (separate future feature)
