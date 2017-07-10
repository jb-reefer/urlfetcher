# urlfetcher

Fun url fetcher for a tech test, pulls urls from a csv file, gets the page described, and returns whether it contains a provided term. 

## Design

This is a fairly straightforward design - I spawn 20 worker threads (by default, you can change this value with --threads) and they pull URLs off of a channel, and search for the term. They then write responses to a result channel, so we can output to results.txt. I used a results channel and a waitgroup instead of a reader/writer pattern because the size of the data meant the buffer was small, and not worth the headaches of race conditions. It would have been simple to read off of the results channel as it was written, but I'm very against premature optimization, and the collecting->reading code is very fast.

The only part of this project that I don't love is that http.Get in Go is somewhat limited, so I had to perform a bit of url banging myself. The quick logic I threw together seems to work most of the time, but I'd love critiscm and tips for the future. Happily, most of the errors it produces seem to be from URLs that *actually* do not exist, but a few certainly do. 

## Testing

There are unit tests in [urlfetcher_test.go](urlfetcher_test.go) - these are vanilla go unit tests, and all run happily. I originally had a few for failing URLs, but due to the way httpClient works, they would hang until a timeout, which doesn't produce great logging. Otherwise, I test a few of my methods and public members. I love the built-in testing - more languages should ship with a testing framework as a part of the standard library.

## Lessons learned/Conclusion

Go is fun, and go channels are a ton of fun. The builtin HTTP libraries seem more more geared towards known good URIs, not short urls, which I suppose is very in line with the KISS philosophy of the language in general - they return what you ask for, not what a remote server wants to give you. Otherwise, working with data and the error handling are great and very intuitive. One of the most interesting parts was that the list of URLs appears to be the Alexa Top 500 from a few years ago, and a LOT of the addresses don't resolve anymore - who knew there was so much churn in those sites?! Topsy.com, Sitemeter - those were pretty big names at one point  - this project was cool!
