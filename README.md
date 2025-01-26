# Web Scrape Greek news websites

The app goes and fetches data from various news websites from Greece. It has intergration with my personal hosted postgres bbut the web scrapping part can be used by anyone.
My main motivation was I was bored to look each page in the browser so i am using the data locally.

All websites are server side rennder but a lot of them supported RSS feed which was nice. But for the ones that didnät , they needed quite a bit of searching in the html source code to take the data.
Dates harmonisation was needed extra care in order to save o have a common format for the db as well with timezone . Some where UTC some local time.
