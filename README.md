# Honkai Star Rail Theorycrafting Tool

<picture>
  <img style="width: 128px;" alt="A picture of Hook" src="images/hook.png">
</picture>

## Input sheet

**[Example (my own input sheet)](https://docs.google.com/spreadsheets/d/1rn1X0IpP8FRU7MoBQ0SFO1BZ14MRRlgNSy-XW9rRM34)**

Scenarios example (cropped):
![A cropped screenshot of the Scenarios page from the input excel](/images/inputScenarios1.png)

Relic sets example (cropped):
![A cropped screenshot of the Relic Sets page from the input excel](/images/inputRelicSets1.png)

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
 - [ ] Add stuff like "Crit dmg taken" on enemies
