# go-shopping-list

## Summary

Do you find making shopping lists BORING? Because I sure do... This repo is to automate the creation of shopping lists and to practice my Golang skills! 

## Table Of Contents

- [go-shopping-list](#go-shopping-list)
  * [Summary](#summary)
  * [Table Of Contents](#table-of-contents)
  * [Demos](#demos)
    + [Excel](#excel)
    + [Reminders](#reminders)
    + [Terminal](#terminal)
  * [How does it work?](#how-does-it-work)
  * [Maintaining Professional Standards](#maintaining-professional-standards)
  * [GitHub Actions](#github-actions)
    + [Testing](#testing)
    + [Linters](#linters)
    + [Content Checking](#content-checking)
    + [Project Management](#project-management)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>

## Demos

### Excel 

If no shopping workflow has been specified the code will create an excel sheet which users can print off to take shopping!

https://user-images.githubusercontent.com/32711718/190857227-7b52b057-c60a-4650-a868-770b92e90f28.mov

### Reminders

When you are running on mac and have the shopping.workflow file present this is how the code operates:

https://user-images.githubusercontent.com/32711718/189308917-77185ca7-6811-4a2f-b3c4-a85158703dde.mov

### Terminal 

If you are not running on a macbook and don't have the shopping.workflow present then the system will just print that it is pretending to add the ingredients so you can still use the system. See the terminal demo below:

https://user-images.githubusercontent.com/32711718/189309910-072d7b0d-bffa-4661-ad5a-01fb5aaff30e.mov

## How does it work? 

1. The codebase works by firstly reading all the JSON files in the recipes folder. 
2. Once the JSON's have been processed they are passed to the GUI and displayed in a list
3. Once multiple recipes are selected and the submit button the workflow is ran.
4. There are currently 3 workflows as demoed above.

## Maintaining Professional Standards

To ensure code is professional and extendable the following rules should be followed:

1. Apply unit testing wherever possible and aim for 80% coverage in packages
2. Use `golangci-lint run` and the revive linter with all rules enabled: https://github.com/mgechev/revive 
3. Using Interfaces to mock results using the following module: https://github.com/vektra/mockery 

See the PR Template fo more checks to follow.

## GitHub Actions

### Testing

The pkg-cov workflow runs all go tests and ensures pkg coverage is above 80%.

![example event parameter](https://github.com/gowhale/go-shopping-list/actions/workflows/pkg-cov.yml/badge.svg?event=push)

The pages workflow publishes a test coverage website everytime there is a push to the main branch. The website can be found here: https://gowhale.github.io/go-shopping-list/#file0

![example event parameter](https://github.com/gowhale/go-shopping-list/actions/workflows/pages.yml/badge.svg?event=push)

### Linters

The revive workflow is executed to statically analsye go files: https://github.com/mgechev/revive

![example event parameter](https://github.com/gowhale/go-shopping-list/actions/workflows/revive.yml/badge.svg?event=push)

The golangci-lint workflow runs the golangci-lint linter: https://github.com/golangci/golangci-lint

![example event parameter](https://github.com/gowhale/go-shopping-list/actions/workflows/golangci-lint.yml/badge.svg?event=push)

### Content Checking

The authenticate workflow checks that the recipe .json files are in the correct format.

![example event parameter](https://github.com/gowhale/go-shopping-list/actions/workflows/authenitcate.yml/badge.svg?event=push)

### Project Management

The issue workflow adds a new issue to the projects Kanban board: https://github.com/users/gowhale/projects/1

![example event parameter](https://github.com/gowhale/go-shopping-list/actions/workflows/issue.yml/badge.svg?event=push)

The cut release workflow creates a binary executable everytime a release is published. The binary file is attached to the release.

![example event parameter](https://github.com/gowhale/go-shopping-list/actions/workflows/cut-release.yml/badge.svg?event=push)


