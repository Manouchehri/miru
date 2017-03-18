# Miru User Guide

This guide explains how to use each of Miru's features from both the perspective of a regular user/archiver as well as from the perspective of an administrator.

## Archivers

At the time of this writing, there is only one feature supported by Miru for non-administrator archivers to use. 

### Making a request to have a site monitored

![making a request](https://github.com/zsck/miru/blob/master/docs/screenshots/making-requests.png)

Users can click the **Request** link at the top right of the index page to be brought to a page with a form for them to fill out with information about the site they would like to have monitored for changes. Users are asked to provide:

1. The address of the website/page they issuing the request for.
2. Instructions that might help an administrator write a script to accurately check for noteworthy changes.

Note that any query string information (e.g. `?user=sensitive@info.com&ip=127.0.0.1`) is removed by Miru upon receipt of a request, as this information could potentially contain sensitive information identifying the archiver.  Users should be educated about this part of an URL and include information about relevant parts in the instructions section of the form if the data is in fact required to access the page.

## Administrators

Administrative users have access to all of Miru's functionality. Upon logging in, the **Request** link shown to non-administrative users will be replaced with an **Admin Panel** link bringing the administrator to a page containing links to other pages wherein actions of interest can be performed.

### Viewing monitor requests

![viewing pending requests](https://github.com/zsck/miru/blob/master/docs/screenshots/viewing-requests.png)

This page shows a list of all pending requests made to have sites monitored.  Here, administrators can reject requests that have either already been fulfilled or will not be fulfilled.

### Fulfilling monitor requests

![fulfilling monitor requests](https://github.com/zsck/miru/blob/master/docs/screenshots/fulfilling-requests.png)

Upon clicking the **Approve** button on the **Pending monitor requests** page, the administrator will be able to upload a [report-generating monitor script](https://github.com/zsck/miru/blob/master/docs/reporting.md) and specify how and when to run it. Currently, scripts can either be written in Python, Ruby, or Perl.

Miru must also be told how frequently to run the script. The first numeric input for **Time to wait between runs (minutes)** allows you to specify how often to run the script, in minutes. The default here is `1440`, which is precisely one full day, or 24 hours. It is advised that monitor scripts be setup to run as infrequently as possible, to avoid having websites flag Miru for suspicious activity.

Finally, Miru can be told how long the script being uploaded should be expected to run for. The **Expected script runtime (seconds)** input allows you to specify the maximum number of seconds that Miru should allow a monitor script to run for before terminating it in order to prevent system overloads caused by erratic script behavior. *Note that this feature is not currently implemented*.

### Viewing and promoting archivers

![viewing archivers](https://github.com/zsck/miru/blob/master/docs/screenshots/viewing-archivers.png)

The admin panel also links to a page that lists all registered archivers. Those that are not already administrators can be made so by clicking **Make Admin** in the row containing their email address.

### Viewing monitor reports

![viewing reports](https://github.com/zsck/miru/blob/master/docs/screenshots/viewing-reports.png)

Every time a monitor script is run, it is expected to produce a report containing information about the site it checked and specify whether any noteworthy changes occurred. Reports are color-coded based on how significant any changes detected to a site are determined to have been with:

* Green reports indicating little to no change.
* Yellow reports indicating some note-worthy changes having taken place.
* Red reports indicating that a significant change, requiring immediate investigation, has taken place.

By clicking on a report, it will be expanded to show information about the monitor script including when it was last run, the checksum it computed of the information it checked, where the script itself is located on disk, and the message left for the administrator by the script.

