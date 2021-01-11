# GoRoller
A dice-rolling discord bot programmed in Go!
Uses [GoDiscord](https://github.com/bwmarrin/discordgo), by @bwmarrin

## Installation
___

Click on this link: https://discord.com/api/oauth2/authorize?client_id=798054898527633420&permissions=2048&scope=bot

It will bring you to a server selector, where you can install it on any server you have the correct permissions for. 

## Usage
___

In any channel that the bot has permission to view and send messages, input `/r` followed by the roll you want to perform.

### Documentation
___

**Rolling Dice**

`/r XdY` - Roll a number of Y-sided dice equal to X. If omitted, X becomes 1 and Y becomes 0

*Examples*

```
/r 1d20  # Roll one 20-sided die
/r 3d8   # Roll three 8-sided dice
/r d10   # Roll one 10-sided die
```
___
	
**Math**
	
`X [+, -, *, /, ^] Y` - Do math! Currently supports addition, subtraction, multiplication, division, and exponentials. All of this can also be done with dice! If omitted, both X and Y become 0

*Examples*

```
/r 2 + 2     # Add 2 and 2
/r 500/4     # Divide 500 by 4
/r 2d10 * 5  # Roll 2d10 and multiply the result by 5
```
___

**Parentheses**
	
`(X+Y)*Z` - Order of operations is followed, meaning things in parentheses are resolved first, even when nested

*Examples*

```
/r (1d20+3)*2      # Roll a d20, add 3, and multiply the result by 2
/r (2d4)d6         # Roll 2d4, and roll a number of d6 equal to the result
/r (15/(2d6-3))^2  # Roll 2d6, subtract 3 from it, divide 15 by that result, and square the whole thing
```
___