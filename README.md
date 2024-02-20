# Honkai Star Rail Theorycrafting Tool

TDLR: A flexible way of sheeting damage, comparing relic sets or supports.

Has support for debuffed enemies (think Silverwolf, Topaz or Resolution Shines As Pearls of Sweat), buffs on specific attacks or team buffs.

<picture>
  <img style="width: 128px;" alt="A picture of Hook" src="images/hook.png">
</picture>

## Executing it from the source

 - Install Go 1.20+ (https://go.dev/doc/install)
 - Execute the following command to run it without making an exe:
   - `go run .\cmd\hsrtctsheets\` (from the root directory of the repository)
   - `go run {FULL_PATH}\cmd\hsrtctsheets\` (from any directory)

## Input sheet

The app looks for a file named "HSRTCT-config.xlsx", needs to be on the same folder.

**[Example (my own input sheet)](https://docs.google.com/spreadsheets/d/1rn1X0IpP8FRU7MoBQ0SFO1BZ14MRRlgNSy-XW9rRM34)**

Scenarios example (cropped):
![A cropped screenshot of the Scenarios page from the input excel](/images/inputScenarios1.png)

Relic sets example (cropped):
![A cropped screenshot of the Relic Sets page from the input excel](/images/inputRelicSets1.png)

Most stuff is coded as a buff (including enemy debuffs).

A buff consists of:
 - The stat it buffs (or debuffs)
 - The number value of the buff
 - The damage tag condition (if it has one)
 - The element condition (if it has one)

## Output

Scenario results:
![A screenshot of the Scenario results page from the output excel](/images/output1.png)

Scenario explanation example (cropped):
![A screenshot of a Scenario explanation page from the output excel](/images/output2.png)

## Roadmap

 - [x] Create the calculator package hsrtct
 - [x] Decide the database to use
    - [x] A google sheet input kek
    - [x] An in memory db kek
 - [x] Make the config sheet
   - With a **Character**, using a **LC**
   - With 1-5 **Enemies**
   - With **N Attacks** to Enemy with an **amount**
   - **Final damage**: The addition of all attacks damage times their amount
 - [x] Better explanations
 - [x] Add "Dmg reduction" stat
 - [ ] Implement stuff like "Crit dmg taken" on enemies
 - [ ] Implement a way to add break damage on scenarios
 - [ ] Implement a way to calc heals/shields
