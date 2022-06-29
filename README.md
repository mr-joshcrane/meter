# Meter

Meetings can be boring and expensive. **Meter** can't measure participant boredom, but it can measure how expensive a given meeting is! You can use **Meter**:
___
### For Simple Cost Estimation
Handy for calculating the total cost of spending time on x activity given an hourly rate!
```bash
meter -rate=100 -duration=1h
The total current cost of this meeting is $ 100.00
```
___
### As A Running Cost Counter
Handy for keeping an eye on the mounting costs of the meeting, so you can trade it off against the value being generated!
```bash
meter -rate=10000 -duration=1h -ticks=5s

The total current cost of this meeting is $ 13.89
The total current cost of this meeting is $ 27.78
The total current cost of this meeting is $ 41.67
... <one hour later> ...
The total current cost of this meeting is $ 10000.00
```
---
### Getting Participants Rates Interactively
Don't know off the combined rate of all participants off the top of your head? No problem, just omit the -rate flag and be prompted interactively.
```bash
meter -duration=1h
Please enter the hourly rates of all participants, one at a time. ie. 150 OR 1000.50
Please enter the hourly rates of the next participant
If all meeting participants accounted for, type Q and enter to move on.
```
___

From the root of the project, you can **build** the executable with
``` bash
go build -o meter cmd/main.go
```