# scrub

## todo
- split
- insurance

```
{"level":"info","ts":1674729299.797014,"msg":"initialising blackjack"}
{"level":"info","ts":1674729299.801733,"msg":"dealing round","playerBets":[{"Player":{"Name":"Martin","Credits":1000},"BetAmount":50}]}
{"level":"info","ts":1674729299.801846,"msg":"dealer hand","card":"8 of Diamonds"}
{"level":"info","ts":1674729299.8018508,"msg":"player hand","player":1,"hand":1}
{"level":"info","ts":1674729299.801856,"msg":"hand","cards":"7 of Spades, 4 of Hearts","value":"11"}
{"level":"info","ts":1674729299.80186,"msg":"playing round"}
{"level":"info","ts":1674729299.801867,"msg":"turn","player":1,"hand":1}
{"level":"info","ts":1674729299.80187,"msg":"dealer hand","card":"8 of Diamonds"}
{"level":"info","ts":1674729299.8018749,"msg":"hand","cards":"7 of Spades, 4 of Hearts","value":"11"}
Take card? (y/N)
y
{"level":"info","ts":1674729305.912905,"msg":"got card","card":"3 of Spades"}
{"level":"info","ts":1674729305.9129422,"msg":"hand","cards":"7 of Spades, 4 of Hearts, 3 of Spades","value":"14"}
Take card? (y/N)
y
{"level":"info","ts":1674729307.642781,"msg":"got card","card":"8 of Clubs"}
{"level":"info","ts":1674729307.642812,"msg":"player bust"}
{"level":"info","ts":1674729307.642818,"msg":"calculating results"}
{"level":"info","ts":1674729307.642847,"msg":"dealer result","cards":"6 of Spades, 8 of Diamonds, 8 of Diamonds","value":"22"}
{"level":"info","ts":1674729307.642857,"msg":"dealer bust"}
{"level":"info","ts":1674729307.642876,"msg":"player result","hand":1,"result":"win","credits":1050}
```

```
{"level":"info","ts":1674744597.5345159,"msg":"initialising blackjack"}
{"level":"info","ts":1674744597.542655,"msg":"dealing round","playerBets":[{"Player":{"Name":"Martin","Credits":1000},"BetAmount":50}]}
{"level":"info","ts":1674744597.542799,"msg":"dealer hand","card":"8 of Spades"}
{"level":"info","ts":1674744597.542805,"msg":"player hand","player":1,"hand":1}
{"level":"info","ts":1674744597.542811,"msg":"hand","cards":"9 of Diamonds, Jack of Diamonds","value":"19"}
{"level":"info","ts":1674744597.542816,"msg":"playing round"}
{"level":"info","ts":1674744597.542824,"msg":"turn","player":1,"hand":1}
{"level":"info","ts":1674744597.542828,"msg":"dealer hand","card":"8 of Spades"}
{"level":"info","ts":1674744597.5428321,"msg":"hand","cards":"9 of Diamonds, Jack of Diamonds","value":"19"}
Double down? (y/N)
y
{"level":"info","ts":1674744600.011611,"msg":"player bust"}
{"level":"info","ts":1674744600.01164,"msg":"calculating results"}
{"level":"info","ts":1674744600.011659,"msg":"dealer result","cards":"5 of Hearts, 8 of Spades, 6 of Spades","value":"19"}
{"level":"info","ts":1674744600.011702,"msg":"player result","hand":1,"result":"push","credits":1000}
```
