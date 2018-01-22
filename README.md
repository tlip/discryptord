discryptord v0.4.2
===========

A Discord bot that will amaze your friends by emitting a chart for any magic internet currency.

![Example](assets/graph.png)

## Usage

```run
!<TICKER> [?<BASE>] [?-(3d|w|m|3m|6m|y)] [?-(rsi)]
```

### Options

#### Time Range

 Examples:
```run
# ETH/BTC ()
!eth btc      # 24h/1d in 10m ticks 
!eth btc -w   # 7d/1w in 1h ticks
!eth btc -m   # 30d/1m in 6h ticks
!eth btc -3m  # 90d/3m in 6h ticks

# ETH/USD
!eth          # 24h/1d in 10m ticks 
```

As you can see, if `<BASE>` is ommited, then it defaults to `"usd"`.
The `-3d` changes the output range from 24h/1d to 72h/3d. 
The `-w`  changes the output range from 24h/1d to 7d/1w. 
The `-m`  changes the output range from 24h/1d to 30d/1m. 
The `-3m` changes the output range from 24h/1d to 90d/3m. 
The `-6m` changes the output range from 24h/1d to 180d/6m. 
The `-y`  changes the output range from 24h/1d to 365d/1y. 


#### Technical Analysis

As of right now, we're testing customizable TA support.
Only RSI is available right now, but I expect to bring many more indicators at some point in the future.
These can be chained in tandem with `-3d`, `-w`, `-m`, etc.

Examples:
```run
!eth -rsi
!eth -w -rsi
!eth -rsi -m
```

## Install

The ancient method of invoking this minor god has long been forgotten...
So if you want to invite our evil friend into your realm, our leading scientists believe the magic cantation with the most promise goes something like:

```zsh
% go get github.com/flamingyawn/discryptord
```

It's simple to summon our bot onces its dwelling in your proximity.
Just recite the following on a night with a new moon during the eve of planets' alignment after drawing salt circle around 3 candles with black flames:

```zsh
% go install
% go build
```

No promises, doe.

## Getting Started

First, create a bot friend on discord.
Then invite it to the channel(s) you want it to stalk...
Next, copy its **BOT_TOKEN** to your magic clipboard.

Go into this package's home directory and whisper this into your terminal:

```zsh
% echo "<BOT_TOKEN>" >> .token
```

Finally, speak

```zsh
% sh start.sh
```

to start the _real_ fun.

## halp donate plz

Of course, I made this little guy with :heart: and offer him completely for free!
But if you find him useful and want to help keep me developing, then I love you long time.
ETH: `0x0f8c31fa23b21f23565db1e0938ebf41dd2ec5cd`

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for full license text.
