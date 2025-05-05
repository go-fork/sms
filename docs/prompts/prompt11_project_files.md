# Prompt 11: Project Files

## Objective
Set up essential project files including license, gitignore, GitHub templates, and CI/CD configuration.

## Required Files to Create

1. Project files:
   - `LICENSE` - Open source license
   - `.gitignore` - Git ignore file

2. GitHub-specific files:
   - `/.github/ISSUE_TEMPLATE/bug_report.md` - Bug report template
   - `/.github/ISSUE_TEMPLATE/feature_request.md` - Feature request template
   - `/.github/PULL_REQUEST_TEMPLATE.md` - Pull request template
   - `/.github/workflows/ci.yml` - GitHub Actions CI workflow

## Implementation Requirements

### License
- Create a `LICENSE` file with:
  - MIT License text
  - Current year and copyright holder information

### Git Configuration
- Create a `.gitignore` file with common Go exclusions:
  - Compiled binaries
  - Vendor directory
  - IDE-specific files
  - Temporary files and directories
  - Local environment and configuration files

### GitHub Templates
- Bug report template with sections for:
  - Bug description
  - Steps to reproduce
  - Expected behavior
  - Environment details
  - Additional context

- Feature request template with sections for:
  - Problem description
  - Proposed solution
  - Alternatives considered
  - Additional context

- Pull request template with sections for:
  - Description of changes
  - Related issues
  - Type of change
  - Checklist for the submitter

### CI/CD Configuration
- GitHub Actions workflow that:
  - Runs on push to main and pull requests
  - Sets up Go environment
  - Runs linting
  - Executes unit tests
  - Checks code coverage
  - Builds the module and examples

## Deliverables
- Complete license file
- Configured .gitignore
- GitHub issue and PR templates
- CI/CD workflow configuration
