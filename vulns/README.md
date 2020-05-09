# Vulns

## General

Look at CVE-2019-1315 for a commented "template".

## Add new CVE

- search kb for major version on https://www.catalog.update.microsoft.com/Search.aspx
- click cumulative update, goto package details
- add all "this update has been replaced by"
- repeat for all major versions

There is also a script to automate this in the wesng repository:
https://github.com/bitsadmin/wesng/blob/master/muc_lookup.py

## Todo

- Common other privesc vectors, e.g. SeImpersonate, PrivEscChecker etc.