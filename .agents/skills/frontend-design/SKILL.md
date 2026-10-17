---
name: frontend-design
description: Art direction, UX architecture, distinct visual design, responsive implementation, and browser-based visual QA for goserver/GOSH websites. Use for new sites, landing pages, UI features, major styling changes, or requested redesigns; not for API-only tasks.
---

# Frontend Design Director

Create a coherent, usable, distinctive website—not a mechanically assembled generic AI template. Act as product designer, art director, interaction designer, and frontend engineer. Follow repository `AGENTS.md` for goserver/GOSH architecture and tooling; this skill governs visual and UX decisions. Explicit user instructions and supplied references override aesthetic defaults.

## Scope and effort

- **New website / redesign:** use the complete workflow, including browser screenshot review and iteration.
- **New section / significant component:** establish local direction and check it in context; reuse existing tokens where appropriate.
- **Small visual fix:** diagnose the affected layout and verify nearby breakpoints; avoid introducing a new design language.
- **Backend/API-only:** do not apply this skill.

## 1. Understand the product before styling

- Inspect the brief, existing pages, content, assets, screenshots, brand constraints, and actual user tasks. Identify audience, page purpose, essential information, primary action, and friction points.
- Respect existing design language when not explicitly redesigning. Never manufacture testimonials, metrics, clients, reviews, UI screenshots, or business claims.
- Write usable real copy in the requested language *before* sizing compositions. Handle long headings, Cyrillic, navigation labels, translations, and realistic content density.
- Define the first viewport around a meaningful subject, useful information, and an obvious next step—not an obligatory centered hero.

## 2. Art direction before code

For new visual work, consider **three meaningfully different** directions internally (e.g. editorial, technical/minimal, expressive/product-led). Different directions require different composition, typography, information density, and imagery—not merely different accent colors. Select the strongest based on audience, brand, content, and implementation constraints.

Before writing UI, briefly establish concrete choices:

- **Composition:** asymmetry/symmetry, hierarchy, section rhythm, alignment, density, whitespace.
- **Typography:** font strategy, scale, weights, measures, title wrapping, navigation hierarchy.
- **Color:** actual neutral/accent values, surfaces, states, accessible contrast.
- **Media:** real relevant images/graphics, crop strategy, aspect ratios, licensing.
- **Interaction:** primary action, menu/forms/states, restrained motion, touch behavior.
- **Mobile:** re-composition decisions, reading order, action placement, media cropping.

Do not turn an implementation request into a consultation; make reasonable decisions and implement them.

## 3. Reject generic UI defaults

Do not automatically produce a centered hero, three rounded feature cards, icon-title-description repetitions, benefit strip, invented testimonial carousel, and generic final CTA. Every section needs a content or task rationale. Avoid ubiquitous card containers, ornamental borders/shadows, purple-blue gradients, glassmorphism, glowing blobs, stickers, oversized radii, and default SaaS-dashboard composition unless the direction specifically calls for them.

- Create identity with **typography, proportion, grouping, pacing, and subject-specific content** before ornament.
- Use cards only for distinct grouped entities, not as containers for every paragraph. Vary treatment by information hierarchy while maintaining system consistency.
- Avoid looking interchangeable with an unrelated brand after swapping its logo. Ask what is genuinely distinctive about the product.
- Do not force a 'creative' aesthetic onto dense operational UIs; clear, efficient information architecture can be distinctive too.
- Do not imitate another product's exact branded composition or copyrighted assets. Extract design principles from user-provided references, rather than cloning unrelated content.

## 4. Build a small, deliberate design system

- Put shared, actually used tokens in `web/css/default.css`: background, surface, text, muted text, borders, accent, focus, semantic states, typography, widths, spacing, radius, and shadow. Map relevant inherited UI variables and inspect computed styles.
- Use a coherent spacing rhythm (e.g. 4/8/12/16/24/32/48/64/96px) without mechanically making every section identical. Keep related content visibly closer than unrelated content.
- Select typography for the concept, readability, platform support, and Cyrillic coverage. Prefer robust font stacks by default; use carefully selected project-provided or appropriately licensed fonts when justified. Verify loading and fallback rather than assuming a font exists.
- Use purposeful type roles: body/form approximately 16–18px, supporting text 14–16px, controls 14–16px; section and display sizes scale with content and viewport rather than a fixed giant heading. Prefer `rem` and bounded `clamp()`; test resulting wraps.
- Aim for readable prose measures (~45–75 characters on roomy layouts), body line-height ~1.45–1.65, and display line-height ~1.05–1.2 as starting points. Do not shrink essential text merely to reproduce a screenshot.
- Keep accessible contrast; style light/dark themes coherently if exposed. Never use faint text or low-opacity controls as an aesthetic shortcut.
- Use real provided or suitably licensed images; verify files load, crops preserve their subjects, and responsive sizes do not shift layout. Do not fabricate product screenshots or leave broken/placeholder assets.
- One consistent icon vocabulary. Motion should clarify feedback and hierarchy, not impede use; respect reduced motion and SSR visibility.

