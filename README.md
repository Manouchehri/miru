# Miru
A website monitoring tool that periodically checks a site for changes.

Miru, pronounced roughly like me-roo, is a tool developed largely for use by the [Environmental Data & Governance Initiative](https://envirodatagov.org/), who initiated in Toronto a movement with the goal of [archiving climate data](http://www.cbc.ca/news/technology/university-toronto-guerrilla-archiving-event-trump-climate-change-1.3896167) before President Trump, who [denies the existence of climate change](https://www.washingtonpost.com/news/energy-environment/wp/2016/12/11/trump-says-nobody-really-knows-if-climate-change-is-real/), has the opportunity to have the data [removed from public access](http://www.reuters.com/article/us-usa-trump-epa-climatechange-idUSKBN15906G) or destroyed entirely.

Miru functions as something of a "glorified [cron job](https://en.wikipedia.org/wiki/Cron) runner" with a web interface.  It allows for participants of archiving events to register an account and make requests to have websites worth archiving monitored for changes, so that other tools can scrape and archive said sites.  Users with administrator privileges are able to then review such requests, write a Python/Ruby/Perl script to check the requested site for changes, and upload their script to Miru, which will run the script in specified intervals to [generate reports](https://github.com/zsck/miru/blob/master/docs/reporting.md) which administrators and other tools will be able to use to determine when a site needs to be revisited.

## Getting Started

1. Read our [code of conduct]() and learn [how to get in touch with us]().
2. Learn how to [build and run Miru locally]() for development.
3. Read the project's [contributing guide]() to and [outstanding issues](https://github.com/zsck/miru/issues) learn how to help build Miru.
4. Read about [how to deploy Miru]() and related advice.