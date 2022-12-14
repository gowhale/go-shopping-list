# Description

Please include a summary of the change and which issue is fixed. Please also include relevant motivation and context.

Fixes # (issue)

## Type of change

Please delete options that are not relevant.

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] This change requires a documentation update

# How Has This Been Tested?

Please describe the tests that you ran to verify your changes. Provide instructions so we can reproduce. Please also list any relevant details for your test configuration

- [ ] Test A
- [ ] Test B

# Checklist:

- [ ] I have added test coverage to new lines of code
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] I have ran `go test ./...` locally
- [ ] I have ran `golangci-lint run`
- [ ] I have ran the revive linter using the `revive.toml` settings
- [ ] I have ran `go run ./cmd/pkg-cover`
- [ ] I have ran `go run ./cmd/authenticate`

Once all these checks have been ticked off the GitHub actions should pass and the PR is ready for review.
- [ ] The GitHub Actions On This PR Have Passed
