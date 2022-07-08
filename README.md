# MapgGridBot
Is a Telegram bot that generates coordinate grid for a chosen area of a map.
## Usage
You send a message to the bot with a parameters as described below. Bot sends you a file of a grid in .gpx format.
Description:

> Enter latitude and longitude in decimal format: xx.xxxxxx, yy.yyyyyy
> First coordinates are for north-west (left top) corner.
> Second coordinates are limiting grid size by south-east corner (optional).
> Next parameter is the name of the grid (optional).
> Last parameter is step size of the grid in meters (optional). Ex: 50, 100, 200, 500 etc.
> Maximal grid size is 26 x 26 points.
> Default grid step is 100 meters.

## Examples of messages to the bot
#### Example 1:
```
60.0000, 30.0000
```
#### Example 2: 
```
61.231715, 30.024984 61.223463, 30.051560 Heposaari
```
#### Example 3: 
```
59.975964, 30.268458; 59.949592, 30.336530 Petrogradsky_District 200
```
#### Example 4: 
```
60.061522, 30.142147 Saint-Petersburg 1000
```
## Deploy
```
docker-compose up --build -d
```