## 5. Design responsive layouts, not stacked desktop mockups

- Begin with a usable narrow layout, enhance with fluid containers, grid/flex and sensible gutters (~16–24px on phones). Choose breakpoints when composition fails, not by device labels. Prefer normal flow over unnecessary absolute positioning.
- Reconsider navigation, section order, headings, cards, short benefit strips, image crops, controls and CTAs at mobile widths. Keep labels and icons grouped; no orphaned arrows or one-word-per-line benefit items.
- Allow text to wrap naturally. No clipping, ellipsis, `nowrap`, hard-coded heading `<br>`, fixed-height text cards, or whole-page overflow hiding to conceal composition failures.
- Give flex/grid children `min-width: 0` where necessary; use `minmax(0, 1fr)`, contain long URLs, and diagnose the true overflow source. Don't use `word-break: break-all` on regular prose.
- Use touch-friendly controls (around 44px where practical), visible keyboard focus, accessible mobile menus, and sticky elements that do not obscure content.
- Maintain essential content and actions at 200% zoom, short landscape viewports, and with reduced motion. Never disable zoom or scale the whole site down to fit.

## 6. Implement using native project conventions

- Build page-specific GOSH compositions, using documented routing/layouts/components and inherited primitives when they truly fit. Do not inherit the welcome page's layout as a design mandate.
- Keep reusable semantic structure separate from project-specific visual arrangements. Use existing Tailwind v4/PostCSS setup, design tokens, and scoped CSS; obey `AGENTS.md` restrictions.
- Make controls real and functional. Account for pending, empty, error, disabled and success states. Preserve data correctness and server-side authority.
- For new marketing pages, use documented `data-gosh-reveal` sparingly; verify animations on initial SSR load, scrolling, and client navigation, with visible content when JS is disabled or reduced motion is selected.

## 7. Visual QA: mandatory for substantial website work

Run the **built** website when possible (`task build` + `task server`; `task preview` is not the production build). Use a browser/Playwright when available. Do not declare success after JSX/GOSH/CSS compilation alone.

1. Inspect **320, 360, 390, 768, 1024, and 1440 CSS-px widths**, including realistic heights and a short landscape viewport. For small fixes, use affected width plus adjacent breakpoints.
2. Capture and actually examine full-page screenshots, at minimum one narrow and one wide. Review the first viewport, middle sections, footer and each distinct page template; resize between sample widths.
3. Evaluate against the design direction: hierarchy, composition, typography, alignment, rhythm, content density, imagery, product identity, consistency, and clarity of primary actions.
4. Name the **three to five highest-impact defects** (including aesthetic weaknesses, not just CSS bugs). Fix those defects and reinspect the affected viewports. Iterate until obvious issues are gone; do not merely list recommendations.
5. Check clipping, overflow, orphaned words, font fallback, broken assets, contrast, unintended scroll, sticky overlays, controls, and 200% zoom. When content fails, fix owning composition rather than masking with CSS hacks.
6. Exercise navigation, menus, forms, primary actions, focus/keyboard, touch, pending/empty/error/success states, reduced motion, and enabled themes. Inspect console errors and failed asset requests.
7. Test fragment links (`#section`, `/#section`, and cross-page variants) by clicking and direct loading, after client navigation, with back/forward, and under sticky headers. Confirm target is unique and genuinely scrolled into view; use documented goserver behavior rather than competing global click handlers.
8. Verify entrance reveals after first load, scroll and client navigation; ensure primary content remains readable with reduced motion and JavaScript disabled.

If browser tools/runtime are unavailable, perform build/static verification and clearly report that **visual and interaction acceptance remains unverified**. Never claim that media queries, HTTP 200, or a build prove a responsive design. Treat a user-provided screenshot of a bug as regression evidence and verify the fix at its width and adjacent widths.

## Final acceptance bar

The result must be product-specific, visually composed rather than template-assembled, understandable, accessible, responsive, functional, and consistent. Keep implementation focused and report only the concrete design choices, verification performed, remaining limitations, and necessary follow-up. Do not present a design plan as the finished implementation.
