# Loan Calculator
Try it out here: https://tekkamanendless.github.io/loan-calculator/

## History
I was frustrated with how limited the typical online loan calculators were.
My original use case was:

1. I had a solar loan for $39,125.00.
1. The monthly payment of $187.48 was assuming that 30% would be paid back in a lump sum within the first 18 months ($11,737.50).
1. Starting at the first payment, I added $52.52 in principal to bring the total up to a round $240/month.
1. About a year in, I received a $5,000 solar grant, so I wanted to see what impact that money would have on the loan.

A lot of online calculators would handle one or two of my details, but not all of them.
And so that's why I made this tool (tl;dr: I should immediately pay the $5,000 down as principal).

## Usage
Feed the tool a loan string and any number of optional extra-payment strings.

A loan string looks like this:

```
amount <amount> rate <rate> months <count> payment <amount> starting <date>
```

Exmaples:

* My solar loan (see above).
   ```
   amount 39,125.00 rate 4.99 months 240 payment 187.48 starting 2019-10-01
   ```

An extra-payment string looks like this:

```
<amount> monthly [starting <date>] [ending <date>|count <count>]
<amount> once on <date>
```

Examples:

* Add an additional $52.52 to principal per month.
   ```
   52.52 monthly
   ```
* Make a one-time $12,000 payment to principal on November 10, 2020.
   ```
   12,000 once on 2020-11-10
   ```
* Make a one-time $5,000 payment to principal on December 20, 2020.
   ```
   5,000 once on 2020-12-20
   ```
* Split that same $5,000 by making 50 additional monthly payments of $100.
   ```
   100 monthly starting 2020-12-20 count 50
   ```
