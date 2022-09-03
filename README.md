# go-shopping-list

## Summary

Do you find making shopping lists BORING? Becuase I sure do... This repo is to automate the creation of shopping lists and to practice my Golang skills! 

## Demo

https://user-images.githubusercontent.com/32711718/188265124-ef39dce0-62ee-43b5-974f-4c0bfc5b3687.mov

## How does it work? 

1. The codebase works by firstly reading all the JSON files in the recipes folder. 
2. Once the JSON's have been processed they are passed to the GUI and displayed in a list
3. When one of the list items is clicked it will run automation scripts which add a reminder to my reminders app
4. The reminders app is an app available on mac and iphone.
5. The iCloud then syncs the items in my list to my phone so I can go to the shops and get the stuff I need.

## Maintaining Proffesional Standards

To ensure my code is proffesional and extendable I followed these rules when making changes:

1. Apply unit testing wherever possible and aim for 80% coverage in packages
2. Use `golangci-lint run` and the revive linter with all rules enabled: https://github.com/mgechev/revive 
3. Using Intefaces to mock results using the following module: https://github.com/vektra/mockery 
