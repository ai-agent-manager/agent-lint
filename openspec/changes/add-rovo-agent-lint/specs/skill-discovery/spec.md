## MODIFIED Requirements

### Requirement: Exactly One Definition File

Each candidate directory SHALL contain exactly one of `SKILL.md` or `rovo-agent.yaml`.

#### Scenario: Skill directory

- **WHEN** a directory contains `SKILL.md` and not `rovo-agent.yaml`
- **THEN** it is validated as an Agent Skill

#### Scenario: Rovo-only directory validated

- **WHEN** a directory contains `rovo-agent.yaml` and not `SKILL.md`
- **THEN** it is validated as a Rovo agent (structural `rovo-agent.yaml` lint)

#### Scenario: Both definition files

- **WHEN** a directory contains both `SKILL.md` and `rovo-agent.yaml`
- **THEN** an error diagnostic with rule ID `skill/asset-type-conflict` is emitted

#### Scenario: Neither definition file

- **WHEN** a directory contains neither `SKILL.md` nor `rovo-agent.yaml`
- **THEN** an error diagnostic with rule ID `skill/asset-type-missing` is emitted
