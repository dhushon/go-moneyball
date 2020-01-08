# go-moneyball
Project to Grab and Un-Marshal Sports Event Data to support downstream player, team, squad, workout, performance improving advanced analytics and machine learning.  First core step is preparing some basic structures to fetch remote data - both historic and game/practice flow.

## to build
```
cd moneyball
go build
ls /moneyball
./                 espnScoreBoard.go  moneyball.go
../                moneyball*
./moneyball
```

```
ScoreBoardService: 6 scores for date 2019-12-30 retrieved
Response: &main.ScoreBoard{Leagues:[]main.League{main.League{ID:"46", UID:"s:40~l:46", Name:"National Basketball Association", Abbreviation:"NBA", Slug:"nba", Season:main.SeasonDef{Year:2020, StartDate:main.espnTime{wall:0x0, ext:63705250800, loc:(*time.Location)(nil)}, EndDate:main.espnTime{wall:0x0, ext:63729183540, loc:(*time.Location)(nil)}, Type:main.SeasonType{ID:"2", Type:2, Name:"Regular Season", Abbreviation:"reg"}}...
```

## to test
```
cd moneyball
go test
```

You can find examples of the raw JSON from the ESPN and NBA API's in the [examples/json](https://github.com/dhushon/go-moneyball/tree/master/examples/json) directory

## to contribute
pulls are very welcome
