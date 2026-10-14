# Team skills

Project skills for Claude Code. Everyone who opens this repo in Claude Code gets them, with no per-user install.

They are vendored unchanged from upstream. Don't edit
them here. To update them, follow the steps under Updating.

## How they are used

| Skill | When |
|---|---|
| `ponytail` | Every change. Pick the laziest solution that is actually correct: reuse before writing, delete before adding, no speculative abstractions. |
| `code-review-and-quality`, `code-simplification` | Before every commit: re-check the diff with both. |
| `security-audit` | Security questions and reviews. Use guidance mode by default. Run the full six-phase audit only when someone asks for one, and keep its output outside the repo. |
| The rest | When the task matches the skill description, e.g. `frontend-ui-engineering` for UI, `debugging-and-error-recovery` for failures, `test-driven-development` for new logic. |

The hard rules in `AGENTS.md` and `CONSTRAINTS.md` win when a skill disagrees with them: no comments, complexity ≤ 15, files ≤ 600 lines.

## Sources

| Skill | Upstream | Folder hash |
|---|---|---|
| `ponytail` | [DietrichGebert/ponytail](https://github.com/DietrichGebert/ponytail) | `cb6e53477eaf` |
| `api-and-interface-design` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `3234ce9609d7` |
| `browser-testing-with-devtools` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `42be59520831` |
| `ci-cd-and-automation` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `34db321777a1` |
| `code-review-and-quality` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `2c0b966d411b` |
| `code-simplification` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `15b6e2192418` |
| `constraint-driven-development` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `8aa3f6a661bc` |
| `context-engineering` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `5da634152562` |
| `debugging-and-error-recovery` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `c3b5c340bae3` |
| `deprecation-and-migration` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `452a130309af` |
| `documentation-and-adrs` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `e0cbd8df50a6` |
| `doubt-driven-development` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `2a25398c4f73` |
| `frontend-ui-engineering` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `8fdd24f2b93c` |
| `git-workflow-and-versioning` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `92d87ae1fe98` |
| `idea-refine` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `14f474d8f513` |
| `incremental-implementation` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `a3e2fab3e3b9` |
| `interview-me` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `2a391376521d` |
| `observability-and-instrumentation` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `1e8f1c0144e2` |
| `performance-optimization` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `c0cf6ae718e8` |
| `planning-and-task-breakdown` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `3045f95a156d` |
| `security-and-hardening` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `79b5adb3e84a` |
| `shipping-and-launch` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `0be8f859c53d` |
| `source-driven-development` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `29e1dcc180d5` |
| `spec-driven-development` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `4d36b30efbe4` |
| `test-driven-development` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `b812b07dc0c5` |
| `using-agent-skills` | [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) | `81e8196f4dd4` |
| `security-audit` | [cloudflare/security-audit-skill](https://github.com/cloudflare/security-audit-skill) | `ccbc33ed119b` |

## Updating

```sh
npx skills add DietrichGebert/ponytail
npx skills add addyosmani/agent-skills
npx skills add https://github.com/cloudflare/security-audit-skill --skill security-audit
```

Then copy the updated folders over the ones here, re-read the diff, and commit them in a `chore(skills): …` commit of their own.

## Licenses

All three upstreams are MIT licensed.

- ponytail: Copyright (c) 2026 DietrichGebert
- agent-skills: Copyright (c) 2025 Addy Osmani
- security-audit-skill: Copyright (c) 2025-2026 Cloudflare, Inc.

> Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
>
> The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
