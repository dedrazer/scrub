# scrub

## todo
- test refactored methods (before: 73.1% files, 53% statements)
- desired functionality:
  - main input args: starting stack, cashout at (min/max)
  - output: ideal strategy, oneCreditAmount, chances of reaching min/max cashout amount
- strategies become map[streak int] betmultiplier float
- bank
- fix SimulationResults when there are no rebuys
- auto stop when results are stable